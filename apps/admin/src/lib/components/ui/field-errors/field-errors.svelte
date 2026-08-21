<script lang="ts">
	import type { RemoteFormIssue } from '@sveltejs/kit';

	let { issues }: { issues: RemoteFormIssue[] | undefined } = $props();
</script>

{#each issues ?? [] as issue}
	<p class="text-destructive text-sm" role="alert">{issue.message}</p>
{/each}
