<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';

	type Model = 'gpt-4o-mini' | 'gemini-2.0-flash' | 'claude-haiku-4-5';

	const MODEL_OPTIONS: { value: Model; label: string }[] = [
		{ value: 'gpt-4o-mini', label: 'GPT-4o mini' },
		{ value: 'gemini-2.0-flash', label: 'Gemini 2.0 Flash' },
		{ value: 'claude-haiku-4-5', label: 'Claude Haiku 4.5' }
	];

	let {
		value = $bindable<Model>('gpt-4o-mini'),
		disabled = false
	}: {
		value?: Model;
		disabled?: boolean;
	} = $props();
</script>

<div class="mt-6 flex flex-col gap-2">
	<label for="analyse-model" class="text-sm font-medium text-[#b7b7c8]">
		{m.analyse_model_label()}
	</label>
	<select
		id="analyse-model"
		bind:value
		{disabled}
		class="w-full max-w-xs rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-sm text-white focus:border-white/30 focus:outline-none disabled:opacity-50"
	>
		{#each MODEL_OPTIONS as option (option.value)}
			<option value={option.value} class="bg-[#1a1a24] text-white">{option.label}</option>
		{/each}
	</select>
</div>
