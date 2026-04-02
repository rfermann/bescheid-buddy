<script>
	import Card from '$lib/components/ui/card/Card.svelte';
	import CardContent from '$lib/components/ui/card/CardContent.svelte';
	import CardHeader from '$lib/components/ui/card/CardHeader.svelte';
	import { m } from '$lib/paraglide/messages.js';
	import ContactForm from '$lib/components/ContactForm.svelte';
</script>

<svelte:head>
	<title>MachKinfach - {m.partner_title()}</title>
</svelte:head>

<Card>
	<CardHeader title={m.partner_title()} />
	<CardContent>
		<p>{m.partner_description()}</p>
	</CardContent>
</Card>
<Card>
	<CardHeader title={m.kontakt_title()} />
	<CardContent>
		<ContactForm type="partner" key="partner" />
	</CardContent>
</Card>
