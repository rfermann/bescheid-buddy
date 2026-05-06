<script lang="ts">
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	console.log('page.svelte');
</script>
<h1>Welcome to SvelteKit, {data?.user?.name ?? ''}</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>
