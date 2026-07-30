<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';
	import AnalysisCard from './AnalysisCard.svelte';
	import AddressCard from './AddressCard.svelte';
	import type { Address, Metadata } from '../generated';

	let { result }: { result: { text?: string; metadata?: Metadata } | null } = $props();

	let receiver = $derived<Address>(result?.metadata?.receiver ?? {});
	let sender = $derived<Address>(result?.metadata?.sender ?? {});
	let details = $derived(result?.metadata?.part_2_ai_analysis);
</script>

{#if result}
	<div class="mt-8 flex flex-col gap-6">
		<div class="grid gap-4 sm:grid-cols-2">
			<AddressCard title={m.analyse_result_receiver()} {...receiver} />
			<AddressCard title={m.analyse_result_sender()} {...sender} />
		</div>
		<div class="rounded-lg border border-white/10 bg-white/5 p-4">
			<h3 class="text-lg font-semibold text-white">
				{result.metadata?.topic || m.analyse_result_translation_title()}
			</h3>
			<div class="mt-2 text-sm whitespace-pre-wrap text-[#e5e5ef]">
				{result.metadata?.part_1_full_translation || result.text}
			</div>
			<div
				id="resultDisclaimer"
				class="mt-2.5 rounded-xl border border-[rgba(255,196,0,0.35)] bg-[rgba(255,196,0,0.07)] px-3 py-2.5 text-xs leading-[1.45] text-[rgba(255,228,160,0.96)]"
			>
				{m.analyse_result_disclaimer()}
			</div>
		</div>
		<div>
			<h3 class="mb-2 text-lg font-semibold text-white">{m.analyse_result_overview()}</h3>
			<div class="grid gap-4 sm:grid-cols-2">
				{#if result.metadata?.letter_type}
					<AnalysisCard title={m.analyse_result_letter_type()} text={result.metadata.letter_type} />
				{/if}
				{#if result.metadata?.importance}
					<AnalysisCard title={m.analyse_result_importance()} text={result.metadata.importance} />
				{/if}
				{#if details?.summary}
					<AnalysisCard title={m.analyse_result_summary()} text={details.summary} />
				{/if}
				{#if details?.deadlines?.length}
					<AnalysisCard title={m.analyse_result_deadlines()} items={details.deadlines} />
				{/if}
				{#if details?.recommended_next_steps?.length}
					<AnalysisCard title={m.analyse_result_next_steps()} items={details.recommended_next_steps} />
				{/if}
				{#if details?.documents_requested?.length}
					<AnalysisCard
						title={m.analyse_result_documents_requested()}
						items={details.documents_requested}
					/>
				{/if}
				{#if details?.possible_consequences?.length}
					<AnalysisCard
						title={m.analyse_result_possible_consequences()}
						items={details.possible_consequences}
					/>
				{/if}
			</div>
		</div>
		{#if result.metadata}
			<div>
				<h3 class="mb-2 text-lg font-semibold text-white">
					{m.analyse_result_metadata_title()}
				</h3>
				<pre
					class="max-h-80 overflow-auto rounded-lg border border-white/10 bg-black/30 p-4 text-xs text-[#b7b7c8]">{JSON.stringify(
						result.metadata,
						null,
						2
					)}</pre>
			</div>
		{/if}
	</div>
{/if}
