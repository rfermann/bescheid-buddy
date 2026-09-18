<script lang="ts">
	type Props = {
		label: string;
		name: string;
		rows: number;
		value?: string | number;
		disabled?: boolean;
		required?: boolean;
		maxlength?: number;
		error?: string;
	};

	let {
		label,
		name,
		rows,
		value = $bindable(''),
		disabled = false,
		required = false,
		maxlength,
		error
	}: Props = $props();
</script>

<div class="flex flex-col gap-1">
	<label class="mb-2 text-xs font-extrabold tracking-widest uppercase" for={name}>{label}</label>
	<textarea
		{rows}
		id={name}
		{name}
		bind:value
		{disabled}
		{required}
		{maxlength}
		aria-invalid={error ? 'true' : undefined}
		aria-describedby={error ? `${name}-error` : undefined}
		class="rounded-xl border border-white/22 bg-white/5 disabled:opacity-50"
	></textarea>
	{#if error}
		<span id={`${name}-error`} class="text-sm text-red-300">{error}</span>
	{/if}
</div>
