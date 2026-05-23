import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { getRequestEvent } from '$app/server';

export const getClientConfig = () => {
	const requestEvent = getRequestEvent();
	const cookies = requestEvent.cookies.getAll();

	return {
		baseUrl: PUBLIC_API_BASE_URL,
		headers: {
			Cookie: cookies.map((cookie) => `${cookie.name}=${cookie.value}`).join('; ')
		}
	};
};
