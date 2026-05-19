import { betterAuth } from 'better-auth/minimal';
import { drizzleAdapter } from 'better-auth/adapters/drizzle';
import { sveltekitCookies } from 'better-auth/svelte-kit';
import { env } from '$env/dynamic/private';
import { getRequestEvent } from '$app/server';
import { db } from '$lib/server/db';
import { telegram } from 'better-auth-telegram';
import { emailOTP, openAPI } from 'better-auth/plugins';
import { getOtpEmail, renderEmail } from '$lib/server/email-renderer';
import { sendEmail } from '$lib/server/mailer';
import { m } from '$lib/paraglide/messages.js';
import { stripeClient } from '$lib/stripe';
import { stripe } from '@better-auth/stripe';

export const auth = betterAuth({
	baseURL: env.ORIGIN,
	secret: env.BETTER_AUTH_SECRET,
	database: drizzleAdapter(db, { provider: 'pg' }),
	user: {
		additionalFields: {
			monthlyPoints: {
				input: false,
				defaultValue: 5,
				required: true,
				type: 'number'
			},
			lifetimePoints: {
				input: false,
				defaultValue: 0,
				required: true,
				type: 'number'
			}
		}
	},
	advanced: {
		database: {
			generateId: 'uuid'
		},
		crossSubDomainCookies: {
			enabled: !!env.CROSS_DOMAIN_ORIGIN,
			domain: env.CROSS_DOMAIN_ORIGIN
		}
	},
	plugins: [
		emailOTP({
			sendVerificationOTP: async ({ email, otp, type }) => {
				if (type === 'sign-in') {
					const OtpEmailComponent = await getOtpEmail();
					const html = await renderEmail(OtpEmailComponent, { code: otp, mins: 5 });
					await sendEmail({
						to: email,
						subject: m.email_loginCodeSubject(),
						html
					});
				}
			}
		}),
		telegram({
			allowUserToLink: false,
			botToken: env.TELEGRAM_BOT_TOKEN,
			botUsername: env.TELEGRAM_BOT_USERNAME,
			mapTelegramDataToUser: (data) => ({
				name: data.username || data.first_name,
				image: data.photo_url,
				email: 'undefined'
			}),
			miniApp: {
				enabled: true,
				validateInitData: false,
				allowAutoSignin: true,
				mapMiniAppDataToUser: (user) => ({
					name: `${user.first_name} ${user.last_name || ''}`.trim(),
					image: user.photo_url,
					email: user.id.toString()
				})
			}
		}),
		openAPI(),
		stripe({
			stripeClient,
			stripeWebhookSecret: env.STRIPE_WEBHOOK_SECRET,
			createCustomerOnSignUp: true
		}),
		sveltekitCookies(getRequestEvent) // make sure this is the last plugin in the array
	]
});
