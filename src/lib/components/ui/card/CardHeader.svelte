<script lang="ts">
	import type { Snippet } from 'svelte';

	let { title }: { title: string } = $props();
</script>

<h2 class="text-secondary uppercase text-2xl font-bold mb-3">{title}</h2>