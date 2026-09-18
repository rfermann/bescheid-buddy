<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';

	let {
		models,
		value = $bindable(''),
		disabled = false,
		id = 'analyse-model',
		label = m.analyse_model_label()
	}: {
		models: Array<{ id: string; name: string }>;
		value?: string;
		disabled?: boolean;
		id?: string;
		label?: string;
	} = $props();
</script>

<div class="mt-6 flex flex-col gap-2">
	<label for={id} class="text-sm font-medium text-[#b7b7c8]">
		{label}
	</label>
	<select
		{id}
		bind:value
		disabled={disabled || models.length === 0}
		class="w-full max-w-xs rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-sm text-white focus:border-white/30 focus:outline-none disabled:opacity-50"
	>
		{#each models as model (model.id)}
			<option value={model.id} class="bg-[#1a1a24] text-white">{model.name}</option>
		{/each}
	</select>
</div>
