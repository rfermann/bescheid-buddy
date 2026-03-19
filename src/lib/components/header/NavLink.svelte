<script lang="ts">
	import { page } from '$app/state';

	const { href, text } = $props();
</script>

<a
	class="text-text text-lg py-2.5 px-4 rounded-xl bg border-main {page.url.pathname === href ? 'bg-active' : ''}"
	href={href}
>{text}</a>

<style>
    .bg {
        background: rgba(14, 14, 22, 0.9);
        transition: all 0.3s ease;
    }
</style>
