<script>
	import Card from '$lib/components/ui/card/Card.svelte';
	import CardContent from '$lib/components/ui/card/CardContent.svelte';
	import CardHeader from '$lib/components/ui/card/CardHeader.svelte';
	import { m } from '$lib/paraglide/messages.js';
</script>

<svelte:head>
	<title>MachKinfach - {m.about_us_main_title()}</title>
</svelte:head>

<Card>
	<CardHeader title={m.about_us_main_title()} />
	<CardContent>
		<p>{m.about_us_main_content_1()}</p>
	</CardContent>
</Card>
