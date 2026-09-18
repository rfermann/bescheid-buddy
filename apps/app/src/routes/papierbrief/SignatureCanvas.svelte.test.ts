import { expect, test } from 'vitest';
import { userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';
import SignatureCanvas from './SignatureCanvas.svelte';

test('exposes its description and clear action to keyboard users', async () => {
	const screen = await render(SignatureCanvas, {
		label: 'Signature',
		description: 'Draw a signature with pointer or touch.',
		clearLabel: 'Clear signature',
		value: 'data:image/png;base64,test'
	});

	const canvas = screen.getByLabelText('Signature');
	await expect
		.element(canvas)
		.toHaveAttribute('aria-describedby', 'papierbrief-signature-description');
	await expect.element(screen.getByText('Draw a signature with pointer or touch.')).toBeVisible();

	const clear = screen.getByRole('button', { name: 'Clear signature' });
	(clear.element() as HTMLButtonElement).focus();
	await userEvent.keyboard('{Enter}');
	await expect.element(clear).toBeDisabled();
});
