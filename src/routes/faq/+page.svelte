<script>
	import Card from '$lib/components/ui/card/Card.svelte';
	import CardContent from '$lib/components/ui/card/CardContent.svelte';
	import CardHeader from '$lib/components/ui/card/CardHeader.svelte';
	import { m } from '$lib/paraglide/messages.js';
</script>

<svelte:head>
	<title>MachKinfach - {m.faq_main_title()}</title>
</svelte:head>

<Card>
	<CardHeader title={m.faq_main_title()} />
	<CardContent>
		<p>{m.faq_introduction()}</p>
	</CardContent>
</Card>

<Card>
	<CardHeader title={m.faq_1_question()} />
	<CardContent>
		<p>{m.faq_1_answer()}</p>
	</CardContent>
</Card>

<Card>
	<CardHeader title={m.faq_2_question()} />
	<CardContent>
		<p>{m.faq_2_answer()}</p>
	</CardContent>
</Card>

<Card>
	<CardHeader title={m.faq_3_question()} />
	<CardContent>
		<p>{m.faq_3_answer()}</p>
	</CardContent>
</Card>

<Card>
	<CardHeader title={m.faq_4_question()} />
	<CardContent>
		<p>{m.faq_4_answer()}</p>
	</CardContent>
</Card>

<Card>
	<CardHeader title={m.faq_5_question()} />
	<CardContent>
		<p>{m.faq_5_answer()}</p>
	</CardContent>
</Card>

<Card>
	<CardHeader title={m.faq_6_question()} />
	<CardContent>
		<p>{m.faq_6_answer()}</p>
	</CardContent>
</Card>

<Card>
	<CardHeader title={m.faq_7_question()} />
	<CardContent>
		<p>{m.faq_7_answer()}</p>
	</CardContent>
</Card>

<Card>
	<CardHeader title={m.faq_8_question()} />
	<CardContent>
		<p>{m.faq_8_answer()}</p>
	</CardContent>
</Card>

<Card>
	<CardHeader title={m.faq_9_question()} />
	<CardContent>
		<p>{m.faq_9_answer()}</p>
	</CardContent>
</Card>

<Card>
	<CardHeader title={m.faq_10_question()} />
	<CardContent>
		<p>{m.faq_10_answer()}</p>
	</CardContent>
</Card>

