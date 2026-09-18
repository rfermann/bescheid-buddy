<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	type Props = {
		children?: Snippet;
		onclick?: () => void;
		text: string;
		type: HTMLButtonAttributes['type'];
		disabled?: boolean;
	};

	let { children, onclick, text, type = 'button', disabled = false }: Props = $props();
</script>

<button
	onclick={onclick ? onclick : undefined}
	{type}
	{disabled}
	class="relative rounded-2xl border border-primary/45 bg-primary/15 px-3.5 py-2.5 text-sm font-bold text-primary transition-colors hover:border-primary/60 disabled:cursor-not-allowed disabled:opacity-50"
>
	<span>{text}</span>
	{#if children}
		{@render children()}
	{/if}
</button>
