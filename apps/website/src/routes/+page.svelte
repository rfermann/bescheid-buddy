<script>
	import Card from '$lib/components/ui/card/Card.svelte';
	import CardContent from '$lib/components/ui/card/CardContent.svelte';
	import CardHeader from '$lib/components/ui/card/CardHeader.svelte';
	import { m } from '$lib/paraglide/messages.js';
	import { env } from '$env/dynamic/public';
</script>

<svelte:head>
	<title>MachKinfach</title>
</svelte:head>
<div
	class="mb-5 rounded-xl border border-accent p-5 text-gray-200 shadow-[0_12px_30px_rgba(0,0,0,0.7),0_0_28px_rgba(0,255,238,0.2)] transition-shadow duration-300"
>
	<CardHeader title={m.landing_main_title()} />
	<CardContent>
		<p>{m.landing_main_content_1()}</p>
		<p>{m.landing_main_content_2()}</p>
		<p>{m.landing_main_content_3()}</p>
		<div>
			<p>{m.landing_main_content_4()}</p>
			<p>{m.landing_main_content_5()}</p>
			<p>{m.landing_main_content_6()}</p>
			<p>{m.landing_main_content_7()}</p>
		</div>
		<p>{m.landing_main_content_8()}</p>
		<a
			class="self-start rounded-2xl border border-purple-400 bg-violet-500/20 px-4 py-3.5 text-xl font-extrabold whitespace-nowrap text-purple-200"
			href={env.PUBLIC_APP_URL}>{m.landing_call_to_action()}</a
		>
	</CardContent>
</div>
<Card>
	<CardHeader title={m.landing_how_it_works_title()} />
	<CardContent>
		<div>
			<p>{m.landing_how_it_works_content_1()}</p>
			<p>{m.landing_how_it_works_content_2()}</p>
			<p>{m.landing_how_it_works_content_3()}</p>
		</div>
		<p>{m.landing_how_it_works_content_4()}</p>
	</CardContent>
</Card>
<Card>
	<CardHeader title={m.landing_security_title()} />
	<CardContent>
		<p>{m.landing_security_content_1()}</p>
		<p>{m.landing_security_content_2()}</p>
		<p>{m.landing_security_content_3()}</p>
	</CardContent>
</Card>
<Card>
	<CardHeader title={m.landing_important_note_title()} />
	<CardContent>
		<p>{m.landing_important_note_content_1()}</p>
		<div>
			<p>{m.landing_important_note_content_2()}</p>
			<ul class="list-inside list-disc">
				<li>
					{m.landing_important_note_content_3()}
				</li>
				<li>
					{m.landing_important_note_content_4()}
				</li>
			</ul>
		</div>
	</CardContent>
</Card>
