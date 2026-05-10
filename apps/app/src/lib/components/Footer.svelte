<script>
	import LanguageSwitcher from '$lib/components/LanguageSwitcher.svelte';
	import { m } from '$lib/paraglide/messages.js';
</script>

<footer class="text-center p-5 shrink-0 text-neutral-400 text-sm">
	<div>
		<div
			class="text-teal-100 mb-2.5 font-bold border-primary px-3 bg-primary-bg inline-flex items-center border rounded-full py-2.5">{m.footer_dsgvo()}</div>
		<div class="flex justify-center flex-col gap-3 md:flex-row">
			<span
				class="text-primary ">MachKInfach Media & Digital Services</span>
			<a class="text-primary border-t border-white/80 pt-3 md:border-t-0 md:pt-0 md:border-s md:border-white/80 md:ps-3" href="/agb">{m.footer_agb()}</a>
			<a class="text-primary border-t border-white/80 pt-3 md:border-t-0 md:pt-0 md:border-s md:border-white/80 md:ps-3" href="/datenschutz">{m.footer_datenschutz()}</a>
			<a class="text-primary border-t border-white/80 pt-3 md:border-t-0 md:pt-0 md:border-s md:border-white/80 md:ps-3" href="/impressum">{m.footer_impressum()}</a>
		</div>
	</div>
	<LanguageSwitcher />
</footer>
