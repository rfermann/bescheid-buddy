<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Separator } from '$lib/components/ui/separator';
	import { Avatar, AvatarImage, AvatarFallback } from '$lib/components/ui/avatar';
	import * as Select from '$lib/components/ui/select';
	import * as Dialog from '$lib/components/ui/dialog';
	import {
		getSupportCaseDetail,
		getSupportCaseMessages,
		sendMessageForm,
		updateStatusCommand
	} from './support-case-detail.remote';
	import ArrowLeftIcon from '@lucide/svelte/icons/arrow-left';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import SendIcon from '@lucide/svelte/icons/send';
	import Loader2Icon from '@lucide/svelte/icons/loader-2';

	const caseId = $derived(page.params.id as string);
	const supportCase = $derived(await getSupportCaseDetail({ id: caseId }));
	const messages = $derived(await getSupportCaseMessages({ caseId }));

	let selectedStatus = $derived(supportCase.status);
	let showOpenDialog = $state(false);

	const statusColors: Record<string, string> = {
		open: 'bg-red-500',
		in_progress: 'bg-yellow-500',
		resolved: 'bg-green-500',
		closed: 'bg-gray-500'
	};

	const statuses = ['open', 'in_progress', 'resolved', 'closed'];

	const formatDate = (date: Date) => {
		return new Date(date).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	};

	const formatTime = (date: Date) => {
		return new Date(date).toLocaleTimeString('en-US', {
			hour: '2-digit',
			minute: '2-digit'
		});
	};

	const isValidStatus = (status: string): status is 'in_progress' | 'resolved' | 'closed' => {
		return status === 'in_progress' || status === 'resolved' || status === 'closed';
	};

	const handleStatusChange = async (status: string) => {
		if (status === 'open') {
			showOpenDialog = true;
			return;
		}
		if (!isValidStatus(status)) {
			return;
		}
		selectedStatus = status;
		await updateStatusCommand({ id: caseId, status });
	};
</script>

<div class="space-y-4">
	<div class="flex items-center gap-2">
		<Button variant="ghost" size="icon" href="/support-cases">
			<ArrowLeftIcon class="h-5 w-5" />
		</Button>
		<h1 class="text-2xl font-bold">Support Case Details</h1>
	</div>

	<Card>
		<div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between px-6 pt-6">
			<div class="space-y-2">
				<CardTitle class="text-xl">{supportCase.subject}</CardTitle>
				<div class="flex flex-wrap items-center gap-2">
					<Badge class={statusColors[supportCase.status] || 'bg-gray-500'}>
						{supportCase.status.replace('_', ' ')}
					</Badge>
					<Badge variant="outline">{supportCase.category.replace('_', ' ')}</Badge>
				</div>
			</div>
			<div class="flex items-center gap-2">
				<Button
					onclick={() => {
						getSupportCaseDetail({ id: caseId }).refresh();
						getSupportCaseMessages({ caseId }).refresh();
					}}
					variant="outline"
					size="sm">
					<RefreshCwIcon class="mr-2 h-4 w-4" />
					Refresh
				</Button>
			</div>
		</div>
		<CardContent class="space-y-4">
			<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
				<div>
					<p class="text-sm font-medium text-muted-foreground">User</p>
					<p class="text-sm font-semibold">{supportCase.userName}</p>
					{#if supportCase.userEmail}
						<p class="text-xs text-muted-foreground">{supportCase.userEmail}</p>
					{/if}
				</div>
				<div>
					<p class="text-sm font-medium text-muted-foreground">Created</p>
					<p class="text-sm">{formatDate(supportCase.created_at)}</p>
				</div>
				<div>
					<p class="text-sm font-medium text-muted-foreground">Last Updated</p>
					<p class="text-sm">{formatDate(supportCase.updated_at)}</p>
				</div>
				<div>
					<p class="text-sm font-medium text-muted-foreground">Status</p>
					<Select.Root type="single" value={selectedStatus} onValueChange={handleStatusChange}>
						<Select.Trigger class="w-full">
							{selectedStatus.replace('_', ' ')}
						</Select.Trigger>
						<Select.Content>
							{#each statuses as status (status)}
								<Select.Item value={status}>{status.replace('_', ' ')}</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
			</div>

			<Separator />

			<div class="space-y-2">
				<h3 class="text-lg font-semibold">Messages</h3>
				<div class="space-y-4 rounded-md border p-4">
					{#each messages as message (message.id)}
						<div
							class="flex gap-3 {message.is_user_message ? 'justify-start' : 'justify-end'}">
							{#if message.is_user_message}
								<Avatar class="h-8 w-8">
									<AvatarImage src="" alt={supportCase.userName} />
									<AvatarFallback>{supportCase.userName.charAt(0).toUpperCase()}</AvatarFallback>
								</Avatar>
							{/if}
							<div
								class="max-w-[70%] rounded-lg p-3 {message.is_user_message
									? 'bg-muted'
									: 'bg-primary text-primary-foreground'}">
								<div class="mb-1 flex items-center gap-2">
									<span class="text-xs font-semibold">
										{message.is_user_message ? supportCase.userName : 'Admin'}
									</span>
									<span class="text-xs opacity-70">{formatTime(message.created_at)}</span>
								</div>
								<p class="text-sm whitespace-pre-wrap">{message.message}</p>
							</div>
							{#if !message.is_user_message}
								<Avatar class="h-8 w-8">
									<AvatarImage src="" alt="Admin" />
									<AvatarFallback>A</AvatarFallback>
								</Avatar>
							{/if}
						</div>
					{:else}
						<p class="text-center text-sm text-muted-foreground">No messages yet.</p>
					{/each}
				</div>
			</div>

			<Separator />

			<div class="space-y-2">
				<h3 class="text-lg font-semibold">Send Message</h3>
				<form {...sendMessageForm} class="flex gap-2">
					<input {...sendMessageForm.fields.caseId.as('hidden', caseId)} />
					<textarea
						{...sendMessageForm.fields.message.as('text')}
						disabled={!!sendMessageForm.pending}
						class="flex min-h-20 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
					></textarea>
					<Button type="submit" disabled={!!sendMessageForm.pending} class="self-end">
						{#if sendMessageForm.pending}
							<Loader2Icon class="mr-2 h-4 w-4 animate-spin" />
						{:else}
							<SendIcon class="mr-2 h-4 w-4" />
						{/if}
						Send
					</Button>
				</form>
			</div>
		</CardContent>
	</Card>
</div>

<Dialog.Root open={showOpenDialog} onOpenChange={(open) => showOpenDialog = open}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Status Change Not Allowed</Dialog.Title>
			<Dialog.Description>
				Changing the status back to "open" is not allowed.
			</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Button onclick={() => showOpenDialog = false}>OK</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
