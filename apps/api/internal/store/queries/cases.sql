-- name: CreateCase :one
INSERT INTO cases (user_id, title)
VALUES ($1, $2)
RETURNING *;

-- name: GetCasesByUserID :many
SELECT *
FROM cases
WHERE user_id = $1
ORDER BY created_at DESC, id DESC;

-- name: GetCaseByIDAndUserID :one
SELECT *
FROM cases
WHERE id = $1
  AND user_id = $2;

-- name: UpdateCaseTitle :exec
UPDATE cases
SET title = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: CreateCaseDocument :one
INSERT INTO case_documents (case_id, file_name)
VALUES ($1, $2)
RETURNING *;

-- name: GetCaseDocumentByID :one
SELECT *
FROM case_documents
WHERE id = $1;

-- name: GetInitialCaseDocument :one
SELECT *
FROM case_documents
WHERE case_id = $1
ORDER BY created_at, id
LIMIT 1;

-- name: GetCaseDocumentsForDetail :many
SELECT document.id AS document_id,
       document.file_name,
       document.upload_status,
       extraction.status AS extraction_status,
       analysis.status AS analysis_status,
       analysis.analysis_text,
       analysis.analysis_metadata,
       analysis.error_message AS analysis_error_message
FROM case_documents AS document
LEFT JOIN document_extractions AS extraction ON extraction.document_id = document.id
LEFT JOIN document_analysis_results AS analysis ON analysis.document_id = document.id
WHERE document.case_id = $1
ORDER BY document.created_at, document.id;

-- name: MarkCaseDocumentUploaded :one
UPDATE case_documents
SET upload_status = 'uploaded',
    uploaded_at   = NOW(),
    updated_at    = NOW()
WHERE id = $1
RETURNING *;

-- name: MarkCaseDocumentUploadFailed :one
UPDATE case_documents
SET upload_status = 'failed',
    updated_at    = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateDocumentExtraction :one
INSERT INTO document_extractions (document_id, job_id, status)
VALUES ($1, $2, 'pending')
RETURNING id, document_id, status, job_id, provider, mime, page_count, char_count, text_sha256, error_message, started_at, completed_at, created_at, updated_at;

-- name: GetDocumentExtractionByDocumentID :one
SELECT id,
       document_id,
       status,
       job_id,
       provider,
       mime,
       page_count,
       char_count,
       text_sha256,
       error_message,
       started_at,
       completed_at,
       created_at,
       updated_at
FROM document_extractions
WHERE document_id = $1;

-- name: MarkDocumentExtractionProcessing :one
UPDATE document_extractions
SET status     = 'processing',
    started_at = NOW(),
    updated_at = NOW()
WHERE document_id = $1
RETURNING id, document_id, status, job_id, provider, mime, page_count, char_count, text_sha256, error_message, started_at, completed_at, created_at, updated_at;

-- name: CompleteDocumentExtraction :one
UPDATE document_extractions
SET status       = 'completed',
    text         = $2,
    provider     = $3,
    mime         = $4,
    page_count   = $5,
    char_count   = $6,
    text_sha256  = $7,
    completed_at = NOW(),
    updated_at   = NOW()
WHERE document_id = $1
RETURNING id, document_id, status, job_id, provider, mime, page_count, char_count, text_sha256, error_message, started_at, completed_at, created_at, updated_at;

-- name: FailDocumentExtraction :one
UPDATE document_extractions
SET status        = 'failed',
    error_message = $2,
    completed_at  = NOW(),
    updated_at    = NOW()
WHERE document_id = $1
RETURNING id, document_id, status, job_id, provider, mime, page_count, char_count, text_sha256, error_message, started_at, completed_at, created_at, updated_at;

-- name: GetDocumentExtractionTextByDocumentID :one
SELECT status, text
FROM document_extractions
WHERE document_id = $1;

-- name: CreateDocumentAnalysis :one
INSERT INTO document_analysis_results (document_id, job_id, model, language, status)
VALUES ($1, $2, $3, $4, 'pending')
RETURNING id, document_id, status, job_id, model, language, provider, error_message, started_at, completed_at, created_at, updated_at;

-- name: GetDocumentAnalysisByDocumentID :one
SELECT id, document_id, status, job_id, model, language, detected_language, provider, analysis_text, analysis_metadata, error_message, started_at, completed_at, created_at, updated_at
FROM document_analysis_results
WHERE document_id = $1;

-- name: MarkDocumentAnalysisProcessing :one
UPDATE document_analysis_results
SET status     = 'processing',
    started_at = NOW(),
    updated_at = NOW()
WHERE document_id = $1
RETURNING id, document_id, status, job_id, model, provider, error_message, started_at, completed_at, created_at, updated_at;

