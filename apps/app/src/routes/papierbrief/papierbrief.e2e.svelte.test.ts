import { beforeEach, expect, test, vi } from 'vitest';
import { userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';
import { setLocale } from '$lib/paraglide/runtime';
import Page from './+page.svelte';

const mocks = vi.hoisted(() => ({
	url: 'http://localhost/papierbrief',
	createGeneration: vi.fn(),
	createPDF: vi.fn(),
	createUpload: vi.fn(),
	confirmUpload: vi.fn(),
	createDelivery: vi.fn(),
	pollGeneration: vi.fn(),
	pollDelivery: vi.fn(),
	putPDF: vi.fn(),
	getAnalysis: vi.fn(),
	refreshPoints: vi.fn()
}));

vi.mock('$app/state', () => ({
	page: {
		get url() {
			return new URL(mocks.url);
		}
	}
}));
vi.mock('../analyse.remote', () => ({
	getModelsQuery: vi.fn().mockResolvedValue([{ id: 'model-1', name: 'Test model' }])
}));
vi.mock('../meine-analysen/analyses.remote', () => ({ getAnalysisQuery: mocks.getAnalysis }));
vi.mock('../meine-buddypoints/points.remote', () => ({
	getBuddyPointsQuery: () =>
		Object.assign(Promise.resolve({ monthly_points: 100, lifetime_points: 0 }), {
			refresh: mocks.refreshPoints
		})
}));
vi.mock('./papierbrief.remote', () => ({
	createGenerationCommand: mocks.createGeneration,
	createPDFCommand: mocks.createPDF,
	createPaperLetterUploadCommand: mocks.createUpload,
	confirmPaperLetterUploadCommand: mocks.confirmUpload,
	createPaperLetterDeliveryCommand: mocks.createDelivery
}));
vi.mock('./papierbrief.client', () => ({
	GenerationCancelledError: class extends Error {},
	GenerationTimedOutError: class extends Error {},
	DeliveryCancelledError: class extends Error {},
	DeliveryTimedOutError: class extends Error {},
	pollGeneration: mocks.pollGeneration,
	pollDelivery: mocks.pollDelivery
}));
vi.mock('./papierbrief.upload.client', () => ({
	maxPaperLetterPDFBytes: 15_728_640,
	paperLetterPDFWithinLimit: (blob: Blob) => blob.size <= 15_728_640,
	putPaperLetterPDF: mocks.putPDF
}));
vi.mock('svelte-sonner', () => ({ toast: { error: vi.fn(), success: vi.fn() } }));

const generated = (germanLetter: string) => ({
	status: 'completed',
	result: {
		subject: 'Widerspruch',
		germanLetter,
		userLanguageLetter: 'Translated letter',
		explanation: 'Explanation'
	}
});

const fillManualLetter = async (screen: Awaited<ReturnType<typeof render>>) => {
	await expect
		.element(screen.getByRole('heading', { name: 'Create a paper letter' }))
		.toBeVisible();
	await userEvent.fill(document.querySelector('[name="sender-name"]')!, 'Max Mustermann');
	await userEvent.fill(
		document.querySelector('[name="sender-address"]')!,
		'Musterweg 1\n10115 Berlin'
	);
	await userEvent.fill(document.querySelector('[name="recipient-name"]')!, 'Jobcenter');
	await userEvent.fill(
		document.querySelector('[name="recipient-address"]')!,
		'Amtstraße 2\n20095 Hamburg'
	);
	await userEvent.fill(document.querySelector('[name="subject"]')!, 'Widerspruch');
	await userEvent.fill(document.querySelector('[name="instruction"]')!, 'Bitte widersprechen.');
};

beforeEach(() => {
	vi.clearAllMocks();
	mocks.url = 'http://localhost/papierbrief';
	setLocale('en', { reload: false });
	mocks.createGeneration.mockResolvedValue({
		ok: true,
		generation: { generationId: 'generation-1' }
	});
	mocks.pollGeneration.mockResolvedValue(generated('Generated German letter'));
	mocks.createPDF.mockResolvedValue({
		ok: true,
		bytes: new Uint8Array([37, 80, 68, 70]),
		filename: 'papierbrief.pdf',
		pageCount: 2
	});
	mocks.createUpload.mockResolvedValue({
		ok: true,
		uploadId: 'upload-1',
		uploadUrl: 'https://upload.test/pdf'
	});
	mocks.putPDF.mockResolvedValue(true);
	mocks.confirmUpload.mockResolvedValue({
		ok: true,
		uploadId: 'upload-1',
		status: 'confirmed',
		sizeBytes: 4,
		pageCount: 2,
		deliveryPrice: 70
	});
	mocks.createDelivery.mockResolvedValue({
		ok: true,
		delivery: { deliveryId: 'delivery-1', status: 'pending' }
	});
	mocks.pollDelivery.mockResolvedValue({
		deliveryId: 'delivery-1',
		status: 'sent',
		pingenId: 'pingen-1'
	});
});

test('runs manual Generate, Improve, PDF, generated upload, and mocked delivery as one flow', async () => {
	const screen = await render(Page);
	await fillManualLetter(screen);
	await userEvent.click(screen.getByRole('button', { name: 'Generate letter' }));
	await expect
		.element(screen.getByLabelText('Current German letter'))
		.toHaveValue('Generated German letter');

	await userEvent.fill(
		screen.getByLabelText('Current German letter'),
		'Locally edited German letter'
	);
	await userEvent.fill(screen.getByLabelText('What should be improved?'), 'Make it shorter.');
	mocks.pollGeneration.mockResolvedValueOnce(generated('Improved German letter'));
	await userEvent.click(screen.getByRole('button', { name: 'Improve letter' }));
	expect(mocks.createGeneration).toHaveBeenLastCalledWith(
		expect.objectContaining({
			mode: 'improve',
			input: expect.objectContaining({
				currentGermanLetter: 'Locally edited German letter',
				feedback: 'Make it shorter.'
			})
		})
	);

	await userEvent.click(screen.getByRole('button', { name: 'Prepare PDF' }));
	await expect.element(screen.getByRole('button', { name: 'Download PDF' })).toBeVisible();
	const download = vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(() => {});
	await userEvent.click(screen.getByRole('button', { name: 'Download PDF' }));
	expect(download).toHaveBeenCalledOnce();
	await userEvent.click(screen.getByRole('button', { name: 'Upload and validate PDF' }));
	expect(mocks.putPDF).toHaveBeenCalledWith('https://upload.test/pdf', expect.any(Blob));
	await userEvent.click(
		screen.getByRole('button', { name: 'Confirm delivery for 70 BuddyPoints' })
	);
	await expect.element(screen.getByText('Sent.')).toBeVisible();
	expect(mocks.createDelivery).toHaveBeenCalledOnce();

	screen.unmount();
	const reloaded = await render(Page);
	await expect
		.element(reloaded.getByRole('heading', { name: 'Create a paper letter' }))
		.toBeVisible();
	expect(document.querySelector('[name="german-letter"]')).toBeNull();
	await expect
		.element(reloaded.getByRole('button', { name: 'Upload and validate PDF' }))
		.toBeDisabled();
});

test('prefills from analysis and supports a ready PDF upload', async () => {
	const caseId = '00000000-0000-4000-8000-000000000001';
	const documentId = '00000000-0000-4000-8000-000000000002';
	mocks.url = `http://localhost/papierbrief?analyse=${caseId}&dokument=${documentId}`;
	mocks.getAnalysis.mockResolvedValue({
		documents: [
			{
				documentId,
				fileName: 'bescheid.pdf',
				analysis: {
					status: 'completed',
					metadata: {
						receiver: { name: 'Max Mustermann', address: 'Musterweg 1, 10115 Berlin' },
						sender: {
							name: 'Jobcenter',
							address: 'Amtstraße 2, 20095 Hamburg',
							category: 'Jobcenter'
						},
						subject: { original: 'Bescheid' },
						reference_number: 'JC-42',
						letter_type: 'widerspruch'
					}
				}
			}
		]
	});

	const screen = await render(Page);
	await expect
		.element(
			screen.getByText(
				'Details were prefilled from bescheid.pdf. Review and edit them before generating.'
			)
		)
		.toBeVisible();
	await expect
		.element(document.querySelector<HTMLElement>('[name="recipient-name"]')!)
		.toHaveValue('Jobcenter');
	await userEvent.upload(
		screen.getByLabelText('Ready PDF'),
		new File(['%PDF'], 'ready.pdf', { type: 'application/pdf' })
	);
	await userEvent.click(screen.getByRole('button', { name: 'Upload and validate PDF' }));
	expect(mocks.createUpload).toHaveBeenCalledWith({ filename: 'ready.pdf' });

	mocks.createDelivery
		.mockResolvedValueOnce({ ok: false, error: 'request_failed' })
		.mockResolvedValueOnce({ ok: true, delivery: { deliveryId: 'delivery-1', status: 'pending' } });
	mocks.pollDelivery.mockResolvedValueOnce({ deliveryId: 'delivery-1', status: 'manual_review' });
	const deliveryButton = screen.getByRole('button', {
		name: 'Confirm delivery for 70 BuddyPoints'
	});
	await userEvent.click(deliveryButton);
	await userEvent.click(deliveryButton);
	const firstRequest = mocks.createDelivery.mock.calls[0][0].clientRequestId;
	expect(mocks.createDelivery.mock.calls[1][0].clientRequestId).toBe(firstRequest);
	await expect
		.element(screen.getByText('Delivery is being reviewed manually. No automatic refund was made.'))
		.toBeVisible();
});

test('shows insufficient BP and terminal-refund states', async () => {
	const screen = await render(Page);
	await fillManualLetter(screen);
	mocks.createGeneration.mockResolvedValueOnce({ ok: false, error: 'insufficient_points' });
	await userEvent.click(screen.getByRole('button', { name: 'Generate letter' }));
	await expect
		.element(screen.getByRole('alert'))
		.toHaveTextContent('You do not have enough BuddyPoints');

	mocks.createGeneration.mockResolvedValueOnce({
		ok: true,
		generation: { generationId: 'generation-2' }
	});
	mocks.pollGeneration.mockResolvedValueOnce({
		status: 'failed',
		errorMessage: 'invalid_generation_output',
		refunded: true
	});
	await userEvent.click(screen.getByRole('button', { name: 'Generate letter' }));
	await expect
		.element(screen.getByRole('alert'))
		.toHaveTextContent('The 10 BuddyPoints were refunded.');
});

test.each([
	['en', 'Create a paper letter', 'ltr'],
	['de', 'Papierbrief erstellen', 'ltr'],
	['ru', 'Создать бумажное письмо', 'ltr'],
	['ua', 'Створити паперовий лист', 'ltr'],
	['pl', 'Utwórz list papierowy', 'ltr'],
	['ar', 'إنشاء رسالة ورقية', 'rtl']
])('renders the %s locale', async (locale, title, direction) => {
	setLocale(locale as 'en', { reload: false });
	const screen = await render(Page);
	await expect.element(screen.getByRole('heading', { name: title })).toBeVisible();
	await expect
		.element(document.querySelector<HTMLElement>(`[dir="${direction}"]`)!)
		.toBeInTheDocument();
});
