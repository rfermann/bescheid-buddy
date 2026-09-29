<script lang="ts">
	import { CalendarDays, FileText, Plus } from '@lucide/svelte';
	import { resolve } from '$app/paths';
	import * as m from '$lib/paraglide/messages.js';
	import { getLocale } from '$lib/paraglide/runtime';
	import PageHeading from '$lib/components/PageHeading.svelte';
	import Section from '$lib/components/Section.svelte';
	import { getPaperLetterDeliveriesQuery } from './deliveries.remote';

	const deliveries = $derived(await getPaperLetterDeliveriesQuery({}));
	const locale = getLocale();
	const dateFormatter = new Intl.DateTimeFormat(locale === 'ua' ? 'uk' : locale, {
		dateStyle: 'medium',
		timeStyle: 'short'
	});

	const statuses = {
		pending: { label: m.paper_deliveries_status_pending, class: 'bg-white/10 text-[#b7b7c8]' },
		processing: {
			label: m.paper_deliveries_status_processing,
			class: 'bg-blue-400/15 text-blue-300'
		},
		sent: { label: m.paper_deliveries_status_sent, class: 'bg-emerald-400/15 text-emerald-300' },
		failed: { label: m.paper_deliveries_status_failed, class: 'bg-red-400/15 text-red-300' },
		manual_review: {
			label: m.paper_deliveries_status_manual_review,
			class: 'bg-amber-400/15 text-amber-300'
		}
	} as const;

	const statusFor = (status: string) =>
		statuses[status as keyof typeof statuses] ?? {
			label: m.paper_deliveries_status_unknown,
			class: 'bg-amber-400/15 text-amber-300'
		};
</script>

<svelte:head>
	<title>{m.paper_deliveries_page_title()} - BescheidBuddy</title>
</svelte:head>

<PageHeading title={m.paper_deliveries_page_title()} />

<Section>
	<div class="mx-auto max-w-4xl">
		<div class="mb-6 flex justify-end">
			<a
				href={resolve('/papierbrief')}
				class="inline-flex items-center gap-2 rounded-2xl border border-primary/45 bg-primary/15 px-3.5 py-2.5 text-sm font-bold text-primary transition-colors hover:border-primary/60"
			>
				<Plus size={17} aria-hidden="true" />
				{m.paper_deliveries_new()}
			</a>
		</div>

		{#if deliveries.length}
			<div class="overflow-hidden rounded-[14px] border border-white/8 bg-white/[0.025]">
				{#each deliveries as delivery (delivery.deliveryId)}
					{@const status = statusFor(delivery.status)}
					<div
						class="grid gap-3 border-b border-white/8 px-4 py-4 last:border-b-0 sm:grid-cols-[auto_minmax(0,1fr)_auto] sm:items-center sm:gap-4"
					>
						<div
							class="flex size-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary"
						>
							<FileText size={20} aria-hidden="true" />
						</div>

						<div class="min-w-0">
							<div class="flex flex-wrap items-center gap-2">
								<h3 class="text-base font-semibold break-words text-[#e9e9f1]">
									{delivery.topic || m.paper_deliveries_topic_missing()}
								</h3>
								<span class={`rounded-full px-2 py-0.5 text-xs font-semibold ${status.class}`}>
									{status.label()}
								</span>
							</div>
							<p class="mt-1 text-sm break-words text-[#b7b7c8]">
								{delivery.recipient || m.paper_deliveries_recipient_missing()}
							</p>
							<p class="mt-2 flex items-center gap-1.5 text-xs text-[#8f8fa3] sm:hidden">
								<CalendarDays size={14} aria-hidden="true" />
								{delivery.sentAt ? m.paper_deliveries_sent_at() : m.paper_deliveries_updated_at()}
								<span class="text-[#b7b7c8]"
									>{dateFormatter.format(delivery.sentAt ?? delivery.updatedAt)}</span
								>
							</p>
						</div>

						<div class="flex items-center justify-between gap-4 sm:block sm:text-right">
							<div class="text-xs">
								<p class="mb-1 hidden text-[#8f8fa3] sm:block">
									{delivery.sentAt ? m.paper_deliveries_sent_at() : m.paper_deliveries_updated_at()}
								</p>
								<p class="hidden text-[#b7b7c8] sm:block">
									{dateFormatter.format(delivery.sentAt ?? delivery.updatedAt)}
								</p>
							</div>
							{#if delivery.caseId}
								<a
									href={resolve('/meine-analysen/[id]', { id: delivery.caseId })}
									class="text-sm font-semibold text-primary hover:underline"
								>
									{m.paper_deliveries_view_analysis()}
								</a>
							{/if}
						</div>
					</div>
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
				<h3 class="mb-2 text-lg font-semibold text-[#e9e9f1]">
					{m.paper_deliveries_empty_title()}
				</h3>
				<p class="max-w-md text-sm leading-relaxed text-[#b7b7c8]">
					{m.paper_deliveries_empty_description()}
				</p>
				<a
					href={resolve('/papierbrief')}
					class="mt-5 text-sm font-semibold text-primary hover:underline"
				>
					{m.paper_deliveries_new()}
				</a>
			</div>
		{/if}
	</div>
</Section>
