export const load = async ({ locals }) => {
	console.log('layout.server.ts');
	return {
		bot: locals.bot,
		user: locals.user
	};
};
