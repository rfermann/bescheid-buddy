<script>
	import NavLink from '$lib/components/header/NavLink.svelte';
	import { m } from '$lib/paraglide/messages.js';
</script>

<header class="text-center pt-8 pb-2.5 px-4 flex flex-col items-center relative">
	<img
		class="w-36 h-auto mb-4 animate-[pulseGlow_3s_ease-in-out_infinite] filter-[drop-shadow(0_0_18px_rgba(192,140,255,0.85))]"
		src="/img/logo.png" alt="">
	<h1 class="text-primary leading-14 font-bold text-3xl [text-shadow:0_0_10px_rgba(192,140,255,0.9)]">MachKinfach</h1>
	<nav class="mt-5 flex flex-wrap gap-3 px-2.5 justify-center">
		<NavLink href="/" text={m.nav_bescheidbuddy()} />
		<NavLink href="/ueber-uns" text={m.nav_ueber_uns()} />
		<NavLink href="/partner" text={m.nav_fuer_partner()} />
		<NavLink href="/kontakt" text={m.nav_kontakt()} />
		<NavLink href="/faq" text={m.nav_fragen_antworten()} />
	</nav>
</header>
