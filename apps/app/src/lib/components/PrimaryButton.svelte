<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	type Props = {
		children?: Snippet
		onclick?: () => void
		text: string
		type: HTMLButtonAttributes['type']
	}

	let { children, onclick, text, type = 'button' }: Props = $props();
</script>

<button
	onclick={onclick ? onclick : undefined}
	{type}
	class="px-3.5 py-2.5 rounded-2xl border font-bold border-primary/45 bg-primary/15 text-primary text-sm hover:border-primary/60 transition-colors relative">
	<span>{text}</span>
	{#if children}
		{@render children()}
	{/if}
</button>
