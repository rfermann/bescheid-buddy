<script lang="ts">
	import { page } from '$app/state';
	import { logout } from '../../routes/auth/auth.remote.ts';
	import { m } from '$lib/paraglide/messages.js';

	interface NavItem {
		id: string;
		label: () => string;
		href: string;
		disabled?: boolean;
		hidden?: boolean;
		badge?: () => string;
	}

	interface Props {
		isOpen: boolean;
		currentLabel: string;
		menuLabel: () => string;
		navItems: NavItem[];
		onToggle: () => void;
	}

	let { isOpen, currentLabel, menuLabel, navItems, onToggle }: Props = $props();

	let unreadCount = $state(2);
</script>

<button
	class="w-full flex items-center gap-2.5 p-3.5 rounded-2xl border border-white/14 bg-zinc-950 text-white shadow-2xl cursor-pointer hover:border-white/20 transition-colors"
	type="button"
	aria-haspopup="true"
	aria-expanded={isOpen}
	onclick={onToggle}
>
	<span class="text-sm tracking-widest uppercase text-white/70 font-black">{menuLabel()}</span>
	<span class="flex-1 min-w-0 font-black text-left whitespace-nowrap overflow-hidden text-ellipsis">
		{currentLabel}
	</span>
	<span class="flex-none text-primary transition-transform duration-150 ease-in-out" class:rotate-180={isOpen}
				aria-hidden="true">▾</span
	>
</button>

{#if isOpen}
	<div
		class="mt-2.5 rounded-3xl bg-zinc-950 border border-white/12 shadow-2xl p-3.5 max-h-[calc(100vh-160px)] overflow-auto"
		role="menu"
	>
		{#each navItems as item (item.id)}
			{#if !item.hidden}
				<a
					href={item.href}
					class={`w-full text-left flex items-center gap-2.5 p-3 rounded-xl border text-white cursor-pointer mb-2.5 transition-all disabled:opacity-55 disabled:cursor-not-allowed ${page.url.pathname === item.href ? 'border-accent bg-primary-bg' : 'border-white/10 bg-white/3'} ${item.disabled ? 'pointer-events-none' : ''}`}
					role="menuitem"
				>
					<div
						class={`w-2 h-2 rounded-full ${page.url.pathname === item.href ? 'bg-primary shadow-[0_0_14px_rgba(0,255,238,0.5)]' : 'bg-white/25'}`}
					></div>
					<span class="flex-1">{item.label()}</span>
					{#if item.badge}
						<span class="text-xs opacity-55 ml-2">{item.badge()}</span>
					{/if}
				</a>
			{/if}
		{/each}

		<a href="/buddypoints-verwalten"
			 class="w-full text-left flex items-center gap-2.5 p-3 rounded-xl border border-primary/30 bg-primary/10 text-white cursor-pointer mb-2.5 transition-all md:hidden"
			 role="menuitem">
			<span class="flex-1 font-semibold">{m.header_manage_buddypoints()}</span>
		</a>

		<a href="/nachrichten"
			 class="w-full text-left flex items-center justify-between gap-2.5 p-3 rounded-xl border border-primary/30 bg-primary/10 text-white cursor-pointer mb-2.5 transition-all md:hidden"
			 role="menuitem">
			<span class="flex-1 font-semibold">{m.header_messages()}</span>
			{#if unreadCount > 0}
				<span class="bg-[#ff4d6d] text-white text-xs rounded-full px-1.5 py-0.5">{unreadCount}</span>
			{/if}
		</a>

		<form {...logout.for("dropdown")} class="w-full md:hidden">
			<button type="submit"
							class="w-full text-left flex items-center gap-2.5 p-3 rounded-xl border border-primary/30 bg-primary/10 text-white cursor-pointer mb-2.5 transition-all"
							role="menuitem">
				<span class="flex-1 font-semibold">{m.header_logout()}</span>
			</button>
		</form>
	</div>
{/if}

<style>
    .rotate-180 {
        transform: rotate(180deg);
    }
</style>
