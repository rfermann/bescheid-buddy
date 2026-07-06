<script lang="ts">
	import Section from '$lib/components/Section.svelte';
	import { getSupportCasesQuery } from './support.remote.ts';

	const supportCases = $derived(await getSupportCasesQuery({}));

	const formatDate = (dateString: Date | null): string => {
		if (!dateString) {
			return '-';
		}

		return dateString.toLocaleDateString();
	};

	const getCategoryDisplay = (category: string): string => {
		const categoryMap = {
			general: 'General',
			technical: 'Technical',
			billing: 'Billing',
			feature_request: 'Feature Request'
		};
		return categoryMap[category as keyof typeof categoryMap] || category;
	};

	const getStatusDisplay = (status: string): string => {
		const statusMap = {
			open: 'Open',
			in_progress: 'In Progress',
			resolved: 'Resolved',
			closed: 'Closed'
		};
		return statusMap[status as keyof typeof statusMap] || status;
	};

	const getStatusClass = (status: string): string => {
		const statusClasses = {
			open: 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30',
			in_progress: 'bg-blue-500/20 text-blue-400 border-blue-500/30',
			resolved: 'bg-green-500/20 text-green-400 border-green-500/30',
			closed: 'bg-gray-500/20 text-gray-400 border-gray-500/30'
		};
		return statusClasses[status as keyof typeof statusClasses] || '';
	};

	const handleCaseClick = (caseId: string) => {
		// TODO: Navigate to case detail page or handle case selection
		console.log('Case selected:', caseId);
	};
</script>

<Section>
	<div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6 mb-6">
		<div>
			<h2 class="text-2xl font-semibold mb-2 text-[#e9e9f1]">Support Cases</h2>
			<p class="text-sm text-[#b7b7c8]">An overview of your support cases</p>
		</div>
	</div>
	<div class="overflow-x-auto">
		<table class="w-full text-sm">
			<thead>
			<tr class="border-b border-white/10">
				<th class="text-left py-3 px-4 text-[#b7b7c8] font-medium">Date Created</th>
				<th class="text-left py-3 px-4 text-[#b7b7c8] font-medium">Category</th>
				<th class="text-left py-3 px-4 text-[#b7b7c8] font-medium">Status</th>
				<th class="text-left py-3 px-4 text-[#b7b7c8] font-medium">Date Closed</th>
				<th class="text-left py-3 px-4 text-[#b7b7c8] font-medium">Actions</th>
			</tr>
			</thead>
			<tbody>
			{#each supportCases.cases || [] as supportCase (supportCase.id)}
				<tr
					class="border-b border-white/5 hover:bg-white/5 cursor-pointer transition-colors"
					onclick={() => handleCaseClick(supportCase.id)}
				>
					<td class="py-3 px-4 text-[#e9e9f1]">{formatDate(supportCase.created_at)}</td>
					<td class="py-3 px-4 text-[#e9e9f1]">{getCategoryDisplay(supportCase.category)}</td>
					<td class="py-3 px-4">
							<span
								class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium border {getStatusClass(supportCase.status)}">
								{getStatusDisplay(supportCase.status)}
							</span>
					</td>
					<td class="py-3 px-4 text-[#e9e9f1]">{formatDate(supportCase.resolved_at)}</td>
					<td class="py-3 px-4">
						<a 
							href="/support/{supportCase.id}"
							class="inline-flex items-center px-3 py-1 rounded-2xl border font-bold border-primary/45 bg-primary/15 text-primary text-xs hover:border-primary/60 transition-colors"
						>
							Details
						</a>
					</td>
				</tr>
			{:else}
				<tr>
					<td colspan="5" class="py-8 text-center text-[#b7b7c8]">
						No support cases found
					</td>
				</tr>
			{/each}
			</tbody>
		</table>
	</div>
</Section>
