import { error } from '@sveltejs/kit';
import * as v from 'valibot';
import {
	getCasesByUser,
	type CaseListItem,
	type GetDocumentStatusOutputBody
} from '../../generated';
import { client } from '../../generated/client.gen';
import { getClientConfig } from '$lib/server/apiClient';
import { authenticatedQuery } from '$lib/server/authenticatedRemotes';

const getCases = async () => {
	const result = await getCasesByUser({
		...getClientConfig()
	});

	if (result.error || !result.data) {
		throw new Error('Failed to fetch cases');
	}

	return result.data.cases ?? [];
};

export const getCasesQuery = authenticatedQuery('unchecked', getCases);

type CaseDetailResponse = Omit<CaseListItem, 'createdAt' | 'updatedAt'> & {
	createdAt: string;
	updatedAt: string;
	documents: Array<
		Pick<
			GetDocumentStatusOutputBody,
			'documentId' | 'fileName' | 'uploadStatus' | 'extractionStatus'
		> & {
			analysis: Pick<
				GetDocumentStatusOutputBody['analysis'],
				'status' | 'text' | 'metadata' | 'errorMessage'
			>;
		}
	>;
};

const getCase = async (id: string) => {
	const result = await client.get<{ 200: CaseDetailResponse }>({
		...getClientConfig(),
		url: '/v1/cases/{caseId}',
		path: { caseId: id }
	});

	if (result.response?.status === 404) {
		error(404, 'Analysis not found');
	}
	if (result.error || !result.data) {
		error(500, 'Failed to fetch analysis');
	}

	return {
		...result.data,
		createdAt: new Date(result.data.createdAt),
		updatedAt: new Date(result.data.updatedAt)
	};
};

export const getAnalysisQuery = authenticatedQuery(
	v.object({ id: v.pipe(v.string(), v.uuid()) }),
	async ({ id }) => {
		const analysisCase = await getCase(id);
		return { analysisCase, documents: analysisCase.documents };
	}
);
