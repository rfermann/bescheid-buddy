import { beforeEach, expect, test, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import { setLocale } from '$lib/paraglide/runtime';
import Page from './+page.svelte';

const mocks = vi.hoisted(() => ({ getDetail: vi.fn() }));

vi.mock('$app/state', () => ({ page: { params: { id: '11111111-1111-4111-8111-111111111111' } } }));
vi.mock('../deliveries.remote', () => ({
	getPaperLetterDeliveryDetailQuery: mocks.getDetail
}));

beforeEach(() => {
	vi.clearAllMocks();
	setLocale('en', { reload: false });
	mocks.getDetail.mockResolvedValue({
		delivery: {
			deliveryId: '11111111-1111-4111-8111-111111111111',
			topic: 'Benefit appeal',
			recipient: 'Jobcenter Berlin',
			status: 'sent',
			createdAt: new Date('2026-01-01T10:00:00Z'),
			updatedAt: new Date('2026-01-02T10:00:00Z'),
			sentAt: new Date('2026-01-02T09:00:00Z'),
			pageCount: 2,
			buddyPointsCost: 70,
			deliveryProduct: 'fast',
			printMode: 'simplex',
			printSpectrum: 'color',
			addressPosition: 'left'
		},
		pdfURLs: {
			inlineUrl: 'about:blank#inline-pdf',
			downloadUrl: '/letter.pdf?signature=test'
		}
	});
});

test('renders owned delivery metadata, inline PDF, and download link', async () => {
	const screen = await render(Page);

	await expect.element(screen.getByRole('heading', { name: 'Benefit appeal' })).toBeVisible();
	await expect.element(screen.getByText('Jobcenter Berlin')).toBeVisible();
	await expect.element(screen.getByText('Sent', { exact: true })).toBeVisible();
	await expect.element(screen.getByText('2', { exact: true })).toBeVisible();
	await expect.element(screen.getByText('70 BuddyPoints')).toBeVisible();
	await expect
		.element(screen.getByRole('link', { name: 'Download PDF' }))
		.toHaveAttribute('href', '/letter.pdf?signature=test');
	await expect
		.element(screen.getByTitle('Paper letter PDF preview'))
		.toHaveAttribute('src', 'about:blank#inline-pdf');
	await expect
		.element(screen.getByRole('link', { name: 'Back to paper letters' }))
		.toHaveAttribute('href', '/meine-papierbriefe');
});
