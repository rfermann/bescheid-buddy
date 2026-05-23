<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		children: Snippet
	}

	let { children }: Props = $props();
</script>

<section class="[background:rgba(255,255,255,0.06)] rounded-2xl border-white/8 border p-4.5 mb-12">
	{@render children()}
</section>