<script lang="ts">
	import { Html, Head, Body, Container, Section, Text } from '@better-svelte-email/components';
	import { m } from '$lib/paraglide/messages.js';

	let { code, mins }: { code: string; mins: number } = $props();
</script>

<Html>
<Head />
<Body class="bg-zinc-100">
<Container class="m-8 mx-auto max-w-lg rounded-2xl bg-white p-8">
	<Section class="mx-auto text-center">
		<Text class="mt-4 text-zinc-600 whitespace-pre-line">
			{m.email_loginCodeText({ code, mins })}
		</Text>
	</Section>
</Container>
</Body>
</Html>
