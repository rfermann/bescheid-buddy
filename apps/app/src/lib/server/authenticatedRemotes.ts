import { command, form, query, getRequestEvent } from '$app/server';
import { error } from '@sveltejs/kit';
import { auth } from '$lib/server/auth';

async function assertAuthenticated() {
	const { request } = getRequestEvent();
	const session = await auth.api.getSession({ headers: request.headers });
	if (!session?.user) {
		error(401, 'Unauthorized');
	}
}

function authenticate(factory: (...a: never[]) => unknown, args: unknown[]): unknown {
	const handler = args[args.length - 1] as (...handlerArgs: unknown[]) => unknown;

	const guarded = async (...handlerArgs: unknown[]) => {
		await assertAuthenticated();
		return handler(...handlerArgs);
	};

	return (factory as (...a: unknown[]) => unknown)(...args.slice(0, -1), guarded);
}

export const authenticatedCommand = ((...args: unknown[]) =>
	authenticate(command, args)) as typeof command;

export const authenticatedForm = ((...args: unknown[]) => authenticate(form, args)) as typeof form;

export const authenticatedQuery = Object.assign((...args: unknown[]) => authenticate(query, args), {
	batch: (...args: unknown[]) => authenticate(query.batch, args),
	live: (...args: unknown[]) => authenticate(query.live, args)
}) as typeof query;
