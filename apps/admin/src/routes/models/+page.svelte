<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import AddModelDialog from './add-model-dialog.svelte';
	import ModelEnabledButton from './model-enabled-button.svelte';
	import { getModels } from './models.remote';

	const models = $derived(await getModels());
</script>

<div class="space-y-4">
	<Card>
		<CardHeader class="flex-row items-center justify-between">
			<div>
				<CardTitle>Models</CardTitle>
				<CardDescription>View the AI models configured for the application</CardDescription>
			</div>
			<AddModelDialog />
		</CardHeader>
		<CardContent>
			<div class="overflow-x-auto rounded-md border">
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>Name</Table.Head>
							<Table.Head>Provider model ID</Table.Head>
							<Table.Head>Provider</Table.Head>
							<Table.Head>State</Table.Head>
							<Table.Head class="text-right">Sort order</Table.Head>
							<Table.Head class="text-right">Actions</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each models as model (model.id)}
							<Table.Row>
								<Table.Cell class="font-medium">{model.name}</Table.Cell>
								<Table.Cell><code>{model.provider_model_id}</code></Table.Cell>
								<Table.Cell class="capitalize">{model.provider}</Table.Cell>
								<Table.Cell>
									<Badge variant={model.enabled ? 'default' : 'secondary'}>
										{model.enabled ? 'Enabled' : 'Disabled'}
									</Badge>
								</Table.Cell>
								<Table.Cell class="text-right">{model.sort_order}</Table.Cell>
								<Table.Cell class="text-right">
									<ModelEnabledButton id={model.id} enabled={model.enabled} />
								</Table.Cell>
							</Table.Row>
						{:else}
							<Table.Row>
								<Table.Cell colspan={6} class="h-24 text-center">No models found.</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
		</CardContent>
	</Card>
</div>
