<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';

	let {
		open = $bindable(false),
		onConfirm,
		onCancel
	}: {
		open?: boolean;
		onConfirm: (dontShowAgain: boolean) => void;
		onCancel: () => void;
	} = $props();

	let dialog = $state<HTMLDialogElement>();
	let dontShowAgain = $state(false);

	$effect(() => {
		if (!dialog) return;
		if (open && !dialog.open) dialog.showModal();
		else if (!open && dialog.open) dialog.close();
	});
</script>

<dialog
	bind:this={dialog}
	oncancel={(e) => {

		e.preventDefault();
		onCancel();
	}}
	onclick={(e) => {
		if (e.target === dialog) onCancel();
	}}
	class="m-auto w-full max-w-md rounded-xl border border-white/10 bg-[#1a1a24] p-0 text-white backdrop:bg-black/60"
>
	<div class="flex flex-col gap-4 p-6">
		<h2 class="text-xl font-semibold text-white">{m.analyse_consent_title()}</h2>
		<p class="text-sm text-[#b7b7c8]">{m.analyse_consent_body()}</p>

		<div class="flex items-center gap-3 text-sm">
			<a href="/datenschutz" target="_blank" rel="noopener" class="text-primary hover:underline">
				{m.analyse_consent_privacy_link()}
			</a>
			<span class="text-white/25">·</span>
			<a href="/agb" target="_blank" rel="noopener" class="text-primary hover:underline">
				{m.analyse_consent_terms_link()}
			</a>
		</div>

		<label class="flex items-center gap-2 text-sm text-[#b7b7c8]">
			<input
				type="checkbox"
				bind:checked={dontShowAgain}
				class="h-4 w-4 rounded border-white/20 bg-white/5 accent-primary"
			/>
			{m.analyse_consent_dont_show_again()}
		</label>

		<div class="mt-2 flex justify-end gap-3">
			<button
				type="button"
				onclick={onCancel}
				class="rounded-2xl border border-white/10 px-3.5 py-2.5 text-sm font-medium text-[#b7b7c8] transition-colors hover:border-white/25 hover:text-white"
			>
				{m.analyse_consent_cancel()}
			</button>
			<PrimaryButton
				type="button"
				onclick={() => onConfirm(dontShowAgain)}
				text={m.analyse_consent_continue()}
			/>
		</div>
	</div>
</dialog>
