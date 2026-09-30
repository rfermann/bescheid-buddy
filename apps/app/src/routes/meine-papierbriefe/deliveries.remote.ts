import { error } from '@sveltejs/kit';
import * as v from 'valibot';
import { getPaperLetterDelivery, listPaperLetterDeliveries } from '../../generated';
import { client } from '../../generated/client.gen';
import { getClientConfig } from '$lib/server/apiClient';
import { authenticatedQuery } from '$lib/server/authenticatedRemotes';

export const getPaperLetterDeliveriesQuery = authenticatedQuery('unchecked', async () => {
	const result = await listPaperLetterDeliveries({ ...getClientConfig() });

	if (result.error || !result.data) {
		throw new Error('Failed to fetch paper letter deliveries');
	}

	return result.data.deliveries ?? [];
});

type DeliveryPDFURLs = { inlineUrl: string; downloadUrl: string };

export const getPaperLetterDeliveryDetailQuery = authenticatedQuery(
	v.object({ id: v.pipe(v.string(), v.uuid()) }),
	async ({ id }) => {
		const config = getClientConfig();
		const [delivery, pdfURLs] = await Promise.all([
			getPaperLetterDelivery({ ...config, path: { deliveryId: id } }),
			client.get<{ 200: DeliveryPDFURLs }>({
				...config,
				url: '/v1/paper-letters/deliveries/{deliveryId}/pdf-urls',
				path: { deliveryId: id }
			})
		]);

		if (delivery.response?.status === 404 || pdfURLs.response?.status === 404) {
			error(404, 'Paper letter delivery not found');
		}
		if (!delivery.data || !pdfURLs.data) {
			error(500, 'Failed to fetch paper letter delivery');
		}
		return { delivery: delivery.data, pdfURLs: pdfURLs.data };
	}
);
