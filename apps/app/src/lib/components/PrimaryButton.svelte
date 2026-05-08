<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		children?: Snippet
		text: string
	}

	let { children, text }: Props = $props();
</script>

<button
	class="px-3.5 py-2.5 rounded-2xl border font-bold border-primary/45 bg-primary/15 text-primary text-sm hover:border-primary/60 transition-colors relative">
	<span>{text}</span>
	{#if children}
		{@render children()}
	{/if}
</button>
