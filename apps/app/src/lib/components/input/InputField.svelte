<script lang="ts">
	import type { HTMLInputAttributes, HTMLInputTypeAttribute } from 'svelte/elements';

	type Props = {
		label: string;
		name: string;
		type: HTMLInputTypeAttribute;
		value?: string | number;
		disabled?: boolean;
		required?: boolean;
		maxlength?: number;
		autocomplete?: HTMLInputAttributes['autocomplete'];
		error?: string;
	};

	let {
		label,
		name,
		type,
		value = $bindable(''),
		disabled = false,
		required = false,
		maxlength,
		autocomplete,
		error
	}: Props = $props();
</script>

<div class="flex flex-col gap-1">
	<label class="mb-2 text-xs font-extrabold tracking-widest uppercase" for={name}>{label}</label>
	<input
		{type}
		id={name}
		{name}
		bind:value
		{disabled}
		{required}
		{maxlength}
		{autocomplete}
		aria-invalid={error ? 'true' : undefined}
		aria-describedby={error ? `${name}-error` : undefined}
		class="rounded-xl border border-white/22 bg-white/5 focus:border-primary/45 focus:ring-primary/45 focus:outline-none disabled:opacity-50"
	/>
	{#if error}
		<span id={`${name}-error`} class="text-sm text-red-300">{error}</span>
	{/if}
</div>
