import { command, form, query } from '$app/server';
import * as v from 'valibot';
import {
	createSupportCase,
	getSupportCasesByUser,
	getSupportCaseById,
	getSupportCaseMessagesByCaseId,
	createSupportCaseMessage,
	updateSupportCaseStatus
} from '../../generated';
import { getClientConfig } from '$lib/server/apiClient';
import {
	vGetSupportCaseByIdPath,
	vGetSupportCaseMessagesByCaseIdPath
} from '../../generated/valibot.gen.ts';

const vCreateSupportCaseFormInput = v.object({
	message: v.pipe(v.string(), v.minLength(3, 'The string must be 3 or more characters long.')),
	subject: v.pipe(v.string(), v.minLength(3, 'The string must be 3 or more characters long.'))
});

const createSupportCaseForm = form(vCreateSupportCaseFormInput, async ({ message, subject }) => {
	const res = await createSupportCase({
		body: {
			category: 'general',
			message,
			subject
		}
	});

	if (res.error) {
		console.log(JSON.stringify(res.error, null, 2));
	}
});

const getSupportCaseByIdQuery = query(vGetSupportCaseByIdPath, async ({ id }) => {
	const res = await getSupportCaseById({
		...getClientConfig(),
		path: {
			id: id
		}
	});

	if (res.error) {
		console.log({ error: res.error });
		throw new Error('Failed to fetch support case');
	}

	return res.data;
});

const getSupportCaseMessagesQuery = query(
	vGetSupportCaseMessagesByCaseIdPath,
	async ({ caseId }) => {
		const res = await getSupportCaseMessagesByCaseId({
			...getClientConfig(),
			path: {
				caseId
			}
		});

		if (res.error) {
			throw new Error('Failed to fetch support cases');
		}

		return res.data;
	}
);

const getSupportCasesQuery = query('unchecked', async () => {
	const res = await getSupportCasesByUser({
		...getClientConfig()
	});

	if (res.error) {
		console.log({ error: res.error });
		throw new Error('Failed to fetch support cases');
	}

	return res.data;
});

const vSendMessageFormInput = v.object({
	caseId: v.string(),
	message: v.pipe(v.string(), v.minLength(1, 'Message is required'), v.maxLength(255))
});

const sendMessageForm = form(vSendMessageFormInput, async ({ caseId, message }) => {
	const res = await createSupportCaseMessage({
		...getClientConfig(),
		path: { caseId },
		body: {
			message,
			isUserMessage: true,
			isInternal: false
		}
	});

	if (res.error) {
		console.log({ error: res.error });
		throw new Error('Failed to send message');
	}

	void getSupportCaseMessagesQuery({ caseId }).refresh();

	return res.data;
});

const vUpdateStatusInput = v.string();

const resolveSupportCaseCommand = command(vUpdateStatusInput, async (id) => {
	const res = await updateSupportCaseStatus({
		...getClientConfig(),
		path: { id },
		body: { status: 'resolved' }
	});

	if (res.error) {
		console.log({ error: res.error });
		throw new Error('Failed to resolve support case');
	}

	void getSupportCaseByIdQuery({ id }).refresh();

	return res.data;
});

const closeSupportCaseCommand = command(vUpdateStatusInput, async (id) => {
	const res = await updateSupportCaseStatus({
		...getClientConfig(),
		path: { id },
		body: { status: 'closed' }
	});

	if (res.error) {
		console.log({ error: res.error });
		throw new Error('Failed to close support case');
	}

	void getSupportCaseByIdQuery({ id }).refresh();

	return res.data;
});

export {
	getSupportCaseByIdQuery,
	getSupportCaseMessagesQuery,
	getSupportCasesQuery,
	createSupportCaseForm,
	sendMessageForm,
	resolveSupportCaseCommand,
	closeSupportCaseCommand
};
