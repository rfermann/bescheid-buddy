import { Renderer } from '@better-svelte-email/server';

let renderer: Renderer | null = null;

const getRenderer = (): Renderer => {
	if (!renderer) {
		renderer = new Renderer();
	}
	return renderer;
};

export const renderEmail = async (
	component: Parameters<Renderer['render']>[0],
	props?: Record<string, unknown>
): Promise<string> => {
	const renderer = getRenderer();
	return await renderer.render(component, { props });
};
