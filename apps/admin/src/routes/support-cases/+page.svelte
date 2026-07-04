<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table';
	import * as Select from '$lib/components/ui/select';
	import * as Card from '$lib/components/ui/card';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import { getSupportCases } from './support-cases.remote';
	import { createSvelteTable, FlexRender } from '$lib/components/ui/data-table';
	import { getCoreRowModel } from '@tanstack/table-core';
	import { createColumns } from './support-cases.columns';

	const supportCases = $derived(await getSupportCases());
	let selectedCategory = $state('all');
	let selectedStatus = $state('all');

	const categories = ['all', 'general', 'technical', 'billing', 'feature_request'];
	const statuses = ['all', 'open', 'in_progress', 'resolved', 'closed'];

	const categoryLabel = $derived(
		categories.find((c) => c === selectedCategory) === 'all' ? 'All Categories' : selectedCategory
	);
	const statusLabel = $derived(
		statuses.find((s) => s === selectedStatus) === 'all' ? 'All Statuses' : selectedStatus.replace('_', ' ')
	);

	const filteredCases = $derived(
		(supportCases || []).filter(
			(c) =>
				(selectedCategory === 'all' || c.category === selectedCategory) &&
				(selectedStatus === 'all' || c.status === selectedStatus)
		)
	);

	const statusColors: Record<string, string> = {
		open: 'bg-red-500',
		in_progress: 'bg-yellow-500',
		resolved: 'bg-blue-500',
		closed: 'bg-gray-500'
	};

	const columns = createColumns({
		statusSnippet
	});

	const table = $derived(
		createSvelteTable({
			get data() {
				return filteredCases;
			},
			columns,
			getCoreRowModel: getCoreRowModel()
		})
	);
</script>

{#snippet statusSnippet(params: { status: string })}
	<span
		class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium text-white {statusColors[params.status] || 'bg-gray-500'}">
		{params.status.replace('_', ' ')}
	</span>
{/snippet}

<div class="space-y-4">
	<Card.Root>
		<Card.Header>
			<Card.Title>Support Cases</Card.Title>
			<Card.Description>View and manage all support cases</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
				<div class="flex flex-col gap-2 sm:flex-row sm:gap-4">
					<div class="flex flex-col gap-1">
						<label for="category-select" class="text-sm font-medium">Category</label>
						<Select.Root type="single" bind:value={selectedCategory}>
							<Select.Trigger class="w-45 capitalize">
								{categoryLabel.replace('_', ' ')}
							</Select.Trigger>
							<Select.Content>
								{#each categories as category (category)}
									<Select.Item value={category} class="capitalize">
										{category === 'all' ? 'All Categories' : category.replace('_', ' ')}
									</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</div>
					<div class="flex flex-col gap-1">
						<label for="status-select" class="text-sm font-medium">Status</label>
						<Select.Root type="single" bind:value={selectedStatus}>
							<Select.Trigger class="w-45">
								{statusLabel}
							</Select.Trigger>
							<Select.Content>
								{#each statuses as status (status)}
									<Select.Item value={status}>
										{status === 'all' ? 'All Statuses' : status.replace('_', ' ')}
									</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</div>
				</div>
				<Button onclick={() => getSupportCases().refresh()} variant="outline" size="default" class="mt-2 sm:mt-0">
					<RefreshCwIcon class="mr-2 h-4 w-4" />
					Refresh
				</Button>
			</div>

			<div class="rounded-md border">
				<Table.Root>
					<Table.Header>
						{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
							<Table.Row>
								{#each headerGroup.headers as header (header.id)}
									<Table.Head
										class={header.id === 'user_id' ? 'hidden md:table-cell' : header.id === 'created_at' || header.id === 'updated_at' ? 'hidden lg:table-cell' : ''}>
										{#if header.column.getCanSort()}
											<button
												onclick={() => header.column.toggleSorting()}
												class="flex items-center gap-2 hover:text-accent-foreground"
											>
												<FlexRender
													content={header.column.columnDef.header}
													context={header.getContext()}
												/>
												{#if header.column.getIsSorted() === 'asc'}
													↑
												{:else if header.column.getIsSorted() === 'desc'}
													↓
												{/if}
											</button>
										{:else}
											<FlexRender
												content={header.column.columnDef.header}
												context={header.getContext()}
											/>
										{/if}
									</Table.Head>
								{/each}
							</Table.Row>
						{/each}
					</Table.Header>
					<Table.Body>
						{#each table.getRowModel().rows as row (row.id)}
							<Table.Row>
								{#each row.getVisibleCells() as cell (cell.id)}
									<Table.Cell
										class={cell.column.id === 'user_id' ? 'hidden md:table-cell' : cell.column.id === 'created_at' || cell.column.id === 'updated_at' ? 'hidden lg:table-cell' : ''}>
										<FlexRender
											content={cell.column.columnDef.cell}
											context={cell.getContext()}
										/>
									</Table.Cell>
								{/each}
							</Table.Row>
						{:else}
							<Table.Row>
								<Table.Cell colspan={columns.length} class="h-24 text-center">
									No support cases found.
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>

			<div class="text-sm text-muted-foreground">
				Showing {filteredCases.length} of {Array.isArray(supportCases) ? supportCases.length : 0}
				support
				cases
			</div>
		</Card.Content>
	</Card.Root>
</div>
