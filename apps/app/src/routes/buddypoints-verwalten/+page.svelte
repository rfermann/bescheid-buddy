<script lang="ts">
	import * as m from '$lib/paraglide/messages';
	import PriceCard from './PriceCard.svelte';
	import PurchaseModal from './PurchaseModal.svelte';
	import TransactionHistory from './TransactionHistory.svelte';
	import {
		getBuddyPackagesQuery,
		getBuddyPointsQuery,
		getBuddyPointsLedgerQuery,
		getCheckoutSessionUrlQuery
	} from './points.remote.ts';
	import { onMount, tick } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { replaceState } from '$app/navigation';
	import type { PageData } from './$types';
	import Section from '$lib/components/Section.svelte';

	let { data }: { data: PageData } = $props();

	const packages = $derived(await getBuddyPackagesQuery({}));
	const buddyPoints = $derived(await getBuddyPointsQuery({}));
	const ledger = $derived(await getBuddyPointsLedgerQuery({}));

	const useModal = true;
	let modalOpen = $derived.by(() => data.purchaseSuccess && useModal);

	onMount(async () => {
		await tick();

		if (data.purchaseSuccess && data.packageName) {
			if (!useModal) {
				toast.success(m.buddypoints_purchase_success_toast({ name: data.packageName }));
			}
		} else if (data.purchaseCanceled) {
			toast.error(m.buddypoints_purchase_canceled());
		}
		replaceState(window.location.pathname, {});
	});

	const closeModal = () => {
		modalOpen = false;
	};

	let selectedPlan = $state<string | null>('Buddy Package');

	const selectPlan = (planName: string) => {
		selectedPlan = planName;
	};

	let isLoading = $state(false);
	const handleBuy = async () => {
		isLoading = true;
		const pkg = packages?.packages?.find((p) => p.name === selectedPlan);
		if (!pkg) {
			return;
		}
		const url = await getCheckoutSessionUrlQuery({ name: pkg.name, stripePriceId: pkg.stripe_price_id }).run();
		if (!url) {
			return;
		}
		window.location.href = url;
	};

</script>

<svelte:head>
	<title>{m.buddypoints_page_title()} - BescheidBuddy</title>
</svelte:head>

<div
	class="container mx-auto px-4 py-8 max-w-6xl min-h-screen text-[#e9e9f1]">
	<h1 class="text-3xl font-bold mb-8 text-[#e9e9f1]">{m.buddypoints_page_title()}</h1>

	<Section>
		<div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6 mb-6">
			<div>
				<h2 class="text-2xl font-semibold mb-2 text-[#e9e9f1]">{m.buddypoints_overview_title()}</h2>
				<p class="text-sm text-[#b7b7c8]">{m.buddypoints_overview_subtitle()}</p>
			</div>
			<div class="text-right">
				<p class="text-5xl font-bold text-primary">{buddyPoints.monthly_points + buddyPoints.lifetime_points}</p>
				<p class="text-sm text-[#b7b7c8]">{m.buddypoints_points()} {m.buddypoints_points_available()}</p>
			</div>
		</div>
		<div class="grid grid-cols-1 md:grid-cols-2 gap-3">
			<div class="border-primary-bg bg-primary-bg rounded-[14px] p-4">
				<p class="text-sm mb-1 text-[#b7b7c8]">{m.buddypoints_purchased()}</p>
				<p class="text-2xl font-bold text-primary">{buddyPoints.lifetime_points}</p>
				<p class="text-xs mt-1 text-primary/70">{m.buddypoints_purchased_subtitle()}</p>
			</div>
			<div class="border-purple-500/15 bg-purple-500/6 rounded-[14px] p-4">
				<p class="text-sm mb-1 text-[#b7b7c8]">{m.buddypoints_monthly()}</p>
				<p class="text-2xl font-bold text-[#a78bfa]">{buddyPoints.monthly_points}</p>
				<p class="text-xs mt-1 text-purple-400/70">{m.buddypoints_monthly_subtitle()}</p>
			</div>
		</div>
	</Section>

	<Section>
		<div class="flex justify-between items-center mb-1.5">
			<h3 class="text-xl font-bold text-[#e9e9f1]">{m.buddypoints_buy_title()}</h3>
		</div>
		<p class="text-[#b7b7c8] leading-relaxed mb-3">{m.buddypoints_buy_subtitle()}</p>
		<div class="grid grid-cols-1 md:grid-cols-3 gap-3">
			{#each packages.packages as pkg, index}
				<PriceCard
					{...pkg}
					active={pkg.name === selectedPlan}
					onclick={(name) => {selectPlan(name)}}
					badge={index === 0 ? m.buddypoints_package_starter_badge() : index === 1 ? m.buddypoints_package_popular_badge() : m.buddypoints_package_value_badge()}
					note={index === 0 ? m.buddypoints_package_starter_note() : index === 1 ? m.buddypoints_package_popular_note() : m.buddypoints_package_value_note()}
				/>
			{/each}
		</div>
		<div class="text-xs text-[rgba(183,183,200,0.9)] mt-3"></div>
		<div class="flex justify-end mt-3.5">

			<button
				class="inline-flex items-center gap-2.5 py-3 px-3.5 border border-solid border-primary/40 rounded-[14px] text-primary bg-primary/[0.161] text-[13.333px] font-black tracking-[0.267px] leading-normal transition-all duration-150 cursor-pointer whitespace-nowrap select-none hover:-translate-y-px text-center disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:translate-y-0"
				onclick={() => handleBuy()}
				disabled={isLoading}
				type="button">
				{#if isLoading}
					<svg
						class="animate-spin h-4 w-4"
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24">
						<circle
							class="opacity-25"
							cx="12"
							cy="12"
							r="10"
							stroke="currentColor"
							stroke-width="4"></circle>
						<path
							class="opacity-75"
							fill="currentColor"
							d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>
				{/if}
				{m.buddypoints_create_purchase_request()}
			</button>
		</div>
	</Section>

	<TransactionHistory entries={ledger?.entries || []} />
</div>

{#if modalOpen && data.packageName && data.packagePoints}
	<PurchaseModal
		packageName={data.packageName}
		points={data.packagePoints}
		onClose={closeModal}
	/>
{/if}