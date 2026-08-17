import * as v from 'valibot';
import { query } from '$app/server';
import { vGenerateUploadUrlInputBody } from '../generated/valibot.gen';
import {
	createCase,
	generateUploadUrl,
	confirmDocumentUpload,
	enqueueDocumentExtraction,
	getDocumentStatus
} from '../generated';
import { client } from '../generated/client.gen';
import { getClientConfig } from '$lib/server/apiClient';
import { authenticatedCommand, authenticatedQuery } from '$lib/server/authenticatedRemotes';

const vDocumentRef = v.object({
	caseId: v.pipe(v.string(), v.uuid()),
	documentId: v.pipe(v.string(), v.uuid())
});

const createUploadUrlCommand = authenticatedCommand(
	vGenerateUploadUrlInputBody,
	async ({ fileName }) => {
		const caseCreationResult = await createCase({
			...getClientConfig(),
			body: {
				title: fileName
			}
		});

		const caseId = caseCreationResult.data?.id;
		if (!caseId) {
			throw new Error('Case creation failed');
		}

		const generateUploadUrlResult = await generateUploadUrl({
			...getClientConfig(),
			path: {
				caseId
			},
			body: {
				fileName
			}
		});

		const uploadUrl = generateUploadUrlResult.data?.uploadUrl;
		const documentId = generateUploadUrlResult.data?.documentId;
		if (!uploadUrl || !documentId) {
			throw new Error('Failed to generate upload URL');
		}

		return {
			caseId,
			documentId,
			uploadUrl
		};
	}
);

const confirmUploadCommand = authenticatedCommand(vDocumentRef, async ({ caseId, documentId }) => {
	const result = await confirmDocumentUpload({
		...getClientConfig(),
		path: {
			caseId,
			documentId
		}
	});

	const uploadStatus = result.data?.uploadStatus;
	if (uploadStatus !== 'uploaded') {
		throw new Error('Upload confirmation failed');
	}

	return {
		uploadStatus
	};
});

const enqueueExtractionCommand = authenticatedCommand(
	vDocumentRef,
	async ({ caseId, documentId }) => {
		const result = await enqueueDocumentExtraction({
			...getClientConfig(),
			path: {
				caseId,
				documentId
			}
		});

		if (!result.data?.documentId) {
			throw new Error('Failed to start extraction');
		}

		return {
			documentId: result.data.documentId,
			jobId: result.data.jobId?.toString() ?? null
		};
	}
);

const vAnalysisRequest = v.object({
	...vDocumentRef.entries,
	model: v.string(),
	language: v.optional(v.string())
});

const enqueueAnalysisCommand = authenticatedCommand(
	vAnalysisRequest,
	async ({ caseId, documentId, model, language }) => {
		const result = await client.post<{ 200: { documentId: string; jobId?: bigint } }>({
			...getClientConfig(),
			url: '/v1/cases/{caseId}/documents/{documentId}/analyse',
			path: {
				caseId,
				documentId
			},
			body: {
				model,
				language
			}
		});

		if (!result.data?.documentId) {
			throw new Error('Failed to start analysis');
		}

		return {
			documentId: result.data.documentId,
			jobId: result.data.jobId?.toString() ?? null
		};
	}
);

type ModelCatalog = {
	models: Array<{ id: string; name: string; provider: string }>;
};

const getModelsQuery = query('unchecked', async () => {
	const result = await client.get<{ 200: ModelCatalog }>({
		...getClientConfig(),
		url: '/v1/models'
	});

	if (result.error || !result.data) {
		throw new Error('Failed to fetch models');
	}

	return result.data.models;
});

const getDocumentStatusQuery = authenticatedQuery(vDocumentRef, async ({ caseId, documentId }) => {
	const result = await getDocumentStatus({
		...getClientConfig(),
		path: {
			caseId,
			documentId
		}
	});

	if (!result.data) {
		throw new Error('Failed to fetch document status');
	}

	return result.data;
});

export {
	createUploadUrlCommand,
	confirmUploadCommand,
	enqueueExtractionCommand,
	enqueueAnalysisCommand,
	getModelsQuery,
	getDocumentStatusQuery
};
