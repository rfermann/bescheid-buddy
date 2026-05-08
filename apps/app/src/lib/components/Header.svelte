<script lang="ts">
	import logo from '$lib/assets/logo.png';
	import { m } from '$lib/paraglide/messages.js';
	import { logout } from '../../routes/auth/auth.remote.ts';
	import type { User } from 'better-auth';
	import PrimaryButton from './PrimaryButton.svelte';
	import PrimaryLink from '$lib/components/PrimaryLink.svelte';

	type Props = {
		user?: User
	}

	let { user }: Props = $props();

	let unreadCount = $state(2);
</script>

<div
	class="text-white border-b border-white/8 bg-black/35 fixed top-0 left-0 right-0 z-50 flex items-center justify-between py-3.5 px-4.5">
	<div class="flex gap-3 items-center">
		<div class="w-10 h-10 rounded-xl border border-accent flex items-center justify-center">
			<img class="w-7 h-7" src={logo} alt="MachKinfach Logo">
		</div>
		<div>
			<div class="font-extrabold text-sm tracking-widest">BescheidBuddy
			</div>
			<div class="text-xs">Dein KI-Assistent</div>
		</div>
	</div>
	{#if user}
		<div class="flex gap-2.5 items-center">
			<PrimaryLink href="/buddypoints-verwalten" text={m.header_buy_buddypoints()} />
			<PrimaryLink href="/nachrichten" text={m.header_messages()}>
				{#if unreadCount > 0}
				<span
					class="absolute -top-2 -right-2 bg-[#ff4d6d] text-white text-xs rounded-full px-1.5 py-0.5 ml-2">
					{unreadCount}
				</span>
				{/if}
			</PrimaryLink>
			<form {...logout}>
				<button
					class="px-4 py-3 rounded-2xl border border-white/12 bg-white/4 text-white text-sm font-black hover:border-primary/22 hover:bg-white/7 transition-colors">
					{m.header_logout()}
				</button>
			</form>
		</div>
	{/if}
</div>