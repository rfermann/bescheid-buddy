<script lang="ts">
	import NavDropdown from './NavDropdown.svelte';
	import { m } from '$lib/paraglide/messages.js';
	import { page } from '$app/state';

	interface NavItem {
		id: string;
		label: () => string;
		href: string;
		disabled?: boolean;
		hidden?: boolean;
		badge?: () => string;
	}

	let isOpen = $state(false);

	$effect(() => {
		page.url.pathname;
		isOpen = false;
	});

	const navItems: NavItem[] = [
		{ id: 'analysis', label: m.nav_analysis, href: '/' },
		{ id: 'paper-letter', label: m.nav_send_paper_letter, href: '/paper-letter' },
		{ id: 'archive', label: m.nav_archive, href: '/archive' },
		{ id: 'reply-pdf', label: m.nav_reply_pdf, href: '/reply-pdf', hidden: true },
		{ id: 'solutions', label: m.nav_solutions, href: '/solutions' },
		{ id: 'partners', label: m.nav_partners, href: '/partners', disabled: true, badge: m.nav_soon },
		{ id: 'expert-work', label: m.nav_expert_work, href: '/expert-work' },
		{ id: 'datenschutz', label: m.nav_datenschutz, href: '/datenschutz' },
		{ id: 'support', label: m.nav_support, href: '/support' },
		{ id: 'data', label: m.nav_data_delete, href: '/data' }
	];

	const toggleDropdown = () => {
		isOpen = !isOpen;
	};


	const getCurrentLabel = () => {
		const item = navItems.find((i) => i.href === page.url.pathname);
		return item ? item.label() : m.nav_menu();
	};
</script>

<aside class="fixed left-4 top-20 z-50 w-[min(320px,calc(100vw-36px))]">
	<NavDropdown
		{isOpen}
		currentLabel={getCurrentLabel()}
		{navItems}
		menuLabel={m.nav_menu}
		onToggle={toggleDropdown}
	/>
</aside>
