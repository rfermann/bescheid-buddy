import type { PageServerLoad } from './$types';
import { getAllBuddyPackages } from '../../generated';
import { getClientConfig } from '$lib/server/apiClient';

export const load: PageServerLoad = async ({ url }) => {
	const sessionId = url.searchParams.get('session_id');
	const packageName = url.searchParams.get('name');
	const canceled = url.searchParams.get('canceled');

	let packagePoints: number | null = null;

	if (packageName) {
		try {
			const res = await getAllBuddyPackages({
				...getClientConfig()
			});

			if (!res.error && res.data?.packages) {
				const pkg = res.data.packages.find((p) => p.name === packageName);
				if (pkg) {
					packagePoints = pkg.points;
				}
			}
		} catch (error) {
			console.error('Failed to fetch package data:', error);
		}
	}

	return {
		purchaseSuccess: sessionId !== null && packageName !== null,
		packageName: packageName || null,
		packagePoints,
		purchaseCanceled: canceled === 'true'
	};
};
