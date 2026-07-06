import { form, query } from '$app/server';
import { getClientConfig } from '$lib/server/apiClient';
import {
	getAllBuddyPackages,
	getBuddyPoints,
	getBuddyPointsLedger,
	createCheckoutSession
} from '../../generated';
import { vCreateCheckoutSessionBody } from '../../generated/valibot.gen.ts';
import { redirect } from '@sveltejs/kit';

export const getBuddyPackagesQuery = query('unchecked', async () => {
	const res = await getAllBuddyPackages({
		...getClientConfig()
	});

	if (res.error) {
		throw new Error(`Failed to fetch buddy packages: ${res.error}`);
	}

	return res.data;
});

export const getBuddyPointsQuery = query('unchecked', async () => {
	const res = await getBuddyPoints({
		...getClientConfig()
	});

	if (res.error) {
		throw new Error(`Failed to fetch buddy points: ${res.error}`);
	}

	return res.data;
});

export const getBuddyPointsLedgerQuery = query('unchecked', async () => {
	const res = await getBuddyPointsLedger({
		...getClientConfig()
	});

	if (res.error) {
		throw new Error(`Failed to fetch buddy points ledger: ${res.error}`);
	}

	return res.data;
});

export const createCheckoutSessionUrlForm = form(
	vCreateCheckoutSessionBody,
	async ({ package_name, stripe_price_id }) => {
		const res = await createCheckoutSession({
			...getClientConfig(),
			body: {
				package_name,
				stripe_price_id
			}
		});

		if (res.error) {
			throw new Error(`Failed to create checkout session: ${res.error}`);
		}

		redirect(303, res.data.url);
	}
);
