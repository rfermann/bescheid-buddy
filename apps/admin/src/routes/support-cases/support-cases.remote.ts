import { query, getRequestEvent } from '$app/server';
import { error } from '@sveltejs/kit';
import { getAllSupportCases } from '../../generated';
import { getClientConfig } from '$lib/apiClient.server';
import { auth } from '$lib/server/auth';

export const getSupportCases = query(async () => {
	const supportCasesResponse = await getAllSupportCases({
		...getClientConfig()
	});

	const usersResponse = await auth.api.listUsers({
		query: {},
		headers: getRequestEvent().request.headers
	});

	if (supportCasesResponse.error || !supportCasesResponse.data || !supportCasesResponse.data.cases) {
		error(500, 'Failed to fetch support cases');
	}

	return supportCasesResponse?.data?.cases.map((c) => {
		const userName = usersResponse.users.find((u) => u.id === c.user_id)?.name;
		console.log({ user: userName });
		return {
			...c,
			userName: userName || 'Unknown'
		};
	}) ?? [];
});

