<script lang="ts">
	import './layout.css';
	import CanvasNeonBackground from '$lib/components/CanvasNeonBackground.svelte';
	import Header from '$lib/components/header/Header.svelte';
	import Footer from '$lib/components/footer/Footer.svelte';
	import { page } from '$app/state';
	import { fly, fade } from 'svelte/transition';

	let { children } = $props();
</script>

<CanvasNeonBackground />
<Header />
<main class="w-full p-4 max-w-225 mx-auto grow">
	{#key page.url.pathname}
		<div
			class="flex w-full flex-col"
			in:fly={{ y: 150, delay: 150, duration: 200, }}
			out:fade={{ duration: 150 }}
		>
			{@render children()}
		</div>
	{/key}
</main>
<Footer />
