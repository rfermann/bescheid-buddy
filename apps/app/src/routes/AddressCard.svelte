<script lang="ts">
	let {
		title,
		name,
		address,
		street,
		house_number,
		postal_code,
		city,
		state,
		country,
		category
	}: {
		title: string;
		name?: string;
		address?: string;
		street?: string;
		house_number?: string;
		postal_code?: string;
		city?: string;
		state?: string;
		country?: string;
		category?: string;
	} = $props();

	let streetLine = $derived([street, house_number].filter(Boolean).join(' '));
	let cityLine = $derived([postal_code, city].filter(Boolean).join(' '));
	let hasStructuredAddress = $derived(!!(streetLine || cityLine || state || country));
</script>

<div class="rounded-lg border border-white/10 bg-white/5 p-4">
	<h3 class="text-sm font-medium text-[#b7b7c8]">{title}</h3>
	<p class="mt-1 text-base font-semibold text-white">{name || '—'}</p>
	{#if hasStructuredAddress}
		<div class="mt-1 text-sm text-[#b7b7c8]">
			{#if streetLine}<p>{streetLine}</p>{/if}
			{#if cityLine}<p>{cityLine}</p>{/if}
			{#if state}<p>{state}</p>{/if}
			{#if country}<p>{country}</p>{/if}
		</div>
	{:else if address}
		<p class="mt-1 whitespace-pre-line text-sm text-[#b7b7c8]">{address}</p>
	{/if}
	{#if category && category !== 'unknown'}
		<p class="mt-1 text-sm text-[#b7b7c8]">{category}</p>
	{/if}
</div>
