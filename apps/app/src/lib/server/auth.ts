import { betterAuth } from 'better-auth/minimal';
import { apiKey } from '@better-auth/api-key';
import { drizzleAdapter } from 'better-auth/adapters/drizzle';
import { sveltekitCookies } from 'better-auth/svelte-kit';
import { env } from '$env/dynamic/private';
import { getRequestEvent } from '$app/server';
import { db } from '$lib/server/db';
import { telegram } from 'better-auth-telegram';
import { admin, emailOTP, openAPI } from 'better-auth/plugins';
import { getOtpEmail, renderEmail } from '$lib/server/email-renderer';
import { m } from '$lib/paraglide/messages.js';
import { stripeClient } from '$lib/stripe';
import { stripe } from '@better-auth/stripe';
import { sendEmail } from '../../generated';
import { users } from '$lib/server/db/auth.schema';
import { sendEmail as sendEmail1 } from './mailer';
import { eq } from 'drizzle-orm';

export const auth = betterAuth({
	baseURL: env.ORIGIN,
	secret: env.BETTER_AUTH_SECRET,
	database: drizzleAdapter(db, { provider: 'pg' }),
	account: {
		modelName: 'accounts'
	},
	session: {
		modelName: 'sessions'
	},
	user: {
		modelName: 'users',
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
	verification: {
		modelName: 'verifications'
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
					const accessKey = await getAdminApiKey();
					if (!accessKey) {
						return;
					}

					const OtpEmailComponent = await getOtpEmail();
					const html = await renderEmail(OtpEmailComponent, { code: otp, mins: 5 });

					// intentionally not awaited
					sendEmail({
						headers: {
							'X-API-Key': accessKey
						},
						body: {
							html,
							receiver: email,
							sender: env.SMTP_USER,
							subject: m.email_loginCodeSubject()
						}
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
		apiKey({
			enableSessionForAPIKeys: true,
			schema: {
				apikey: {
					modelName: 'apikeys'
				}
			}
		}),
		openAPI(),
		stripe({
			stripeClient,
			stripeWebhookSecret: env.STRIPE_WEBHOOK_SECRET,
			createCustomerOnSignUp: true
		}),
		admin({
			defaultRole: 'admin'
		}),
		sveltekitCookies(getRequestEvent) // make sure this is the last plugin in the array
	]
});

const getAdminApiKey = async () => {
	const userId = await db
		.select({ id: users.id })
		.from(users)
		.where(eq(users.email, 'admin@admin.dev'))
		.limit(1);

	if (!userId.length) {
		return null;
	}

	const res = await auth.api.createApiKey({
		body: {
			name: 'key',
			userId: userId[0].id, // server-only
			prefix: 'ak_',
			remaining: 2
		}
	});

	return res.key;
};

export { getAdminApiKey };
