<script lang="ts">
	import { CalendarDays, FileText, Plus } from '@lucide/svelte';
	import * as m from '$lib/paraglide/messages.js';
	import { getLocale } from '$lib/paraglide/runtime';
	import PageHeading from '$lib/components/PageHeading.svelte';
	import Section from '$lib/components/Section.svelte';
	import { getCasesQuery } from './analyses.remote';

	const cases = $derived(await getCasesQuery({}));
	const locale = getLocale();
	const dateFormatter = new Intl.DateTimeFormat(locale === 'ua' ? 'uk' : locale, {
		dateStyle: 'medium',
		timeStyle: 'short'
	});
</script>

<svelte:head>
	<title>{m.analyses_page_title()} - BescheidBuddy</title>
</svelte:head>

<PageHeading title={m.analyses_page_title()} />

<Section>
	<div class="mx-auto max-w-4xl">
		<div class="mb-6 flex justify-end">
			<a
				href="/"
				class="inline-flex items-center gap-2 rounded-2xl border border-primary/45 bg-primary/15 px-3.5 py-2.5 text-sm font-bold text-primary transition-colors hover:border-primary/60"
			>
				<Plus size={17} aria-hidden="true" />
				{m.analyses_new()}
			</a>
		</div>

		{#if cases.length}
			<div class="overflow-hidden rounded-[14px] border border-white/8 bg-white/[0.025]">
				{#each cases as analysis (analysis.id)}
					<a
						href="/meine-analysen/{analysis.id}"
						class="grid grid-cols-[auto_minmax(0,1fr)] items-center gap-4 border-b border-white/8 px-4 py-4 transition-colors last:border-b-0 hover:bg-primary/[0.045] focus-visible:bg-primary/[0.045] focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-primary/70 md:grid-cols-[auto_minmax(0,1fr)_auto]"
					>
						<div
							class="flex size-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary"
						>
							<FileText size={20} aria-hidden="true" />
						</div>

						<div class="min-w-0">
							<h3 class="truncate text-base font-semibold text-[#e9e9f1]">
								{analysis.title || m.analyses_untitled()}
							</h3>
							<p class="mt-1 flex items-center gap-1.5 text-xs text-[#8f8fa3]">
								<CalendarDays size={14} aria-hidden="true" />
								{m.analyses_created()}
								<span class="text-[#b7b7c8]">{dateFormatter.format(analysis.createdAt)}</span>
							</p>
							<p class="mt-1 text-xs text-[#8f8fa3] md:hidden">
								{m.analyses_updated()}
								<span class="text-[#b7b7c8]">{dateFormatter.format(analysis.updatedAt)}</span>
							</p>
						</div>

						<div class="hidden text-right text-xs md:block">
							<p class="mb-1 text-[#8f8fa3]">{m.analyses_updated()}</p>
							<p class="text-[#b7b7c8]">{dateFormatter.format(analysis.updatedAt)}</p>
						</div>
					</a>
				{/each}
			</div>
		{:else}
			<div
				class="flex flex-col items-center rounded-[14px] border border-dashed border-white/12 px-6 py-14 text-center"
			>
				<div
					class="mb-4 flex size-14 items-center justify-center rounded-2xl bg-primary/10 text-primary"
				>
					<FileText size={28} aria-hidden="true" />
				</div>
				<h3 class="mb-2 text-lg font-semibold text-[#e9e9f1]">{m.analyses_empty_title()}</h3>
				<p class="max-w-md text-sm leading-relaxed text-[#b7b7c8]">
					{m.analyses_empty_description()}
				</p>
			</div>
		{/if}
	</div>
</Section>
