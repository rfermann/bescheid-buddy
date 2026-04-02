import * as v from 'valibot';
import { form } from '$app/server';
import { m } from '$lib/paraglide/messages.js';
import { transporter } from '$lib/mail';
import { getRequestEvent } from '$app/server';
import { getLocale } from '$lib/paraglide/runtime';
import { SMTP_USER } from '$env/static/private';
import { error } from '@sveltejs/kit';

const contactSchema = v.object({
	name: v.pipe(
		v.optional(v.string(), 'asdfasdf'),
		v.minLength(3, () => m.validation_name_min())
	),
	email: v.pipe(
		v.optional(v.string(), 'rene.fermann@gmx.de'),
		v.email(() => m.validation_email_invalid())
	),
	message: v.pipe(
		v.optional(v.string(), 'asdfasdf'),
		v.minLength(3, () => m.validation_message_min())
	),
	type: v.enum({ general: 'general', partner: 'partner' })
});

const submitContact = form(contactSchema, async ({ name, email, message, type }) => {
	const { getClientAddress, request } = getRequestEvent();

	const to = 'kontakt@machkinfach.de';
	const subject = `MachKinfach ${type.charAt(0).toUpperCase() + type.slice(1)} Contact Form`;
	const from = SMTP_USER;
	const fromName = 'BescheidBuddy';

	const body =
		`New message from machkinfach.de contact form\n\n` +
		`Name: ${name}\n` +
		`Email: ${email}\n` +
		`Language: ${getLocale()}\n` +
		`IP: ${getClientAddress()}\n` +
		`User-Agent: ${request.headers.get('user-agent')}\n\n` +
		`Message:\n${message}\n`;

	const result = await transporter.sendMail({
		to,
		subject,
		text: body,
		from: `${fromName} <${from}>`,
		replyTo: email
	});

	if (!result.accepted.includes(to)) {
		error(
			404,
			'Es trat ein unerwarteter Fehler auf. Bitte probiere es in einigen Minuten noch einmal.'
		);
	}

	return {
		success: true
	};
});

export { submitContact };
