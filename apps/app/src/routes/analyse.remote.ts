import * as v from 'valibot';
import {
	vGenerateUploadUrlInputBody,
	vEnqueueDocumentAnalysisBody
} from '../generated/valibot.gen';
import {
	createCase,
	generateUploadUrl,
	confirmDocumentUpload,
	enqueueDocumentExtraction,
	enqueueDocumentAnalysis,
	getDocumentStatus
} from '../generated';
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
	...vEnqueueDocumentAnalysisBody.entries
});

const enqueueAnalysisCommand = authenticatedCommand(
	vAnalysisRequest,
	async ({ caseId, documentId, model, language }) => {
		const result = await enqueueDocumentAnalysis({
			...getClientConfig(),
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
	getDocumentStatusQuery
};
