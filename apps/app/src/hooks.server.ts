import { sequence } from '@sveltejs/kit/hooks';
import { building } from '$app/environment';
import { auth } from '$lib/server/auth';
import { svelteKitHandler } from 'better-auth/svelte-kit';
import { type Handle, redirect } from '@sveltejs/kit';
import { getTextDirection } from '$lib/paraglide/runtime';
import { paraglideMiddleware } from '$lib/paraglide/server';

const handleParaglide: Handle = ({ event, resolve }) =>
	paraglideMiddleware(event.request, ({ request, locale }) => {
		event.request = request;

		return resolve(event, {
			transformPageChunk: ({ html }) =>
				html
					.replace('%paraglide.lang%', locale)
					.replace('%paraglide.dir%', getTextDirection(locale))
		});
	});

const handleBetterAuth: Handle = async ({ event, resolve }) => {
	const session = await auth.api.getSession({ headers: event.request.headers });
	const pathname = event.url.pathname;
	const isAuthRoute = pathname === '/auth' || pathname.startsWith('/auth/');
	const isBetterAuthRoute = pathname.startsWith('/api/auth/');

	if (session) {
		event.locals.session = session.session;
		event.locals.user = session.user;
	}

	const params = new URL(event.request.url).searchParams;
	const bot = params.get('bot_app');
	if (bot) {
		event.locals.bot = true;
	}

	console.log('hooks.server.ts');

	if (!bot && !event.locals.session && !isAuthRoute && !isBetterAuthRoute) {
		const redirectTo = pathname;
		const queryParams = event.url.search.replace('?', '&');
		redirect(303, `/auth/login?redirectTo=${encodeURIComponent(redirectTo)}${queryParams}`);
	}

	event.url.protocol = 'https:';

	return svelteKitHandler({ event, resolve, auth, building });
};

export const handle: Handle = sequence(handleParaglide, handleBetterAuth);
