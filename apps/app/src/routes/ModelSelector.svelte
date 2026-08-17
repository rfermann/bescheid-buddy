<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';

	let {
		models,
		value = $bindable(''),
		disabled = false
	}: {
		models: Array<{ id: string; name: string }>;
		value?: string;
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
		disabled={disabled || models.length === 0}
		class="w-full max-w-xs rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-sm text-white focus:border-white/30 focus:outline-none disabled:opacity-50"
	>
		{#each models as model (model.id)}
			<option value={model.id} class="bg-[#1a1a24] text-white">{model.name}</option>
		{/each}
	</select>
</div>
