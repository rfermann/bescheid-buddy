<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { getModels, setModelEnabledCommand } from './models.remote';

	let { id, enabled }: { id: string; enabled: boolean } = $props();
	let errorMessage = $state<string>();

	const toggleEnabled = async () => {
		errorMessage = undefined;
		const result = await setModelEnabledCommand({ id, enabled: !enabled });
		if (result.success) {
			void getModels().refresh();
		} else {
			errorMessage = result.message;
		}
	};
</script>

<div class="flex flex-col items-end gap-1">
	<Button
		variant={enabled ? 'outline' : 'default'}
		size="sm"
		disabled={!!setModelEnabledCommand.pending}
		onclick={toggleEnabled}
	>
		{enabled ? 'Disable' : 'Enable'}
	</Button>
	{#if errorMessage}
		<p class="text-destructive text-sm" role="alert">{errorMessage}</p>
	{/if}
</div>
