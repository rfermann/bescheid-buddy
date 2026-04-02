<script>
	import LanguageSwitcher from '$lib/components/LanguageSwitcher.svelte';
	import { m } from '$lib/paraglide/messages.js';
</script>

<footer class="text-center p-5 shrink-0 text-neutral-400 text-sm">
	<div>

		<div
			class="text-teal-100 mb-2.5 font-bold border-accent px-3 bg-accent-bg inline-flex items-center border rounded-full py-2.5">{m.footer_dsgvo()}</div>
		<div>
			<span>MachKInfach Media & Digital Services</span>
			<span> | </span>
			<a class="text-secondary" href="/impressum">{m.footer_impressum()}</a>
			<span> | </span>
			<a class="text-secondary" href="/datenschutz">{m.footer_datenschutz()}</a>
		</div>
	</div>
	<LanguageSwitcher />
</footer>
