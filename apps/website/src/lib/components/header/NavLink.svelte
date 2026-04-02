<script lang="ts">
	import { page } from '$app/state';

	const { href, text } = $props();
</script>

<a
	class="text-text text-lg py-2.5 px-4 rounded-xl {page.url.pathname === href ? 'border bg-purple-300/35 border-violet-300/90 shadow-[0_0_18px_rgba(200,160,255,0.7)]' : 'bg border-main'}"
	href={href}
>{text}</a>

<style>
    .bg {
        background: rgba(14, 14, 22, 0.9);
        transition: all 0.3s ease;
    }
</style>
