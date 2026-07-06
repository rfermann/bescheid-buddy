<script lang="ts">
	import type { PageProps } from './$types';
	import PageHeading from '$lib/components/PageHeading.svelte';
	import Section from '$lib/components/Section.svelte';
	import Textarea from '$lib/components/input/Textarea.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import { getSupportCaseByIdQuery, getSupportCaseMessagesQuery, sendMessageForm, resolveSupportCaseCommand, closeSupportCaseCommand } from '../support.remote.ts';

	let { params }: PageProps = $props();
	const supportCase = $derived(await getSupportCaseByIdQuery({ id: params.id }));
	const supportCaseMessages = $derived(await getSupportCaseMessagesQuery({ caseId: params.id }));


	const formatDate = (dateString: Date): string => {
		return dateString.toLocaleDateString() + ' ' + dateString.toLocaleTimeString();
	};

	const getCategoryDisplay = (category: string): string => {
		const categoryMap = {
			general: 'General',
			technical: 'Technical',
			billing: 'Billing',
			feature_request: 'Feature Request'
		};
		return categoryMap[category as keyof typeof categoryMap] || category;
	};

	const getStatusDisplay = (status: string): string => {
		const statusMap = {
			open: 'Open',
			in_progress: 'In Progress',
			resolved: 'Resolved',
			closed: 'Closed'
		};
		return statusMap[status as keyof typeof statusMap] || status;
	};

	const getStatusClass = (status: string): string => {
		const statusClasses = {
			open: 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30',
			in_progress: 'bg-blue-500/20 text-blue-400 border-blue-500/30',
			resolved: 'bg-green-500/20 text-green-400 border-green-500/30',
			closed: 'bg-gray-500/20 text-gray-400 border-gray-500/30'
		};
		return statusClasses[status as keyof typeof statusClasses] || '';
	};
</script>

<div class="flex items-center gap-4">
	<a
		href="/support"
		class="inline-flex mb-8 items-center px-3 py-1 rounded-2xl border font-bold border-primary/45 bg-primary/15 text-primary text-xs hover:border-primary/60 transition-colors"
	>
		←
	</a>
	<PageHeading title={`Support Case: ${supportCase.supportCase.subject}`} />
</div>

<Section>
	<div class="space-y-6">
		<div class="flex flex-wrap gap-4 text-sm items-center">
			<div class="flex items-center">
				<span class="text-[#b7b7c8]">Category:</span>
				<span class="ml-2">{getCategoryDisplay(supportCase.supportCase.category)}</span>
			</div>
			<div class="flex items-center">
				<span class="text-[#b7b7c8]">Status:</span>
				<span class="ml-2">
					<span
						class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium border {getStatusClass(supportCase.supportCase.status)}">
						{getStatusDisplay(supportCase.supportCase.status)}
					</span>
				</span>
			</div>
			<div class="flex items-center">
				<span class="text-[#b7b7c8]">Created:</span>
				<span class="ml-2">{formatDate(supportCase.supportCase.created_at)}</span>
			</div>
			<div class="flex items-center">
				<span class="text-[#b7b7c8]">Resolved:</span>
				{#if supportCase.supportCase.resolved_at}
					<span class="ml-2">{formatDate(supportCase.supportCase.resolved_at)}</span>
				{/if}
			</div>
		</div>

		{#if supportCase.supportCase.status !== 'resolved' && supportCase.supportCase.status !== 'closed'}
			<div class="flex gap-3">
				<button
					onclick={async () => await resolveSupportCaseCommand(params.id)}
					class="px-4 py-2 rounded-lg bg-green-500/20 border border-green-500/30 text-green-400 text-sm font-medium hover:bg-green-500/30 transition-colors"
				>
					Resolve
				</button>
				<button
					onclick={async () => await closeSupportCaseCommand(params.id)}
					class="px-4 py-2 rounded-lg bg-gray-500/20 border border-gray-500/30 text-gray-400 text-sm font-medium hover:bg-gray-500/30 transition-colors"
				>
					Close
				</button>
			</div>
		{/if}
	</div>
</Section>

<Section>
	<div class="space-y-6">
		<h2 class="text-2xl font-semibold">Messages</h2>

		<div class="space-y-4">
			{#each supportCaseMessages.messages as message (message.id)}
				<div class="flex {message.is_user_message ? 'justify-end' : 'justify-start'}">
					<div class="max-w-lg">
						<div class="text-xs text-[#b7b7c8] mb-1">
							{message.is_user_message ? 'You' : 'Support'} • {formatDate(message.created_at)}
						</div>
						<div class="px-4 py-3 rounded-2xl {message.is_user_message 
							? 'bg-primary/20 border border-primary/30 text-primary' 
							: 'bg-white/10 border border-white/20 text-[#e9e9f1]'}">
							<p class="text-sm">{message.message}</p>
						</div>
					</div>
				</div>
			{/each}
		</div>

		<div class="mt-8">
			<form
				{...sendMessageForm}
				class="space-y-4"
			>
				<input {...sendMessageForm.fields.caseId.as('hidden', params.id)} />
				<Textarea
					{...sendMessageForm.fields.message.as('text')}
					rows={4}
					label="Your message"
				/>
				<div class="flex justify-end">
					<PrimaryButton type="submit" text="Send Message" />
				</div>
			</form>
		</div>
	</div>
</Section>