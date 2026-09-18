import { expect, test, vi } from 'vitest';

const mocks = vi.hoisted(() => ({ getDelivery: vi.fn() }));

vi.mock('./papierbrief.remote', () => ({
	getGenerationStatusQuery: vi.fn(),
	getPaperLetterDeliveryQuery: mocks.getDelivery
}));

import { DeliveryCancelledError, pollDelivery } from './papierbrief.client';

test('stops polling when delivery is sent or cancelled', async () => {
	const sent = { deliveryId: 'delivery-id', status: 'sent', pingenId: 'pingen-id' };
	const refresh = vi.fn();
	mocks.getDelivery.mockReturnValue({ current: sent, refresh });
	await expect(pollDelivery('delivery-id', () => false, vi.fn())).resolves.toBe(sent);
	expect(refresh).toHaveBeenCalledOnce();

	await expect(pollDelivery('delivery-id', () => true, vi.fn())).rejects.toBeInstanceOf(
		DeliveryCancelledError
	);
});
