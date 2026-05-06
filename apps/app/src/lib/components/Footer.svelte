<script>
	import LanguageSwitcher from '$lib/components/LanguageSwitcher.svelte';
	import { m } from '$lib/paraglide/messages.js';
</script>

<footer class="text-center p-5 shrink-0 text-neutral-400 text-sm">
	<div>
		<div
			class="text-teal-100 mb-2.5 font-bold border-primary px-3 bg-primary-bg inline-flex items-center border rounded-full py-2.5">{m.footer_dsgvo()}</div>
		<div class="flex justify-center flex-row gap-3">
			<span class="text-primary">MachKInfach Media & Digital Services</span>
			<span> | </span>
			<a class="text-primary" href="/agb">{m.footer_agb()}</a>
			<span> | </span>
			<a class="text-primary" href="/datenschutz">{m.footer_datenschutz()}</a>
			<span> | </span>
			<a class="text-primary" href="/impressum">{m.footer_impressum()}</a>
		</div>
	</div>
	<LanguageSwitcher />
</footer>
