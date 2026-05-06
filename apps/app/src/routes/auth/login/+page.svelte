<script lang="ts">
	import { goto } from '$app/navigation';
	import { slide } from 'svelte/transition';
	import { sendOTP, verifyOTP } from './auth.remote';
	import { m } from '$lib/paraglide/messages.js';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let step = $state(1);

	const focusOnMount = (node: HTMLInputElement) => {
		queueMicrotask(() => {
			node.focus();
			node.select();
		});

		return {
			destroy() {
			}
		};
	};

	sendOTP.enhance(async ({ submit }) => {
		await submit();
	});

	verifyOTP.enhance(async ({ submit }) => {
		await submit();
	});

	$effect(() => {
		if (sendOTP.result?.success) {
			step = 2;
		}

	});

	$effect(() => {
		if (verifyOTP.result?.success) {
			goto('/dashboard');
		}
	});

	const focusEmailOnMount = (node: HTMLInputElement) => {
		if (step === 2) {
			queueMicrotask(() => {
				node.focus();
				node.select();
			});
		}

		return {
			destroy() {
			}
		};
	};
</script>

<div class="px-4 py-12 sm:px-6 lg:px-8">

	<div
		class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full max-w-md rounded-[18px] border border-white/10 bg-white/6 p-6 shadow-[0_20px_60px_rgba(0,0,0,0.55)] backdrop-blur-sm sm:p-7">
		<div>
			<h1>{data.bot ? 'Bot' : ''}</h1>
			<h2 class="mb-1 text-3xl font-extrabold text-white">
				{step === 1 ? m.login_signInWithEmail() : m.login_enterOTP()}
			</h2>
			{#key step}
				<p transition:slide={{ duration: 180 }} class="text-sm leading-6 text-[#b7b7c8]">
					{step === 1
						? m.login_sendOTPDescription()
						: m.login_checkEmailForCode({ email: sendOTP.fields.email.value() })}
				</p>
			{/key}
		</div>

		<div class="mt-6 space-y-5">
			{#if sendOTP.result?.message}
				<div
					class={sendOTP.result.success
						? 'rounded-[14px] border border-emerald-400/30 bg-emerald-400/10 px-4 py-3 text-sm text-emerald-200'
						: 'rounded-[14px] border border-red-400/30 bg-red-400/10 px-4 py-3 text-sm text-red-200'}
				>
					{sendOTP.result.message}
				</div>
			{/if}

			<div>
				<label for="email" class="mb-2 block text-xs font-extrabold tracking-[0.1em] text-white/85 uppercase">
					{m.login_emailAddress()}
				</label>
				{#each sendOTP.fields.email.issues() as issue, index (issue.message + index)}
					<p class="mt-1 text-sm text-red-300">{issue.message}</p>
				{/each}

				{#if step === 1}
					<form {...sendOTP}>
						<input
							use:focusEmailOnMount
							id="email"
							{...sendOTP.fields.email.as('email')}
							autocomplete="email"
							class="block w-full rounded-[14px] border border-white/12 bg-white/4 px-4 py-3 text-sm text-white placeholder:text-white/35 focus:border-primary/40 focus:bg-white/7 focus:outline-none focus:ring-0"
							placeholder={m.login_enterEmail()}
						/>

						<div class="mt-5">
							<button
								type="submit"
								disabled={!!sendOTP.pending}
								class="inline-flex w-full items-center justify-center rounded-[14px] border border-primary/40 bg-primary/16 px-4 py-3 text-sm font-black text-primary transition-all duration-150 hover:-translate-y-px hover:border-primary/50 disabled:cursor-not-allowed disabled:opacity-60"
							>
								{sendOTP.pending ? m.login_sending() : m.login_sendOTP()}
							</button>
						</div>
					</form>
				{:else}
					<div class="space-y-3">
						<input
							id="email"
							value={sendOTP.fields.email.value()}
							disabled
							class="block w-full rounded-[14px] border border-white/10 bg-white/3 px-4 py-3 text-sm text-white/65 placeholder:text-white/35 focus:outline-none"
							placeholder={m.login_enterEmail()}
						/>

						<form {...sendOTP}>
							<input {...sendOTP.fields.email.as('hidden', sendOTP.fields.email.value())} />
							<button
								type="submit"
								disabled={!!sendOTP.pending}
								class="inline-flex w-full items-center justify-center rounded-[14px] border border-white/12 bg-white/4 px-4 py-3 text-sm font-black text-white transition-all duration-150 hover:-translate-y-px hover:border-primary/22 disabled:cursor-not-allowed disabled:opacity-60"
							>
								{sendOTP.pending ? m.login_resending() : m.login_resendCode()}
							</button>
						</form>
					</div>
				{/if}
			</div>

			{#if step === 2}
				<form {...verifyOTP} transition:slide={{ duration: 220 }} class="space-y-5 overflow-hidden">
					{#if verifyOTP.result?.message}
						<div
							class={verifyOTP.result.success
								? 'rounded-[14px] border border-emerald-400/30 bg-emerald-400/10 px-4 py-3 text-sm text-emerald-200'
								: 'rounded-[14px] border border-red-400/30 bg-red-400/10 px-4 py-3 text-sm text-red-200'}
						>
							{verifyOTP.result.message}
						</div>
					{/if}

					<input {...verifyOTP.fields.email.as('hidden', sendOTP.fields.email.value())} />

					<div>
						<label for="otp" class="mb-2 block text-xs font-extrabold tracking-widest text-white/85 uppercase">
							{m.login_oneTimePassword()}
						</label>
						{#each verifyOTP.fields.otp.issues() as issue, index (issue.message + index)}
							<p class="mt-1 text-sm text-red-300">{issue.message}</p>
						{/each}
						<div>
							<input
								use:focusOnMount
								id="otp"
								{...verifyOTP.fields.otp.as('text')}
								inputmode="numeric"
								pattern="[0-9]*"
								maxlength="6"
								class="block w-full rounded-[14px] border border-white/12 bg-white/4 px-4 py-3 text-center text-lg tracking-[0.35em] text-white placeholder:text-white/35 focus:border-primary/40 focus:bg-white/7 focus:outline-none focus:ring-0"
								placeholder={m.login_enterCode()}
							/>
						</div>
					</div>

					<div class="flex gap-3">
						<button
							type="submit"
							disabled={!!verifyOTP.pending}
							class="inline-flex flex-1 items-center justify-center rounded-[14px] border border-primary/40 bg-primary/16 px-4 py-3 text-sm font-black text-primary transition-all duration-150 hover:-translate-y-px hover:border-primary/50 disabled:cursor-not-allowed disabled:opacity-60"
						>
							{verifyOTP.pending ? m.login_verifying() : m.login_signIn()}
						</button>
					</div>
				</form>
			{/if}
		</div>
	</div>
</div>
