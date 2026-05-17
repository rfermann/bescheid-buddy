import { Renderer } from '@better-svelte-email/server';

let renderer: Renderer | null = null;

const getRenderer = (): Renderer => {
	if (!renderer) {
		renderer = new Renderer();
	}
	return renderer;
};

const renderEmail = async (
	component: Parameters<Renderer['render']>[0],
	props?: Record<string, unknown>
): Promise<string> => {
	const renderer = getRenderer();
	return await renderer.render(component, { props });
};

const getOtpEmail = async () => {
	return (await import('../emails/otp-email.svelte')).default;
};

export { renderEmail, getOtpEmail };
