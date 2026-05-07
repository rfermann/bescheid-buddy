import { command, form, getRequestEvent } from '$app/server';
import { auth } from '$lib/server/auth';

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

export const sendOTP = form('unchecked', async ({ email }) => {
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

export const verifyOTP = form('unchecked', async ({ email, otp }) => {
	try {
		await auth.api.signInEmailOTP({
			body: {
				email: String(email ?? ''),
				otp: String(otp ?? '')
			}
		});

		return {
			success: true,
			message: 'Signed in successfully.'
		};
	} catch (error) {
		return {
			success: false,
			message: getMessage(error, 'Unable to verify the one-time password.')
		};
	}
});

export const logout = command('unchecked', async () => {
	try {
		const event = getRequestEvent();
		await auth.api.signOut({
			headers: event.request.headers
		});
		return {
			success: true,
			message: 'Signed out successfully.'
		};
	} catch (error) {
		return {
			success: false,
			message: getMessage(error, 'Unable to sign out.')
		};
	}
});
