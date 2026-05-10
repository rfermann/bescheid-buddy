<script lang="ts">
	import './layout.css';
	import Footer from '$lib/components/Footer.svelte';
	import Header from '$lib/components/Header.svelte';
	import Navigation from '$lib/components/Navigation.svelte';
	import { init } from '@tma.js/sdk-svelte';
	import { browser } from '$app/environment';
	import { onMount, type Snippet } from 'svelte';
	import { authClient } from '$lib';
	import type { PageData } from './$types';
	import { goto } from '$app/navigation';

	let { children, data }: { children: Snippet, data: PageData } = $props();
	// const { bot, user } = data;

	console.log('layout.svelte');

	$inspect(data);

	onMount(() => {
		if (browser && data.bot && typeof data.user === 'undefined') {
			init();
			onMount(() => {
				const waitForInitData = async () => {
					let attempts = 0;
					while (attempts < 50) {
						if ((window as any).Telegram?.WebApp) {
							try {
								const result = await authClient.autoSignInFromMiniApp();
								await goto('/', {
									invalidateAll: true
								});
							} catch (error) {
								console.error('Not running in Mini App or auth failed');
							}
							break;
						}
						await new Promise((resolve) => setTimeout(resolve, 100));
						console.log('Waiting for Telegram WebApp to initialize...');
						attempts++;
					}
				};

				waitForInitData();
			});
		}
	});
</script>
<div class="min-h-screen flex flex-col pt-16">
	<Header user={data.user} />
	{#if data.user}
		<Navigation />
	{/if}
	<main class="w-full max-w-300 py-5 text-white mx-auto flex-1">
		{@render children()}
	</main>

	<Footer />
</div>
