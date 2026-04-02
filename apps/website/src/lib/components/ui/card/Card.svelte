<script lang="ts">
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<div
	class="p-5 mb-5 rounded-xl border border-main text-gray-200 shadow-[0_12px_30px_rgba(0,0,0,0.7),0_0_24px_rgba(141,116,201,0.25)] hover:shadow-[0_12px_30px_rgba(0,0,0,0.7),0_0_28px_rgba(0,255,238,0.2)] transition-shadow duration-300">

	{@render children()}
</div>


