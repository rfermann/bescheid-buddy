import * as v from 'valibot';
import {
	confirmPaperLetterUpload,
	createPaperLetterDelivery,
	createPaperLetterGeneration,
	createPaperLetterPdf,
	createPaperLetterUpload,
	getPaperLetterDelivery,
	getPaperLetterGeneration
} from '../../generated';
import {
	vConfirmPaperLetterUploadPath,
	vCreatePaperLetterDeliveryBody,
	vCreatePaperLetterUploadBody,
	vGetPaperLetterDeliveryPath
} from '../../generated/valibot.gen';
import { getClientConfig } from '$lib/server/apiClient';
import { authenticatedCommand, authenticatedQuery } from '$lib/server/authenticatedRemotes';

const category = v.picklist([
	'Jobcenter',
	'Krankenkasse',
	'Finanzamt',
	'Ausländerbehörde',
	'Rentenversicherung',
	'Arbeitsagentur',
	'Inkasso',
	'Gericht',
	'Schule/Kita',
	'Polizei',
	'Versicherung',
	'Bank',
	'Telekom/Internet',
	'Energie',
	'Other',
	'Universal'
]);

const subcategory = v.picklist([
	'widerspruch',
	'fristverlaengerung',
	'unterlagen_nachreichen',
	'zahlung_klaeren',
	'termin_verschieben',
	'kuendigung',
	'freie_formulierung'
]);

const party = v.object({
	name: v.pipe(v.string(), v.maxLength(200)),
	address: v.pipe(v.string(), v.maxLength(1000))
});

const createGenerationInput = v.object({
	mode: v.picklist(['generate', 'improve']),
	clientRequestId: v.pipe(v.string(), v.uuid()),
	model: v.string(),
	locale: v.picklist(['en', 'de', 'ru', 'ua', 'pl', 'ar']),
	sourceCaseId: v.optional(v.pipe(v.string(), v.uuid())),
	sourceDocumentId: v.optional(v.pipe(v.string(), v.uuid())),
	input: v.object({
		sender: party,
		recipient: party,
		category,
		subcategory,
		subject: v.pipe(v.string(), v.maxLength(500)),
		reference: v.pipe(v.string(), v.maxLength(200)),
		instruction: v.pipe(v.string(), v.maxLength(2500)),
		currentGermanLetter: v.optional(v.pipe(v.string(), v.maxLength(24000))),
		feedback: v.optional(v.pipe(v.string(), v.maxLength(2500)))
	})
});

export type CreateGenerationRequest = v.InferInput<typeof createGenerationInput>;

export const createGenerationCommand = authenticatedCommand(
	createGenerationInput,
	async (request) => {
		const result = await createPaperLetterGeneration({
			...getClientConfig(),
			body: {
				mode: request.mode,
				clientRequestId: request.clientRequestId,
				model: request.model,
				locale: request.locale,
				sourceCaseId: request.sourceCaseId,
				sourceDocumentId: request.sourceDocumentId,
				input: request.input
			}
		});
		if (result.data) {
			return { ok: true as const, generation: result.data };
		}

		const status = result.response?.status;
		const code =
			status === 401
				? 'unauthorized'
				: status === 402
					? 'insufficient_points'
					: status === 404
						? 'invalid_source'
						: status === 409
							? 'source_not_ready'
							: status === 422
								? 'invalid_input'
								: 'request_failed';
		return { ok: false as const, error: code };
	}
);

const generationStatusInput = v.object({
	generationId: v.pipe(v.string(), v.uuid())
});

export const getGenerationStatusQuery = authenticatedQuery(
	generationStatusInput,
	async ({ generationId }) => {
		const result = await getPaperLetterGeneration({
			...getClientConfig(),
			path: { generationId }
		});
		if (!result.data) {
			throw new Error('Failed to fetch paper letter generation');
		}
		return result.data;
	}
);

const createPDFInput = v.object({
	germanText: v.pipe(v.string(), v.maxLength(24000)),
	subject: v.pipe(v.string(), v.maxLength(500)),
	sender: party,
	recipient: party,
	reference: v.pipe(v.string(), v.maxLength(200)),
	filename: v.pipe(v.string(), v.maxLength(200)),
	signaturePngBase64: v.optional(v.pipe(v.string(), v.maxLength(700000)))
});

