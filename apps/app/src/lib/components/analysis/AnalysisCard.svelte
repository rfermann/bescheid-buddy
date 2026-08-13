<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';

	let { title, text, items }: { title: string; text?: string; items?: string[] | null } = $props();
</script>

<div class="rounded-lg border border-white/10 bg-white/5 p-4">
	<h3 class="text-sm font-medium text-[#b7b7c8]">{title}</h3>
	{#if text}
		<p class="mt-1 text-sm whitespace-pre-wrap text-[#e5e5ef]">{text}</p>
	{/if}
	{#if items?.length}
		<ul class="mt-2 list-disc space-y-1 pl-5 text-sm text-[#e5e5ef]">
			{#each items as item}
				<li>{item}</li>
			{/each}
		</ul>
	{:else if !text}
		<p class="mt-1 text-sm text-[#8f8fa3]">{m.analyse_result_no_information()}</p>
	{/if}
</div>
