import { afterEach, expect, test, vi } from 'vitest';
import {
	maxPaperLetterPDFBytes,
	paperLetterPDFWithinLimit,
	putPaperLetterPDF
} from './papierbrief.upload.client';

afterEach(() => vi.unstubAllGlobals());

test('enforces the size hint and performs the direct PDF PUT', async () => {
	expect(paperLetterPDFWithinLimit(new Blob([new Uint8Array(maxPaperLetterPDFBytes)]))).toBe(true);
	expect(paperLetterPDFWithinLimit(new Blob([new Uint8Array(maxPaperLetterPDFBytes + 1)]))).toBe(
		false
	);

	const pdf = new Blob(['%PDF'], { type: 'application/pdf' });
	const fetchMock = vi
		.fn()
		.mockResolvedValueOnce({ ok: true })
		.mockResolvedValueOnce({ ok: false })
		.mockRejectedValueOnce(new Error('network'));
	vi.stubGlobal('fetch', fetchMock);

	await expect(putPaperLetterPDF('https://storage.example/upload', pdf)).resolves.toBe(true);
	expect(fetchMock).toHaveBeenCalledWith('https://storage.example/upload', {
		method: 'PUT',
		body: pdf,
		headers: { 'Content-Type': 'application/pdf' }
	});
	await expect(putPaperLetterPDF('https://storage.example/upload', pdf)).resolves.toBe(false);
	await expect(putPaperLetterPDF('https://storage.example/upload', pdf)).resolves.toBe(false);
});
