<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';
	import AnalysisCard from './AnalysisCard.svelte';
	import AddressCard from './AddressCard.svelte';
	import type { Address, Metadata } from '../../../generated';

	let { result }: { result: { text?: string; metadata?: Metadata } | null } = $props();

	let receiver = $derived<Address>(result?.metadata?.receiver ?? {});
	let sender = $derived<Address>(result?.metadata?.sender ?? {});
	let details = $derived(result?.metadata?.part_2_ai_analysis);
	let subject = $derived(result?.metadata?.subject);
	let dates = $derived(result?.metadata?.dates);
</script>

{#if result}
	<div class="mt-8 flex flex-col gap-4">
		<h3 class="text-lg font-semibold text-white">{m.analyse_result_overview()}</h3>
		<div class="grid gap-4 sm:grid-cols-2">
			<AddressCard title={m.analyse_result_receiver()} {...receiver} />
			<AddressCard title={m.analyse_result_sender()} {...sender} />
		</div>
		<div>
			<div class="grid gap-4 sm:grid-cols-2">
				<AnalysisCard
					title={m.analyse_result_subject()}
					text={subject?.localized || subject?.original}
				/>
				<AnalysisCard title={m.analyse_result_letter_type()} text={result.metadata?.letter_type} />
				<AnalysisCard title={m.analyse_result_importance()} text={result.metadata?.importance} />
				<AnalysisCard title={m.analyse_result_letter_date()} text={dates?.letter_date} />
				<AnalysisCard title={m.analyse_result_deadline_date()} text={dates?.deadline_date} />
				<AnalysisCard title={m.analyse_result_summary()} text={details?.summary} />
				<AnalysisCard title={m.analyse_result_deadlines()} items={details?.deadlines} />
				<AnalysisCard
					title={m.analyse_result_next_steps()}
					items={details?.recommended_next_steps}
				/>
				<AnalysisCard
					title={m.analyse_result_documents_requested()}
					items={details?.documents_requested}
				/>
				<AnalysisCard
					title={m.analyse_result_possible_consequences()}
					items={details?.possible_consequences}
				/>
			</div>
		</div>
		<div class="rounded-lg border border-white/10 bg-white/5 p-4">
			<h3 class="text-lg font-semibold text-white">{m.analyse_result_translation_title()}</h3>
			{#if result.metadata?.topic}
				<p class="mt-1 text-sm font-medium text-[#b7b7c8]">{result.metadata.topic}</p>
			{/if}
			<div class="mt-2 text-sm whitespace-pre-wrap text-[#e5e5ef]">
				{result.metadata?.part_1_full_translation ||
				result.text ||
				m.analyse_result_no_information()}
			</div>
		</div>
		<div
			id="resultDisclaimer"
			class="rounded-xl border border-[rgba(255,196,0,0.35)] bg-[rgba(255,196,0,0.07)] px-3 py-2.5 text-xs leading-[1.45] text-[rgba(255,228,160,0.96)]"
		>
			{m.analyse_result_disclaimer()}
		</div>
	</div>
{/if}
