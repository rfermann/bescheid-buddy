package paper_letters_pdf_service

import (
	"bytes"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/go-pdf/fpdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

const (
	maxGermanTextRunes = 24000
	maxSubjectRunes    = 500
	maxReferenceRunes  = 200
	maxNameRunes       = 200
	maxAddressRunes    = 1000
	maxFilenameRunes   = 200
	maxSignatureBytes  = 512 * 1024
	maxSignatureSide   = 2000
	maxAddressLines    = 6
)

var (
	ErrInvalidInput     = errors.New("invalid PDF input")
	ErrInvalidAddress   = errors.New("invalid postal address")
	ErrInvalidFilename  = errors.New("invalid PDF filename")
	ErrInvalidSignature = errors.New("invalid PNG signature")

	postalCityPattern  = regexp.MustCompile(`^([0-9]{5})\s+(.+)$`)
	postOfficePattern  = regexp.MustCompile(`(?i)^postfach\s+[0-9][0-9\s-]*$`)
	houseNumberPattern = regexp.MustCompile(`[0-9]+[[:alpha:]]?(?:[-/][0-9]+[[:alpha:]]?)?$`)
)

//go:embed fonts/NotoSans-Regular.ttf fonts/NotoSans-Bold.ttf fonts/OFL.txt
var fontFiles embed.FS

type Party struct {
	Name    string
	Address string
}

type GenerateInput struct {
	GermanText         string
	Subject            string
	Sender             Party
	Recipient          Party
	Reference          string
	Filename           string
	SignaturePNGBase64 string
	Now                time.Time
}

type Result struct {
	Bytes     []byte
	Filename  string
	PageCount int
}

type normalizedInput struct {
	GermanText string
	Subject    string
	Sender     normalizedParty
	Recipient  normalizedParty
	Reference  string
	Filename   string
	Signature  []byte
	Now        time.Time
}

type normalizedParty struct {
	Name  string
	Lines []string
	City  string
}

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Generate(input GenerateInput) (*Result, error) {
	normalized, err := normalizeInput(input)
	if err != nil {
		return nil, err
	}

	pdfBytes, err := render(normalized)
	if err != nil {
		return nil, fmt.Errorf("render paper letter PDF: %w", err)
	}

	configuration := model.NewDefaultConfiguration()
	configuration.Offline = true
	if err := api.Validate(bytes.NewReader(pdfBytes), configuration); err != nil {
		return nil, fmt.Errorf("validate generated paper letter PDF: %w", err)
	}
	pageCount, err := api.PageCount(bytes.NewReader(pdfBytes), configuration)
	if err != nil {
		return nil, fmt.Errorf("count generated paper letter PDF pages: %w", err)
	}
	if pageCount < 1 {
		return nil, errors.New("generated paper letter PDF has no pages")
	}

	return &Result{Bytes: pdfBytes, Filename: normalized.Filename, PageCount: pageCount}, nil
}

func normalizeInput(input GenerateInput) (normalizedInput, error) {
	germanText := normalizeText(input.GermanText)
	subject := normalizeSingleLine(input.Subject)
	reference := normalizeSingleLine(input.Reference)
	if germanText == "" || subject == "" {
		return normalizedInput{}, ErrInvalidInput
	}
	if exceeds(germanText, maxGermanTextRunes) || exceeds(subject, maxSubjectRunes) ||
		exceeds(reference, maxReferenceRunes) {
		return normalizedInput{}, ErrInvalidInput
	}

	sender, err := normalizeParty(input.Sender)
	if err != nil {
		return normalizedInput{}, fmt.Errorf("%w: sender", err)
	}
	recipient, err := normalizeParty(input.Recipient)
	if err != nil {
		return normalizedInput{}, fmt.Errorf("%w: recipient", err)
	}
	filename, err := sanitizeFilename(input.Filename)
	if err != nil {
		return normalizedInput{}, err
	}
	signature, err := decodeSignature(input.SignaturePNGBase64)
	if err != nil {
		return normalizedInput{}, err
	}

	now := input.Now
	if now.IsZero() {
		now = time.Now()
	}
	return normalizedInput{
		GermanText: germanText,
		Subject:    subject,
		Sender:     sender,
		Recipient:  recipient,
		Reference:  reference,
		Filename:   filename,
		Signature:  signature,
		Now:        now,
	}, nil
}

func normalizeParty(input Party) (normalizedParty, error) {
	name := normalizeSingleLine(input.Name)
	address := normalizeText(input.Address)
	if name == "" || address == "" || exceeds(name, maxNameRunes) || exceeds(address, maxAddressRunes) {
		return normalizedParty{}, ErrInvalidAddress
	}

	lines := splitAddressLines(address)
	seen := map[string]bool{fold(name): true}
	normalized := make([]string, 0, len(lines))
	streetOrBox := false
	streetIndex := -1
	postalIndex := -1
	city := ""
	for _, line := range lines {
		key := fold(line)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		normalized = append(normalized, line)
		if postOfficePattern.MatchString(line) || houseNumberPattern.MatchString(line) {
			streetOrBox = true
			streetIndex = len(normalized) - 1
		}
		if match := postalCityPattern.FindStringSubmatch(line); match != nil {
			city = strings.TrimSpace(match[2])
			postalIndex = len(normalized) - 1
		}
	}
	if !streetOrBox || city == "" {
		return normalizedParty{}, ErrInvalidAddress
	}

	lines = append([]string{name}, capAddressLines(normalized, streetIndex, postalIndex)...)
	return normalizedParty{Name: name, Lines: lines, City: city}, nil
}

func capAddressLines(lines []string, requiredIndexes ...int) []string {
	maximum := maxAddressLines - 1
	if len(lines) <= maximum {
		return lines
	}
	selected := make(map[int]bool, maximum)
	for _, index := range requiredIndexes {
		if index >= 0 {
			selected[index] = true
		}
	}
	for index := range lines {
		if len(selected) >= maximum {
			break
		}
		selected[index] = true
	}
	result := make([]string, 0, maximum)
	for index, line := range lines {
		if selected[index] {
			result = append(result, line)
		}
	}
	return result
}

func splitAddressLines(value string) []string {
	value = strings.NewReplacer("\r\n", "\n", "\r", "\n", ",", "\n").Replace(value)
	raw := strings.Split(value, "\n")
	lines := make([]string, 0, len(raw))
	for _, line := range raw {
		if normalized := normalizeSingleLine(line); normalized != "" {
			lines = append(lines, normalized)
		}
	}
	return lines
}

func decodeSignature(value string) ([]byte, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	const prefix = "data:image/png;base64,"
	if strings.HasPrefix(strings.ToLower(value), "data:") {
		if !strings.HasPrefix(strings.ToLower(value), prefix) {
			return nil, ErrInvalidSignature
		}
		value = value[len(prefix):]
	}
	if base64.StdEncoding.DecodedLen(len(value)) > maxSignatureBytes {
		return nil, ErrInvalidSignature
	}
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil || len(decoded) == 0 || len(decoded) > maxSignatureBytes {
		return nil, ErrInvalidSignature
	}
	configuration, err := png.DecodeConfig(bytes.NewReader(decoded))
	if err != nil || configuration.Width < 1 || configuration.Height < 1 ||
		configuration.Width > maxSignatureSide || configuration.Height > maxSignatureSide ||
		configuration.Width*configuration.Height > maxSignatureSide*maxSignatureSide {
		return nil, ErrInvalidSignature
	}
	return decoded, nil
}

func sanitizeFilename(value string) (string, error) {
	value = normalizeSingleLine(value)
	if value == "" {
		return "papierbrief.pdf", nil
	}
	if exceeds(value, maxFilenameRunes) || strings.ContainsAny(value, "\x00\r\n") {
		return "", ErrInvalidFilename
	}
	value = filepath.Base(strings.ReplaceAll(value, "\\", "/"))
	if strings.EqualFold(filepath.Ext(value), ".pdf") {
		value = value[:len(value)-len(filepath.Ext(value))]
	}
	value = strings.TrimSpace(value)
	var filename strings.Builder
	lastSeparator := false
	for _, r := range value {
		allowed := unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_'
		if allowed {
			filename.WriteRune(r)
			lastSeparator = false
			continue
		}
		if !lastSeparator {
			filename.WriteByte('_')
			lastSeparator = true
		}
	}
	base := strings.Trim(filename.String(), "._-")
	if base == "" || base == "." || base == ".." {
		return "", ErrInvalidFilename
	}
	baseRunes := []rune(base)
	if len(baseRunes) > 120 {
		base = string(baseRunes[:120])
	}
	return base + ".pdf", nil
}

func render(input normalizedInput) ([]byte, error) {
	regular, err := fontFiles.ReadFile("fonts/NotoSans-Regular.ttf")
	if err != nil {
		return nil, err
	}
	bold, err := fontFiles.ReadFile("fonts/NotoSans-Bold.ttf")
	if err != nil {
		return nil, err
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8FontFromBytes("NotoSans", "", regular)
	pdf.AddUTF8FontFromBytes("NotoSans", "B", bold)
	pdf.SetMargins(25, 20, 20)
	pdf.SetAutoPageBreak(true, 22)
	pdf.SetTitle(input.Subject, true)
	pdf.SetAuthor(input.Sender.Name, true)
	pdf.AliasNbPages("{nb}")
	pdf.SetHeaderFuncMode(func() {
		if pdf.PageNo() == 1 {
			return
		}
		pdf.SetXY(25, 18)
		pdf.SetFont("NotoSans", "B", 9)
		pdf.CellFormat(165, 5, input.Subject+" – Fortsetzung", "", 0, "L", false, 0, "")
		pdf.SetY(30)
	}, true)
	pdf.SetFooterFunc(func() {
		pdf.SetY(-14)
		pdf.SetFont("NotoSans", "", 8)
		pdf.CellFormat(0, 5, fmt.Sprintf("Seite %d von {nb}", pdf.PageNo()), "", 0, "C", false, 0, "")
	})

	pdf.AddPage()
	writeFirstPageHeader(pdf, input)
	pdf.SetFont("NotoSans", "", 10.5)
	pdf.MultiCell(165, 5.6, input.GermanText, "", "L", false)

	closingIncluded := containsClosing(input.GermanText)
	closingHeight := 24.0
	if len(input.Signature) > 0 {
		closingHeight += 18
	}
	ensureSpace(pdf, closingHeight)
	if !closingIncluded {
		pdf.Ln(5)
		pdf.SetFont("NotoSans", "", 10.5)
		pdf.Cell(0, 5.6, "Mit freundlichen Grüßen")
		pdf.Ln(7)
	}
	if len(input.Signature) > 0 {
		options := fpdf.ImageOptions{ImageType: "PNG", ReadDpi: false}
		info := pdf.RegisterImageOptionsReader("signature", options, bytes.NewReader(input.Signature))
		if info == nil || pdf.Error() != nil {
			return nil, fmt.Errorf("register signature: %w", pdf.Error())
		}
		width, height := info.Extent()
		scale := min(45/width, 18/height)
		width *= scale
		height *= scale
		pdf.ImageOptions("signature", 25, pdf.GetY(), width, height, false, options, 0, "")
		pdf.SetY(pdf.GetY() + height + 2)
	}
	if !closingIncluded || len(input.Signature) > 0 {
		pdf.SetFont("NotoSans", "", 10.5)
		pdf.Cell(0, 5.6, input.Sender.Name)
	}

	var output bytes.Buffer
	if err := pdf.Output(&output); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func writeFirstPageHeader(pdf *fpdf.Fpdf, input normalizedInput) {
	pdf.SetFont("NotoSans", "", 6.5)
	pdf.SetXY(23, 60.5)
	pdf.CellFormat(85, 2.5, strings.Join(input.Sender.Lines, " · "), "", 0, "L", false, 0, "")

	pdf.SetFont("NotoSans", "", 9.5)
	y := 63.5
	for _, line := range input.Recipient.Lines {
		pdf.SetXY(23, y)
		pdf.CellFormat(85, 3.6, line, "", 0, "L", false, 0, "")
		y += 3.6
	}

	pdf.SetFont("NotoSans", "", 9.5)
	pdf.SetXY(110, 63.5)
	pdf.CellFormat(75, 5, fmt.Sprintf("%s, %s", input.Sender.City, input.Now.Format("02.01.2006")), "", 0, "R", false, 0, "")
	if input.Reference != "" {
		pdf.SetXY(25, 108)
		pdf.CellFormat(165, 5, "Aktenzeichen: "+input.Reference, "", 0, "L", false, 0, "")
	}

	pdf.SetXY(25, 120)
	pdf.SetFont("NotoSans", "B", 11)
	pdf.MultiCell(165, 5.8, input.Subject, "", "L", false)
	pdf.Ln(7)
}

func ensureSpace(pdf *fpdf.Fpdf, height float64) {
	if pdf.GetY()+height > 275 {
		pdf.AddPage()
	}
}

func containsClosing(value string) bool {
	value = strings.ToLower(value)
	return strings.Contains(value, "mit freundlichen grüßen") ||
		strings.Contains(value, "mit freundlichen gruessen") ||
		strings.Contains(value, "freundliche grüße")
}

func normalizeText(value string) string {
	value = strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(value)
	lines := strings.Split(value, "\n")
	for index := range lines {
		lines[index] = strings.TrimRightFunc(lines[index], unicode.IsSpace)
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func normalizeSingleLine(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func fold(value string) string {
	return strings.ToLower(normalizeSingleLine(value))
}

func exceeds(value string, maximum int) bool {
	return !utf8.ValidString(value) || utf8.RuneCountInString(value) > maximum
}
