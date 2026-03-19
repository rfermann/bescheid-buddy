<script lang="ts">
	import { page } from '$app/state';
	import { locales, localizeHref } from '$lib/paraglide/runtime';
	import './layout.css';
	import CanvasNeonBackground from '$lib/components/CanvasNeonBackground.svelte';
	import Header from '$lib/components/header/Header.svelte';
	import Footer from '$lib/components/footer/Footer.svelte';

	let { children } = $props();
</script>

<CanvasNeonBackground />
<Header />
<main class="main">
	{@render children()}
</main>

<div style="display:none">
	{#each locales as locale}
		<a href={localizeHref(page.url.pathname, { locale })}>{locale}</a>
	{/each}
</div>
<Footer />
