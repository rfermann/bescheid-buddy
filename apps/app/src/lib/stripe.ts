import { Stripe } from 'stripe';
import { env } from '$env/dynamic/private';

export const stripeClient = new Stripe(env.STRIPE_SECRET_KEY, {
	apiVersion: '2026-04-22.dahlia'
});
