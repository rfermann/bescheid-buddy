<script lang="ts">
	import { page } from '$app/state';
	import { ArrowLeft, CalendarDays, Download, FileText } from '@lucide/svelte';
	import * as m from '$lib/paraglide/messages.js';
	import { getLocale } from '$lib/paraglide/runtime';
	import PageHeading from '$lib/components/PageHeading.svelte';
	import Section from '$lib/components/Section.svelte';
	import { getPaperLetterDeliveryDetailQuery } from '../deliveries.remote';

	const detail = $derived(
		await getPaperLetterDeliveryDetailQuery({ id: page.params.id as string })
	);
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
	const status = $derived(
		statuses[detail.delivery.status as keyof typeof statuses] ?? {
			label: m.paper_deliveries_status_unknown,
			class: 'bg-amber-400/15 text-amber-300'
		}
	);
</script>

<svelte:head>
	<title>{detail.delivery.topic || m.paper_deliveries_topic_missing()} - BescheidBuddy</title>
</svelte:head>

<a
	href="/meine-papierbriefe"
	class="mb-4 inline-flex items-center gap-2 text-sm font-semibold text-primary transition-colors hover:text-primary/80"
>
	<ArrowLeft size={16} aria-hidden="true" />
	{m.paper_delivery_detail_back()}
</a>

<PageHeading title={detail.delivery.topic || m.paper_deliveries_topic_missing()} />

<Section>
	<div class="grid gap-6">
		<article class="rounded-[14px] border border-white/8 bg-white/[0.025] p-5">
			<div class="mb-5 flex items-center justify-between gap-3">
				<h2 class="flex items-center gap-2 font-semibold text-white">
					<FileText size={18} class="text-primary" aria-hidden="true" />
					{m.paper_delivery_detail_heading()}
				</h2>
				<span class={`rounded-full px-2 py-0.5 text-xs font-semibold ${status.class}`}>
					{status.label()}
				</span>
			</div>

			<dl class="space-y-4 text-sm">
				<div>
					<dt class="text-[#8f8fa3]">{m.paper_delivery_detail_recipient()}</dt>
					<dd class="mt-1 break-words text-[#e9e9f1]">
						{detail.delivery.recipient || m.paper_deliveries_recipient_missing()}
					</dd>
				</div>
				<div>
					<dt class="text-[#8f8fa3]">{m.paper_delivery_detail_dates()}</dt>
					<dd class="mt-1 space-y-1 text-[#e9e9f1]">
						<p class="flex items-center gap-2">
							<CalendarDays size={14} class="text-primary" aria-hidden="true" />
							{m.paper_delivery_detail_created()}: {dateFormatter.format(detail.delivery.createdAt)}
						</p>
						<p>
							{m.paper_deliveries_updated_at()}: {dateFormatter.format(detail.delivery.updatedAt)}
						</p>
						{#if detail.delivery.sentAt}
							<p>{m.paper_deliveries_sent_at()}: {dateFormatter.format(detail.delivery.sentAt)}</p>
						{/if}
					</dd>
				</div>
				<div class="grid grid-cols-2 gap-4">
					<div>
						<dt class="text-[#8f8fa3]">{m.paper_delivery_detail_pages()}</dt>
						<dd class="mt-1 text-[#e9e9f1]">{detail.delivery.pageCount}</dd>
					</div>
					<div>
						<dt class="text-[#8f8fa3]">{m.paper_delivery_detail_cost()}</dt>
						<dd class="mt-1 text-[#e9e9f1]">
							{detail.delivery.buddyPointsCost} BuddyPoints
						</dd>
					</div>
				</div>
				<div>
					<dt class="text-[#8f8fa3]">{m.paper_delivery_detail_print_settings()}</dt>
					<dd class="mt-1 text-[#e9e9f1]">
						{detail.delivery.deliveryProduct} · {detail.delivery.printMode} · {detail.delivery
							.printSpectrum}
						· {detail.delivery.addressPosition}
					</dd>
				</div>
			</dl>
		</article>

		<article class="min-w-0 rounded-[14px] border border-white/8 bg-white/[0.025] p-4">
			<div class="mb-4 flex items-center justify-between gap-3">
				<h2 class="font-semibold text-white">{m.paper_delivery_detail_pdf()}</h2>
				<a
					href={detail.pdfURLs.downloadUrl}
					class="inline-flex items-center gap-2 rounded-xl border border-primary/45 bg-primary/15 px-3 py-2 text-sm font-bold text-primary transition-colors hover:border-primary/60"
				>
					<Download size={16} aria-hidden="true" />
					{m.paper_delivery_detail_download()}
				</a>
			</div>
			<iframe
				src={detail.pdfURLs.inlineUrl}
				title={m.paper_delivery_detail_pdf_title()}
				class="h-[70vh] min-h-[32rem] w-full rounded-lg bg-white"
			></iframe>
		</article>
	</div>
</Section>
