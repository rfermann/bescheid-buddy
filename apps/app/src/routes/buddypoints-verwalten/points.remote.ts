import { getRequestEvent, query } from '$app/server';
// import { stripeClient } from '$lib/server/auth';
import { getClientConfig } from '$lib/server/apiClient';
import { getAllBuddyPackages, getBuddyPoints, getBuddyPointsLedger } from '../../generated';
import * as v from 'valibot';
import { stripeClient } from '$lib/stripe';
import { env } from '$env/dynamic/private';
import { auth } from '$lib/server/auth';
import { db } from '$lib/server/db';
import { user } from '$lib/server/db/auth.schema';
import { eq } from 'drizzle-orm';

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

const getCheckoutSessionUrlInput = v.strictObject({
	name: v.string(),
	stripePriceId: v.string()
});

export const getCheckoutSessionUrlQuery = query(
	getCheckoutSessionUrlInput,
	async ({ name, stripePriceId }) => {
		const event = getRequestEvent();
		const session = await auth.api.getSession({
			headers: event.request.headers
		});

		if (!session || !session.user) {
			throw new Error('User not authenticated');
		}

		let stripeCustomerId = session.user.stripeCustomerId;

		if (!stripeCustomerId) {
			const createCustomerResponse = await stripeClient.customers.create({
				name: session.user.name

				// email: session.user.email
			});
			stripeCustomerId = createCustomerResponse.id;
			await db.update(user).set({ stripeCustomerId }).where(eq(user.id, session.user.id));
		}

		const stripeSession = await stripeClient.checkout.sessions.create({
			customer: stripeCustomerId,
			line_items: [{ price: stripePriceId, quantity: 1 }],
			mode: 'payment',

			success_url: `${env.ORIGIN}/buddypoints-verwalten?session_id={CHECKOUT_SESSION_ID}&name=${name}`,
			cancel_url: `${env.ORIGIN}/buddypoints-verwalten?canceled=true`
		});
		return stripeSession.url;
	}
);
