<script lang="ts">
	import * as m from '$lib/paraglide/messages';
	import Section from '$lib/components/Section.svelte';

	let { entries } = $props();

	const getTransactionTypeLabel = (type: string) => {
		switch (type) {
			case 'purchase':
				return m.buddypoints_type_purchase();
			case 'referral':
				return m.buddypoints_type_referral();
			case 'monthly_reset':
				return m.buddypoints_type_monthly();
			case 'usage':
				return m.buddypoints_type_usage();
			default:
				return type;
		}
	};

	const formatDate = (dateString: string) => {
		const date = new Date(dateString);
		return date.toLocaleDateString('de-DE', {
			year: 'numeric',
			month: '2-digit',
			day: '2-digit'
		});
	};
</script>

<Section>
	<h2 class="text-2xl font-semibold p-4 mb-6 text-[#e9e9f1]">{m.buddypoints_transaction_history()}</h2>
	<div class="overflow-hidden  bg-white/6 shadow-[0_20px_60px_rgba(0,0,0,0.55)]">
		<table class="w-full">
			<thead class="bg-white/4">
			<tr>
				<th class="px-6 py-3 text-left text-xs font-extrabold uppercase tracking-widest text-[#b7b7c8]">
					{m.buddypoints_date()}
				</th>
				<th class="px-6 py-3 text-left text-xs font-extrabold uppercase tracking-widest text-[#b7b7c8]">
					{m.buddypoints_type()}
				</th>
				<th class="px-6 py-3 text-left text-xs font-extrabold uppercase tracking-widest text-[#b7b7c8]">
					{m.buddypoints_description()}
				</th>
				<th class="px-6 py-3 text-right text-xs font-extrabold uppercase tracking-widest text-[#b7b7c8]">
					{m.buddypoints_amount()}
				</th>
				<th class="px-6 py-3 text-right text-xs font-extrabold uppercase tracking-widest text-[#b7b7c8]">
					{m.buddypoints_balance_after()}
				</th>
			</tr>
			</thead>
			<tbody class="border-t border-white/8">
			{#each entries as entry (entry.id)}
				<tr class="border-b border-white/8">
					<td class="px-6 py-4 whitespace-nowrap text-sm text-[#e9e9f1]">{formatDate(entry.created_at)}</td>
					<td class="px-6 py-4 whitespace-nowrap text-sm text-[#e9e9f1]">
							<span
								class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-white/6 text-[#e9e9f1]">
								{getTransactionTypeLabel(entry.transaction_type)}
							</span>
					</td>
					<td class="px-6 py-4 whitespace-nowrap text-sm text-[#e9e9f1]">{entry.description}</td>
					<td
						class="px-6 py-4 whitespace-nowrap text-sm text-right {entry.amount > 0 ? 'text-[#42f59b]' : 'text-[#ff4d6d]'}">
						{entry.amount > 0 ? '+' : ''}{entry.amount}
					</td>
					<td class="px-6 py-4 whitespace-nowrap text-sm text-right text-[#e9e9f1]">
						{entry.lifetime_points_after}
					</td>
				</tr>
			{/each}
			</tbody>
		</table>
	</div>
</Section>
