<script lang="ts">
	import * as m from '$lib/paraglide/messages';
	import { X } from '@lucide/svelte';

	let { packageName, points, onClose }: {
		packageName: string,
		points: number,
		onClose: () => void
	} = $props();

</script>

<div class="fixed inset-0 z-50 flex items-center justify-center">
	<!-- Backdrop -->
	<button class="absolute inset-0 bg-black/60 backdrop-blur-sm" onclick={onClose} aria-label="Close modal"></button>

	<!-- Modal -->
	<div class="relative bg-[#1a1a2e] rounded-2xl p-6 max-w-md w-full mx-4 border border-white/10 shadow-2xl">
		<button
			onclick={onClose}
			class="absolute top-4 right-4 text-white/60 hover:text-white transition-colors">
			<X class="w-5 h-5" />
		</button>

		<div class="text-center">
			<div class="w-16 h-16 bg-green-500/20 rounded-full flex items-center justify-center mx-auto mb-4">
				<svg class="w-8 h-8 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
				</svg>
			</div>

			<h2 class="text-2xl font-bold text-white mb-2">{m.buddypoints_purchase_success_title()}</h2>
			<p class="text-white/70 mb-6">
				{m.buddypoints_purchase_success_message({ name: packageName })}
				{#if points}
					<br /><br />
					{m.buddypoints_purchase_success_points_added({ points })}
				{/if}
			</p>

			<button
				onclick={onClose}
				class="w-full py-3 px-4 bg-primary hover:bg-primary/90 text-white font-semibold rounded-xl transition-colors">
				{m.buddypoints_purchase_success_continue()}
			</button>
		</div>
	</div>
</div>