-- name: CompleteDocumentAnalysis :one
UPDATE document_analysis_results
SET status            = 'completed',
    analysis_text     = $2,
    analysis_metadata = $3,
    provider          = $4,
    detected_language = $5,
    completed_at      = NOW(),
    updated_at        = NOW()
WHERE document_id = $1
RETURNING id, document_id, status, job_id, model, provider, error_message, started_at, completed_at, created_at, updated_at;

-- name: FailDocumentAnalysis :one
UPDATE document_analysis_results
SET status        = 'failed',
    error_message = $2,
    completed_at  = NOW(),
    updated_at    = NOW()
WHERE document_id = $1
RETURNING id, document_id, status, job_id, model, provider, error_message, started_at, completed_at, created_at, updated_at;

-- name: CreateDocumentReplyGeneration :one
INSERT INTO document_reply_generations (id, user_id, document_id, mode, client_request_id, job_id, model, provider, language, input)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetDocumentReplyGenerationByID :one
SELECT *
FROM document_reply_generations
WHERE id = $1;

-- name: GetDocumentReplyGenerationByIDAndUserID :one
SELECT *
FROM document_reply_generations
WHERE id = $1
  AND user_id = $2;

-- name: GetDocumentReplyGenerationByUserAndClientRequestID :one
SELECT *
FROM document_reply_generations
WHERE user_id = $1
  AND client_request_id = $2;

-- name: SetDocumentReplyGenerationJobID :exec
UPDATE document_reply_generations
SET job_id    = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: MarkDocumentReplyGenerationProcessing :one
UPDATE document_reply_generations
SET status     = 'processing',
    started_at = COALESCE(started_at, NOW()),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CompleteDocumentReplyGeneration :one
UPDATE document_reply_generations
SET status        = 'completed',
    result        = $2,
    error_message = NULL,
    completed_at  = NOW(),
    updated_at    = NOW()
WHERE id = $1
RETURNING *;

-- name: FailActiveDocumentReplyGeneration :one
UPDATE document_reply_generations
SET status        = 'failed',
    error_message = $2,
    completed_at  = NOW(),
    updated_at    = NOW()
WHERE id = $1
  AND status IN ('pending', 'processing')
RETURNING id, user_id;

-- name: CreatePaperLetterUpload :one
INSERT INTO paper_letter_uploads (id, user_id, sanitized_filename)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPaperLetterUploadByIDAndUserID :one
SELECT *
FROM paper_letter_uploads
WHERE id = $1
  AND user_id = $2;

-- name: ConfirmPaperLetterUpload :one
UPDATE paper_letter_uploads
SET status       = 'confirmed',
    size_bytes   = $3,
    page_count   = $4,
    sha256       = $5,
    confirmed_at = NOW(),
    updated_at   = NOW()
WHERE id = $1
  AND user_id = $2
  AND status = 'pending'
RETURNING *;

-- name: CreatePaperLetterDelivery :one
INSERT INTO paper_letter_deliveries
    (id, user_id, upload_id, client_request_id, case_id, recipient, topic, page_count, bp_cost,
     delivery_product, address_position, print_mode, print_spectrum)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetPaperLetterDeliveryByID :one
SELECT * FROM paper_letter_deliveries WHERE id = $1;

-- name: GetPaperLetterDeliveryByUserAndClientRequestID :one
SELECT * FROM paper_letter_deliveries WHERE user_id = $1 AND client_request_id = $2;

-- name: GetPaperLetterDeliveriesByUserID :many
SELECT *
FROM paper_letter_deliveries
WHERE user_id = $1
ORDER BY created_at DESC, id DESC;

-- name: SetPaperLetterDeliveryJobID :exec
UPDATE paper_letter_deliveries SET job_id = $2, updated_at = NOW() WHERE id = $1;

-- name: MarkPaperLetterDeliveryProcessing :one
UPDATE paper_letter_deliveries
SET status = 'processing', provider_submission_started_at = COALESCE(provider_submission_started_at, NOW()), updated_at = NOW()
WHERE id = $1 AND status IN ('pending', 'processing')
RETURNING *;

-- name: MarkPaperLetterDeliverySent :one
UPDATE paper_letter_deliveries
SET status = 'sent', pingen_id = $2, error_message = NULL, sent_at = NOW(), updated_at = NOW()
WHERE id = $1 AND status = 'processing'
RETURNING *;

-- name: MarkPaperLetterDeliveryManualReview :one
UPDATE paper_letter_deliveries
SET status = 'manual_review', updated_at = NOW()
WHERE id = $1
  AND status = 'processing'
  AND provider_submission_started_at <= sqlc.arg(cutoff)
RETURNING *;

-- name: FailPaperLetterDelivery :one
UPDATE paper_letter_deliveries
SET status = 'failed', error_message = $2, updated_at = NOW()
WHERE id = $1 AND status = 'processing'
RETURNING id;
