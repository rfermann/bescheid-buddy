import { getLocale } from '$lib/paraglide/runtime';
import {
	createUploadUrlCommand,
	confirmUploadCommand,
	enqueueExtractionCommand,
	enqueueAnalysisCommand,
	getDocumentStatusQuery
} from './analyse.remote';
import type { Metadata } from '../generated';

export type AnalysisModel = 'gpt-4o-mini' | 'gemini-2.0-flash' | 'claude-haiku-4-5';
export type AnalysisStep = 'preparing' | 'uploading' | 'extracting' | 'analyzing' | 'finished';

export type AnalysisPipelineResult = {
	text?: string;
	metadata?: Metadata;
};

const POLL_INTERVAL_MS = 1500;
const PIPELINE_TIMEOUT_MS = 3 * 60 * 1000;

const wait = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export const runAnalysisPipeline = async ({
	file,
	model,
	onStep
}: {
	file: File;
	model: AnalysisModel;
	onStep: (step: AnalysisStep) => void;
}): Promise<AnalysisPipelineResult> => {
	onStep('preparing');
	const { caseId, documentId, uploadUrl } = await createUploadUrlCommand({
		fileName: file.name
	});

	onStep('uploading');
	const uploadResponse = await fetch(uploadUrl, {
		method: 'PUT',
		body: file,
		headers: {
			'Content-Type': file.type || 'application/octet-stream'
		}
	});
	if (!uploadResponse.ok) {
		throw new Error('Storage upload failed');
	}

	await confirmUploadCommand({ caseId, documentId });

	onStep('extracting');
	await enqueueExtractionCommand({ caseId, documentId });

	const statusQuery = getDocumentStatusQuery({ caseId, documentId });
	const deadline = Date.now() + PIPELINE_TIMEOUT_MS;
	let analysisEnqueued = false;

	while (true) {
		if (Date.now() > deadline) {
			throw new Error('Analysis timed out');
		}

		await statusQuery.refresh();
		const status = statusQuery.current;

		if (status) {
			if (status.extractionStatus === 'failed') {
				throw new Error('Extraction failed');
			}
			if (status.extractionStatus === 'completed' && !analysisEnqueued) {
				await enqueueAnalysisCommand({
					caseId,
					documentId,
					model,
					language: getLocale()
				});
				analysisEnqueued = true;
				onStep('analyzing');
			}

			if (analysisEnqueued) {
				if (status.analysis.status === 'failed') {
					throw new Error('Analysis failed');
				}
				if (status.analysis.status === 'completed') {
					onStep('finished');
					return {
						text: status.analysis.text,
						metadata: status.analysis.metadata
					};
				}
			}
		}

		await wait(POLL_INTERVAL_MS);
	}
};
