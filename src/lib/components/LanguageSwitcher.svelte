<script lang="ts">
	import { getLocale, localizeHref, setLocale, locales } from '$lib/paraglide/runtime';
	import { goto } from '$app/navigation';
	import { fade, fly } from 'svelte/transition';

	type Locale = typeof locales[number];

	type Flag = {
		code: Locale,
		flag: string
	}

	const localeFlags: Flag[] = [
		{ code: 'de', flag: '🇩🇪' },
		{ code: 'en', flag: '🇬🇧' },
		{ code: 'pl', flag: '🇵🇱' },
		{ code: 'ru', flag: '🇷🇺' },
		{ code: 'ua', flag: '🇺🇦' },
		{ code: 'ar', flag: '🇦🇪' }
	];

	let isOpen = $state(false);
	let currentLocale = $state(localeFlags.find(item => item.code === (getLocale() || 'de')) || localeFlags[0]);

	const handleLocaleChange = (localeCode: Locale) => {
		currentLocale = localeFlags.find(item => item.code === localeCode) || localeFlags[0];
		setLocale(localeCode, {});
		isOpen = false;
		goto(localizeHref(`/${localeCode}`));
	};

	const toggleDropdown = (event: MouseEvent) => {
		event.stopPropagation();
		isOpen = !isOpen;
	};
</script>

<svelte:window on:click={() => isOpen = false} />

<div class="fixed right-4 bottom-4 flex flex-col gap-2 items-center">
	{#if isOpen}
		<div
			class="flex flex-col gap-2 items-center"
			transition:fly={{ y: 50, duration: 400 }}
			role="menu"
			tabindex="-1"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.key === 'Escape' && (isOpen = false)}
		>
			{#each localeFlags as locale (locale)}
				<button
					class="w-11 h-11 rounded-2xl border border-main transition-colors flex items-center justify-center {locale.code === currentLocale.code ? 'opacity-50' : ''}"
					transition:fade={{ delay: localeFlags.filter(l => l.code !== currentLocale.code).indexOf(locale) * 50, duration: 200 }}
					onclick={() => handleLocaleChange(locale.code)}
					aria-label={`Switch to ${locale.code}`}
				>
					{locale.flag}
				</button>
			{/each}
		</div>
	{/if}

	<button
		class="w-11 h-11 rounded-2xl border border-main flex items-center justify-center transition-colors"
		onclick={toggleDropdown}
		aria-label="Change language"
		aria-expanded={isOpen}
	>
		{currentLocale.flag}
	</button>
</div>


