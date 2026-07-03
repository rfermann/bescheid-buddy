<script lang="ts" module>
	import BotIcon from '@lucide/svelte/icons/bot';
	import BriefCaseMedicalIcon from '@lucide/svelte/icons/briefcase-medical';

	const data = {
		user: {
			name: 'shadcn',
			email: 'm@example.com',
			avatar: '/avatars/shadcn.jpg'
		},
		navMain: [
			{
				title: 'Dashboard',
				url: '#',
				icon: BotIcon,
				isActive: true,
				items: []
			},
			{
				title: 'Support Cases',
				url: '/support-cases',
				icon: BriefCaseMedicalIcon,
				items: []
			}
		]
	};
</script>
<script lang="ts">
	import NavMain from './nav-main.svelte';
	import NavUser from './nav-user.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';

	let {
		ref = $bindable(null),
		collapsible = 'icon',
		...restProps
	}: ComponentProps<typeof Sidebar.Root> = $props();
</script>
<Sidebar.Root bind:ref {collapsible} {...restProps}>
	<Sidebar.Content>
		<NavMain items={data.navMain} />
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser user={data.user} />
	</Sidebar.Footer>
	<Sidebar.Rail />
</Sidebar.Root>