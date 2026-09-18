import type { DeliveryBody, GenerationBody } from '../../generated';
import { getGenerationStatusQuery, getPaperLetterDeliveryQuery } from './papierbrief.remote';

const POLL_INTERVAL_MS = 1500;
const GENERATION_TIMEOUT_MS = 3 * 60 * 1000;

const wait = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export class GenerationCancelledError extends Error {}
export class GenerationTimedOutError extends Error {}
export class DeliveryCancelledError extends Error {}
export class DeliveryTimedOutError extends Error {}

export const pollGeneration = async (
	generationId: string,
	isCancelled: () => boolean,
	onStatus: (generation: GenerationBody) => void
) => {
	const statusQuery = getGenerationStatusQuery({ generationId });
	const deadline = Date.now() + GENERATION_TIMEOUT_MS;

	while (Date.now() <= deadline) {
		if (isCancelled()) {
			throw new GenerationCancelledError();
		}
		await statusQuery.refresh();
		if (isCancelled()) {
			throw new GenerationCancelledError();
		}

		const generation = statusQuery.current;
		if (generation) {
			onStatus(generation);
			if (generation.status === 'completed' || generation.status === 'failed') {
				return generation;
			}
		}
		await wait(POLL_INTERVAL_MS);
	}

	throw new GenerationTimedOutError();
};

export const pollDelivery = async (
	deliveryId: string,
	isCancelled: () => boolean,
	onStatus: (delivery: DeliveryBody) => void
) => {
	const statusQuery = getPaperLetterDeliveryQuery({ deliveryId });
	const deadline = Date.now() + GENERATION_TIMEOUT_MS;

	while (Date.now() <= deadline) {
		if (isCancelled()) {
			throw new DeliveryCancelledError();
		}
		await statusQuery.refresh();
		if (isCancelled()) {
			throw new DeliveryCancelledError();
		}

		const delivery = statusQuery.current;
		if (delivery) {
			onStatus(delivery);
			if (['sent', 'failed', 'manual_review'].includes(delivery.status)) {
				return delivery;
			}
		}
		await wait(POLL_INTERVAL_MS);
	}

	throw new DeliveryTimedOutError();
};
