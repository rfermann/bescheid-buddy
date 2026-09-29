import { listPaperLetterDeliveries } from '../../generated';
import { getClientConfig } from '$lib/server/apiClient';
import { authenticatedQuery } from '$lib/server/authenticatedRemotes';

export const getPaperLetterDeliveriesQuery = authenticatedQuery('unchecked', async () => {
	const result = await listPaperLetterDeliveries({ ...getClientConfig() });

	if (result.error || !result.data) {
		throw new Error('Failed to fetch paper letter deliveries');
	}

	return result.data.deliveries ?? [];
});
