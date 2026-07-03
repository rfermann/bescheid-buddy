import { form, getRequestEvent } from '$app/server';
import { auth } from '$lib/server/auth';
import { invalid, redirect } from '@sveltejs/kit';
import * as v from 'valibot';

const getMessage = (error: unknown, fallback: string) => {
	if (error instanceof Error && error.message) {
		return error.message;
	}

	if (typeof error === 'object' && error !== null && 'message' in error) {
		const message = (error as { message?: unknown }).message;
		if (typeof message === 'string' && message) {
			return message;
		}
	}

	return fallback;
};

const sendOTPSchema = v.object({
	email: v.pipe(v.string(), v.email(), v.maxLength(255))
});

const verifyOTPSchema = v.object({
	email: v.pipe(v.string(), v.email(), v.maxLength(255)),
	otp: v.pipe(v.string(), v.minLength(6), v.maxLength(6), v.regex(/^\d{6}$/))
});

export const sendOTP = form(sendOTPSchema, async ({ email }) => {
	try {
		await auth.api.sendVerificationOTP({
			body: {
				email: String(email ?? ''),
				type: 'sign-in'
			}
		});

		return {
			success: true,
			message: 'We sent a one-time password to your email address.'
		};
	} catch (error) {
		return {
			success: false,
			message: getMessage(error, 'Unable to send the one-time password.')
		};
	}
});

export const verifyOTP = form(verifyOTPSchema, async ({ email, otp }) => {
	console.log('verifyOTP called with:', { email, otp });
	try {
		await auth.api.signInEmailOTP({
			body: {
				email: String(email ?? ''),
				otp: String(otp ?? '')
			}
		});
	} catch (error) {
		return invalid(getMessage(error, 'Verification failed. Please check your code and try again.'));
	}

	redirect(303, '/');
});

export const logout = form('unchecked', async () => {
	await auth.api.signOut({
		headers: getRequestEvent().request.headers
	});
	redirect(303, '/');
});
