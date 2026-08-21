import { command, form, query } from '$app/server';
import { error } from '@sveltejs/kit';
import {
	createAiModel,
	getAllModels,
	getModelProviders as fetchModelProviders,
	setModelEnabled
} from '../../generated';
import {
	vCreateAiModelBody,
	vSetModelEnabledBody,
	vSetModelEnabledPath
} from '../../generated/valibot.gen';
import { getClientConfig } from '$lib/apiClient.server';
import * as v from 'valibot';

export const getModels = query(async () => {
	const response = await getAllModels({
		...getClientConfig()
	});


	if (response.error || !response.data) {
		error(500, 'Failed to fetch AI models');
	}

	return response.data.models ?? [];
});

export const getModelProviders = query(async () => {
	const response = await fetchModelProviders({
		...getClientConfig()
	});

	if (response.error || !response.data) {
		error(500, 'Failed to fetch AI model providers');
	}

	return response.data.providers ?? [];
});

export const createModelForm = form(
	vCreateAiModelBody,
	async ({ name, provider, provider_model_id, sort_order }) => {
		const response = await createAiModel({
			...getClientConfig(),
			body: { name, provider, provider_model_id, sort_order }
		});

		if (response.error || !response.data) {
			return {
				success: false,
				message: response.error?.detail ?? 'Failed to create AI model'
			};
		}

		return { success: true };
	}
);

const setModelEnabledSchema = v.object({
	...vSetModelEnabledPath.entries,
	...vSetModelEnabledBody.entries
});

export const setModelEnabledCommand = command(setModelEnabledSchema, async ({ id, enabled }) => {
	const response = await setModelEnabled({
		...getClientConfig(),
		path: { id },
		body: { enabled }
	});

	if (response.error || !response.data) {
		return {
			success: false,
			message: response.error?.detail ?? 'Failed to update AI model'
		};
	}

	return { success: true };
});
