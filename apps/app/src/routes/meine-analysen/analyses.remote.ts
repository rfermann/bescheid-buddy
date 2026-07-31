import { getCasesByUser } from '../../generated';
import { getClientConfig } from '$lib/server/apiClient';
import { authenticatedQuery } from '$lib/server/authenticatedRemotes';

export const getCasesQuery = authenticatedQuery('unchecked', async () => {
	const result = await getCasesByUser({
		...getClientConfig()
	});

	if (result.error) {
		throw new Error('Failed to fetch cases');
	}

	return result.data.cases ?? [];
});
