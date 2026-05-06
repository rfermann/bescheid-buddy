import { betterAuth } from 'better-auth/minimal';
import { drizzleAdapter } from 'better-auth/adapters/drizzle';
import { sveltekitCookies } from 'better-auth/svelte-kit';
import { env } from '$env/dynamic/private';
import { getRequestEvent } from '$app/server';
import { db } from '$lib/server/db';
import { telegram } from 'better-auth-telegram';
import { emailOTP, openAPI } from 'better-auth/plugins';

console.log(env.ORIGIN);

export const auth = betterAuth({
	baseURL: env.ORIGIN,
	secret: env.BETTER_AUTH_SECRET,
	database: drizzleAdapter(db, { provider: 'pg' }),
	advanced: {},
	plugins: [
		emailOTP({
			sendVerificationOTP: async ({ email, otp, type }) => {
				console.log(email, otp, type);
				console.log('sendVerificationOTP');
			}
		}),
		telegram({
			allowUserToLink: false,
			botToken: '7486765220:AAG97zoOx0MseGkrlOewaE-5vsPAClBbwno',
			botUsername: 'countdown2cuddle_bot',
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
		sveltekitCookies(getRequestEvent) // make sure this is the last plugin in the array
	]
});
