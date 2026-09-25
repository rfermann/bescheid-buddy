package paper_letters_handler

import (
	"encoding/json"
	"time"

	"github.com/bored-bots/bored-template/internal/services/paper_letters_service"
	"github.com/bored-bots/bored-template/internal/services/reply_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
)

type DeliveryBody struct {
	DeliveryID        uuid.UUID  `json:"deliveryId"`
	UploadID          uuid.UUID  `json:"uploadId"`
	ClientRequestID   uuid.UUID  `json:"clientRequestId"`
	CaseID            *uuid.UUID `json:"caseId" nullable:"true"`
	Recipient         *string    `json:"recipient" nullable:"true"`
	Topic             *string    `json:"topic" nullable:"true"`
	Status            string     `json:"status"`
	PageCount         int32      `json:"pageCount"`
	BuddyPointsCost   int32      `json:"buddyPointsCost"`
	DeliveryProduct   string     `json:"deliveryProduct"`
	AddressPosition   string     `json:"addressPosition"`
	PrintMode         string     `json:"printMode"`
	PrintSpectrum     string     `json:"printSpectrum"`
	PingenID          *string    `json:"pingenId,omitempty"`
	SubmissionStarted *time.Time `json:"submissionStartedAt,omitempty"`
	SentAt            *time.Time `json:"sentAt,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

func deliveryBody(delivery *db.PaperLetterDelivery) DeliveryBody {
	body := DeliveryBody{
		DeliveryID: delivery.ID, UploadID: delivery.UploadID, ClientRequestID: delivery.ClientRequestID,
		Status: string(delivery.Status), PageCount: delivery.PageCount, BuddyPointsCost: delivery.BpCost,
		DeliveryProduct: delivery.DeliveryProduct, AddressPosition: delivery.AddressPosition,
		PrintMode: delivery.PrintMode, PrintSpectrum: delivery.PrintSpectrum,
		SubmissionStarted: delivery.ProviderSubmissionStartedAt, SentAt: delivery.SentAt,
		CreatedAt: delivery.CreatedAt, UpdatedAt: delivery.UpdatedAt,
	}
	if delivery.PingenID.Valid {
		body.PingenID = &delivery.PingenID.String
	}
	if delivery.CaseID.Valid {
		caseID := uuid.UUID(delivery.CaseID.Bytes)
		body.CaseID = &caseID
	}
	if delivery.Recipient.Valid {
		body.Recipient = &delivery.Recipient.String
	}
	if delivery.Topic.Valid {
		body.Topic = &delivery.Topic.String
	}
	return body
}

type PartyInputBody struct {
	Name    string `json:"name,omitempty" maxLength:"200"`
	Address string `json:"address,omitempty" maxLength:"1000"`
}

type GenerationInputBody struct {
	Sender              *PartyInputBody `json:"sender,omitempty"`
	Recipient           *PartyInputBody `json:"recipient,omitempty"`
	Category            string          `json:"category,omitempty" maxLength:"100" enum:"Jobcenter,Krankenkasse,Finanzamt,Ausländerbehörde,Rentenversicherung,Arbeitsagentur,Inkasso,Gericht,Schule/Kita,Polizei,Versicherung,Bank,Telekom/Internet,Energie,Other,Universal"`
	Subcategory         string          `json:"subcategory,omitempty" maxLength:"100" enum:"widerspruch,fristverlaengerung,unterlagen_nachreichen,zahlung_klaeren,termin_verschieben,kuendigung,freie_formulierung"`
	Subject             string          `json:"subject,omitempty" maxLength:"500"`
	Reference           string          `json:"reference,omitempty" maxLength:"200"`
	Instruction         string          `json:"instruction,omitempty" maxLength:"2500"`
	CurrentGermanLetter string          `json:"currentGermanLetter,omitempty" maxLength:"24000"`
	Feedback            string          `json:"feedback,omitempty" maxLength:"2500"`
}

type GenerationResultBody struct {
	GermanLetter       string `json:"germanLetter"`
	UserLanguageLetter string `json:"userLanguageLetter"`
	Explanation        string `json:"explanation"`
	Subject            string `json:"subject"`
	Tone               string `json:"tone"`
}

type GenerationBody struct {
	GenerationID     uuid.UUID             `json:"generationId"`
	ClientRequestID  uuid.UUID             `json:"clientRequestId"`
	SourceCaseID     *uuid.UUID            `json:"sourceCaseId" nullable:"true"`
	SourceDocumentID *uuid.UUID            `json:"sourceDocumentId" nullable:"true"`
	Mode             string                `json:"mode"`
	Status           string                `json:"status"`
	Model            string                `json:"model"`
	Locale           string                `json:"locale"`
	Result           *GenerationResultBody `json:"result,omitempty"`
	ErrorMessage     string                `json:"errorMessage,omitempty"`
	Refunded         bool                  `json:"refunded"`
	StartedAt        *time.Time            `json:"startedAt,omitempty"`
	CompletedAt      *time.Time            `json:"completedAt,omitempty"`
	CreatedAt        time.Time             `json:"createdAt"`
	UpdatedAt        time.Time             `json:"updatedAt"`
}

func generationBody(result *paper_letters_service.GenerationResult) (GenerationBody, error) {
	generation := result.Generation
	body := GenerationBody{
		GenerationID:     generation.ID,
		ClientRequestID:  generation.ClientRequestID,
		SourceCaseID:     result.SourceCaseID,
		SourceDocumentID: result.SourceDocumentID,
		Mode:             string(generation.Mode),
		Status:           string(generation.Status),
		Model:            generation.Model,
		Locale:           generation.Language,
		ErrorMessage:     publicGenerationError(generation.ErrorMessage.String),
		Refunded:         generation.Status == "failed",
		StartedAt:        generation.StartedAt,
		CompletedAt:      generation.CompletedAt,
		CreatedAt:        generation.CreatedAt,
		UpdatedAt:        generation.UpdatedAt,
	}
	if len(generation.Result) > 0 {
		var generated reply_service.Result
		if err := json.Unmarshal(generation.Result, &generated); err != nil {
			return GenerationBody{}, err
		}
		body.Result = &GenerationResultBody{
			GermanLetter:       generated.GermanLetter,
			UserLanguageLetter: generated.UserLanguageLetter,
			Explanation:        generated.Explanation,
			Subject:            generated.Subject,
			Tone:               generated.Tone,
		}
	}
	return body, nil
}

func publicGenerationError(value string) string {
	switch value {
	case "invalid_generation_input", "source_not_ready", "unsupported_model", "invalid_generation_output", "generation_failed":
		return value
	case "":
		return ""
	default:
		return "generation_failed"
	}
}
