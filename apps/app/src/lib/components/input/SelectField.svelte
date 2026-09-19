<script lang="ts">
	type Props = {
		label: string;
		name: string;
		options: readonly { value: string; label: string }[];
		value?: string;
		disabled?: boolean;
		placeholder?: string;
		error?: string;
		onchange?: (value: string) => void;
	};

	let {
		label,
		name,
		options,
		value = $bindable(''),
		disabled = false,
		placeholder,
		error,
		onchange
	}: Props = $props();

	const change = (event: Event) => {
		value = (event.currentTarget as HTMLSelectElement).value;
		onchange?.(value);
	};
</script>

<div class="flex flex-col gap-1">
	<label class="mb-2 text-xs font-extrabold tracking-widest uppercase" for={name}>{label}</label>
	<select
		id={name}
		{name}
		{value}
		{disabled}
		onchange={change}
		aria-invalid={error ? 'true' : undefined}
		aria-describedby={error ? `${name}-error` : undefined}
		class="w-full rounded-xl border border-white/22 bg-white/5 focus:border-primary/45 focus:ring-primary/45 focus:outline-none disabled:opacity-50"
	>
		{#if placeholder}<option value="" class="bg-[#1a1a24]">{placeholder}</option>{/if}
		{#each options as option (option.value)}
			<option value={option.value} class="bg-[#1a1a24]">{option.label}</option>
		{/each}
	</select>
	{#if error}
		<span id={`${name}-error`} class="text-sm text-red-300">{error}</span>
	{/if}
</div>
