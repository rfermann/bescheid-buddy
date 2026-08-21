<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { FieldErrors } from '$lib/components/ui/field-errors';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { createModelForm, getModelProviders, getModels } from './models.remote';

	const providers = $derived(await getModelProviders());
	let open = $state(false);

	createModelForm.enhance(async ({ submit }) => {
		await submit();
		if (createModelForm.result?.success) {
			createModelForm.fields.set({
				name: '',
				provider_model_id: '',
				provider: providers[0] ?? '',
				sort_order: 0
			});
			open = false;
			void getModels().refresh();
		}
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger>
		{#snippet child({ props })}
			<Button {...props}>Add model</Button>
		{/snippet}
	</Dialog.Trigger>
	<Dialog.Content class="w-[calc(100%-2rem)] max-w-lg sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>Add model</Dialog.Title>
			<Dialog.Description>Add a model supported by one of the configured providers</Dialog.Description>
		</Dialog.Header>
		<form {...createModelForm} class="grid gap-4">
			<div class="space-y-2">
				<Label for="model-name">Name</Label>
				<Input
					id="model-name"
					placeholder="GPT-5.6 Luna"
					maxlength={255}
					{...createModelForm.fields.name.as('text')}
				/>
				<FieldErrors issues={createModelForm.fields.name.issues()} />
			</div>
			<div class="space-y-2">
				<Label for="provider-model-id">Provider model ID</Label>
				<Input
					id="provider-model-id"
					placeholder="gpt-5.6-luna"
					maxlength={255}
					{...createModelForm.fields.provider_model_id.as('text')}
				/>
				<FieldErrors issues={createModelForm.fields.provider_model_id.issues()} />
			</div>
			<div class="space-y-2">
				<Label for="model-provider">Provider</Label>
				<select
					id="model-provider"
					class="border-input dark:bg-input/30 focus-visible:border-ring focus-visible:ring-ring/50 h-9 w-full rounded-md border bg-transparent px-2.5 py-1 text-sm shadow-xs outline-none focus-visible:ring-3"
					{...createModelForm.fields.provider.as('select', providers[0] ?? '')}
				>
					<option value="" disabled>Select a provider</option>
					{#each providers as provider (provider)}
						<option value={provider}>{provider}</option>
					{/each}
				</select>
				<FieldErrors issues={createModelForm.fields.provider.issues()} />
			</div>
			<div class="space-y-2">
				<Label for="sort-order">Sort order</Label>
				<Input
					id="sort-order"
					step="1"
					{...createModelForm.fields.sort_order.as('number', 0)}
				/>
				<FieldErrors issues={createModelForm.fields.sort_order.issues()} />
			</div>
			{#if createModelForm.result?.success === false}
				<p class="text-destructive text-sm" role="alert">{createModelForm.result.message}</p>
			{/if}
			<Dialog.Footer>
				<Button type="submit" disabled={!!createModelForm.pending}>
					{createModelForm.pending ? 'Adding...' : 'Add model'}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
