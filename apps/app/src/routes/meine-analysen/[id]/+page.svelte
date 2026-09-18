<script lang="ts">
	import { page } from '$app/state';
	import { ArrowLeft, CalendarDays, FileText, TriangleAlert } from '@lucide/svelte';
	import * as m from '$lib/paraglide/messages.js';
	import { getLocale } from '$lib/paraglide/runtime';
	import PageHeading from '$lib/components/PageHeading.svelte';
	import Result from '$lib/components/analysis/Result.svelte';
	import Section from '$lib/components/Section.svelte';
	import PrimaryLink from '$lib/components/PrimaryLink.svelte';
	import { getAnalysisQuery } from '../analyses.remote';

	const detail = $derived(await getAnalysisQuery({ id: page.params.id as string }));
	const locale = getLocale();
	const dateFormatter = new Intl.DateTimeFormat(locale === 'ua' ? 'uk' : locale, {
		dateStyle: 'medium',
		timeStyle: 'short'
	});
</script>

<svelte:head>
	<title>{detail.analysisCase.title || m.analyses_untitled()} - BescheidBuddy</title>
</svelte:head>

<a
	href="/meine-analysen"
	class="mb-4 inline-flex items-center gap-2 text-sm font-semibold text-primary transition-colors hover:text-primary/80"
>
	<ArrowLeft size={16} aria-hidden="true" />
	{m.analysis_detail_back()}
</a>

<PageHeading title={detail.analysisCase.title || m.analyses_untitled()} />

<Section>
	<div class="flex flex-wrap gap-x-6 gap-y-2 border-b border-white/8 pb-4 text-sm text-[#b7b7c8]">
		<span class="inline-flex items-center gap-2">
			<CalendarDays size={16} class="text-primary" aria-hidden="true" />
			{m.analyses_created()}
			{dateFormatter.format(detail.analysisCase.createdAt)}
		</span>
	</div>

	<div class="mt-6 flex flex-col gap-6">
		{#each detail.documents as document (document.documentId)}
			<article class="rounded-xl border border-white/10 bg-white/[0.025] p-4">
				<h2 class="flex items-center gap-2 font-semibold text-white">
					<FileText size={18} class="text-primary" aria-hidden="true" />
					{document.fileName}
				</h2>

				{#if document.analysis.status === 'completed'}
					{#if document.analysis.text || document.analysis.metadata}
						<Result result={document.analysis} />
					{/if}
					<div class="mt-5 flex justify-end">
						<PrimaryLink
							href={`/papierbrief?analyse=${detail.analysisCase.id}&dokument=${document.documentId}`}
							text={m.analysis_detail_create_paper_letter()}
						/>
					</div>
				{:else if document.analysis.status === 'failed'}
					<div
						role="alert"
						class="mt-4 flex gap-3 rounded-xl border border-red-400/25 bg-red-400/10 p-4 text-red-200"
					>
						<TriangleAlert class="mt-0.5 shrink-0" size={20} aria-hidden="true" />
						<div>
							<h3 class="font-semibold">{m.analysis_detail_failed_title()}</h3>
							<p class="mt-1 text-sm text-red-200/80">
								{m.analysis_detail_failed_description()}
							</p>
						</div>
					</div>
				{:else}
					<div class="mt-4 rounded-xl border border-white/10 bg-white/5 p-5">
						<h3 class="font-semibold text-white">{m.analysis_detail_pending_title()}</h3>
						<p class="mt-1 text-sm text-[#b7b7c8]">
							{m.analysis_detail_pending_description()}
						</p>
					</div>
				{/if}
			</article>
		{:else}
			<div class="rounded-xl border border-white/10 bg-white/5 p-5">
				<h2 class="font-semibold text-white">{m.analysis_detail_pending_title()}</h2>
				<p class="mt-1 text-sm text-[#b7b7c8]">
					{m.analysis_detail_pending_description()}
				</p>
			</div>
		{/each}
	</div>
</Section>
