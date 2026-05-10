import nodemailer from 'nodemailer';
import { env } from '$env/dynamic/private';

let transporter: nodemailer.Transporter | null = null;

const getTransporter = (): nodemailer.Transporter => {
	if (!transporter) {
		transporter = nodemailer.createTransport({
			host: env.SMTP_HOST,
			port: parseInt(env.SMTP_PORT || '587'),
			secure: true,
			auth: {
				user: env.SMTP_USER,
				pass: env.SMTP_PASS
			}
		});
	}
	return transporter;
};

export const sendEmail = async (options: { to: string; subject: string; html: string }) => {
	const transporter = getTransporter();

	// no need to await here
	transporter.sendMail({
		from: env.SMTP_FROM || env.SMTP_USER,
		to: options.to,
		subject: options.subject,
		html: options.html
	});
};
