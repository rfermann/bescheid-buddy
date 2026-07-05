import { command, form, query, getRequestEvent } from '$app/server';
import { error } from '@sveltejs/kit';
import {
	getSupportCaseById,
	getSupportCaseMessagesByCaseId,
	createSupportCaseMessage,
	updateSupportCaseStatus
} from '../../../generated';
import { getClientConfig } from '$lib/apiClient.server';
import { auth } from '$lib/server/auth';
import type { SupportCase, SupportCaseMessage } from '../../../generated';
import {
	vGetSupportCaseByIdPath,
	vGetSupportCaseMessagesByCaseIdPath,
	vUpdateSupportCaseStatusBody
} from '../../../generated/valibot.gen';
import * as v from 'valibot';

export const getSupportCaseDetail = query(vGetSupportCaseByIdPath, async ({ id }) => {
	const supportCaseResponse = await getSupportCaseById({
		...getClientConfig(),
		path: { id }
	});

	if (supportCaseResponse.error || !supportCaseResponse.data) {
		error(500, 'Failed to fetch support case');
	}

	const supportCase = supportCaseResponse.data.supportCase;

	// Fetch specific user by ID
	const user = await auth.api.getUser({
		query: { id: supportCase.user_id },
		headers: getRequestEvent().request.headers
	});

	return {
		...supportCase,
		userName: user?.name || 'Unknown',
		userEmail: user?.email || null
	} as SupportCase & { userName: string; userEmail: string | null };
});

export const getSupportCaseMessages = query(vGetSupportCaseMessagesByCaseIdPath, async ({ caseId }) => {
	const messagesResponse = await getSupportCaseMessagesByCaseId({
		...getClientConfig(),
		path: { caseId }
	});

	if (messagesResponse.error || !messagesResponse.data) {
		error(500, 'Failed to fetch support case messages');
	}

	return (messagesResponse.data.messages || []) as SupportCaseMessage[];
});

const sendMessageSchema = v.object({
	caseId: v.string(),
	message: v.pipe(v.string(), v.maxLength(255))
});

const updateStatusSchema = v.object({
	id: v.string(),
	...vUpdateSupportCaseStatusBody.entries
});

export const updateStatusCommand = command(updateStatusSchema, async ({ id, status }) => {
	try {
		const response = await updateSupportCaseStatus({
			...getClientConfig(),
			path: { id },
			body: { status }
		});

		if (response.error || !response.data) {
			return {
				success: false,
				message: 'Failed to update status'
			};
		}

		void getSupportCaseDetail({ id }).refresh();

		return {
			success: true,
			case: response.data.case
		};
	} catch (error) {
		return {
			success: false,
			message: 'Failed to update status'
		};
	}
});

export const sendMessageForm = form(sendMessageSchema, async ({ caseId, message }) => {
	try {
		const response = await createSupportCaseMessage({
			...getClientConfig(),
			path: { caseId },
			body: {
				message,
				isUserMessage: false,
				isInternal: false
			}
		});

		if (response.error || !response.data) {
			return {
				success: false,
				message: 'Failed to send message'
			};
		}

		void getSupportCaseDetail({ id: caseId }).refresh();
		void getSupportCaseMessages({ caseId }).refresh();

		return {
			success: true,
			message: response.data.message
		};
	} catch (error) {
		return {
			success: false,
			message: 'Failed to send message'
		};
	}
});
