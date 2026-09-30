import { beforeEach, expect, test, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import { setLocale } from '$lib/paraglide/runtime';
import Page from './+page.svelte';

const mocks = vi.hoisted(() => ({ getDeliveries: vi.fn() }));

vi.mock('./deliveries.remote', () => ({ getPaperLetterDeliveriesQuery: mocks.getDeliveries }));

beforeEach(() => {
	vi.clearAllMocks();
	setLocale('en', { reload: false });
});

test('renders delivery details, statuses, dates, and available analysis links', async () => {
	mocks.getDeliveries.mockResolvedValue([
		{
			deliveryId: 'delivery-pending',
			caseId: 'case-pending',
			recipient: 'Jobcenter Berlin',
			topic: 'Pending letter',
			status: 'pending',
			sentAt: null,
			updatedAt: new Date('2026-01-02T10:00:00Z')
		},
		{
			deliveryId: 'delivery-processing',
			caseId: null,
			recipient: 'Health insurer',
			topic: 'Processing letter',
			status: 'processing',
			sentAt: new Date('2026-01-03T10:00:00Z'),
			updatedAt: new Date('2026-01-03T10:00:00Z')
		},
		{
			deliveryId: 'delivery-sent',
			caseId: 'case-sent',
			recipient: 'Tax office',
			topic: 'Sent letter',
			status: 'sent',
			sentAt: new Date('2026-01-04T10:00:00Z'),
			updatedAt: new Date('2026-01-04T10:00:00Z')
		},
		{
			deliveryId: 'delivery-failed',
			caseId: null,
			recipient: 'Court',
			topic: 'Failed letter',
			status: 'failed',
			sentAt: null,
			updatedAt: new Date('2026-01-05T10:00:00Z')
		},
		{
			deliveryId: 'delivery-review',
			caseId: null,
			recipient: 'Immigration office',
			topic: 'Review letter',
			status: 'manual_review',
			sentAt: null,
			updatedAt: new Date('2026-01-06T10:00:00Z')
		}
	]);

	const screen = await render(Page);
	for (const label of ['Pending', 'Processing', 'Sent', 'Failed', 'Manual review']) {
		await expect.element(screen.getByText(label, { exact: true }).first()).toBeVisible();
	}
	for (const text of ['Pending letter', 'Jobcenter Berlin', '2026']) {
		await expect.element(screen.getByText(text, { exact: text !== '2026' }).first()).toBeVisible();
	}
	await expect
		.element(screen.getByRole('link', { name: 'Pending letter' }))
		.toHaveAttribute('href', '/meine-papierbriefe/delivery-pending');
	await expect
		.element(screen.getByRole('link', { name: 'View analysis' }).first())
		.toHaveAttribute('href', '/meine-analysen/case-pending');
	expect(document.querySelectorAll('a[href^="/meine-analysen/"]')).toHaveLength(2);
});

test('renders nullable fields safely', async () => {
	mocks.getDeliveries.mockResolvedValue([
		{
			deliveryId: 'delivery-nullable',
			caseId: null,
			recipient: null,
			topic: null,
			status: 'future-status',
			sentAt: null,
			updatedAt: new Date('2026-01-02T10:00:00Z')
		}
	]);

	const screen = await render(Page);
	await expect.element(screen.getByText('No topic')).toBeVisible();
	await expect.element(screen.getByText('Unknown recipient')).toBeVisible();
	await expect.element(screen.getByText('Unknown status')).toBeVisible();
	expect(document.querySelector('a[href^="/meine-analysen/"]')).toBeNull();
});

test('renders an empty state with a new-letter link', async () => {
	mocks.getDeliveries.mockResolvedValue([]);

	const screen = await render(Page);
	await expect.element(screen.getByRole('heading', { name: 'No paper letters yet' })).toBeVisible();
	await expect
		.element(screen.getByRole('link', { name: 'Send paper letter' }).last())
		.toHaveAttribute('href', '/papierbrief');
});