export type CreatePDFRequest = v.InferInput<typeof createPDFInput>;

const responseFilename = (header: string | null, fallback: string) => {
	if (!header) {
		return fallback;
	}
	const encoded = /filename\*=UTF-8''([^;]+)/i.exec(header)?.[1];
	if (encoded) {
		try {
			return decodeURIComponent(encoded);
		} catch {
			return fallback;
		}
	}
	return /filename="?([^";]+)"?/i.exec(header)?.[1] ?? fallback;
};

export const createPDFCommand = authenticatedCommand(createPDFInput, async (request) => {
	const result = await createPaperLetterPdf({
		...getClientConfig(),
		body: request,
		parseAs: 'arrayBuffer'
	});
	if (result.data) {
		const buffer = result.data as unknown as ArrayBuffer;
		const headers = result.response?.headers;
		return {
			ok: true as const,
			bytes: new Uint8Array(buffer),
			filename: responseFilename(
				headers?.get('Content-Disposition') ?? null,
				request.filename || 'papierbrief.pdf'
			),
			pageCount: Number.parseInt(headers?.get('X-Page-Count') ?? '', 10) || 1
		};
	}

	return {
		ok: false as const,
		error:
			result.response?.status === 401
				? ('unauthorized' as const)
				: result.response?.status === 422
					? ('invalid_input' as const)
					: ('request_failed' as const)
	};
});

export const createPaperLetterUploadCommand = authenticatedCommand(
	vCreatePaperLetterUploadBody,
	async (request) => {
		const result = await createPaperLetterUpload({
			...getClientConfig(),
			body: request
		});
		if (result.data) {
			return { ok: true as const, ...result.data };
		}
		return {
			ok: false as const,
			error:
				result.response?.status === 401
					? ('unauthorized' as const)
					: result.response?.status === 422
						? ('invalid_filename' as const)
						: ('request_failed' as const)
		};
	}
);

export const confirmPaperLetterUploadCommand = authenticatedCommand(
	vConfirmPaperLetterUploadPath,
	async ({ uploadId }) => {
		const result = await confirmPaperLetterUpload({
			...getClientConfig(),
			path: { uploadId }
		});
		if (result.data) {
			return {
				ok: true as const,
				uploadId: result.data.uploadId,
				status: result.data.status,
				sizeBytes: Number(result.data.sizeBytes),
				pageCount: result.data.pageCount,
				deliveryPrice: Number(result.data.deliveryPrice)
			};
		}
		const status = result.response?.status;
		return {
			ok: false as const,
			error:
				status === 401
					? ('unauthorized' as const)
					: status === 400
						? ('missing_object' as const)
						: status === 404
							? ('upload_not_found' as const)
							: status === 413
								? ('too_large' as const)
								: status === 422
									? ('invalid_pdf' as const)
									: ('request_failed' as const)
		};
	}
);
export const createPaperLetterDeliveryCommand = authenticatedCommand(
	vCreatePaperLetterDeliveryBody,
	async (request) => {
		const result = await createPaperLetterDelivery({ ...getClientConfig(), body: request });

		if (result.data) {
			return { ok: true as const, delivery: result.data };
		}
		const status = result.response?.status;
		return {
			ok: false as const,
			error:
				status === 401
					? ('unauthorized' as const)
					: status === 402
						? ('insufficient_points' as const)
						: status === 404
							? ('upload_not_found' as const)
							: status === 409
								? ('upload_not_confirmed' as const)
								: status === 422
									? ('invalid_input' as const)
									: ('request_failed' as const)
		};
	}
);

export const getPaperLetterDeliveryQuery = authenticatedQuery(
	vGetPaperLetterDeliveryPath,
	async ({ deliveryId }) => {
		const result = await getPaperLetterDelivery({ ...getClientConfig(), path: { deliveryId } });
		if (!result.data) {
			throw new Error('Failed to fetch paper letter delivery');
		}
		return result.data;
	}
);
