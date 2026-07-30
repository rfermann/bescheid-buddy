<script lang="ts">
	type Step = {
		id: string;
		label: string;
	};

	type Props = {
		steps: Step[];
		currentStep: string;
		finished?: boolean;
	};

	let { steps, currentStep, finished = false }: Props = $props();

	let currentIndex = $derived(steps.findIndex((step) => step.id === currentStep));
	let safeIndex = $derived(Math.max(0, currentIndex));
	let progressPercent = $derived(
		steps.length > 1 ? (safeIndex / (steps.length - 1)) * 100 : 0
	);

	let isCompleted = (index: number) => index < currentIndex || (finished && index === currentIndex);
	let isActive = (index: number) => index === currentIndex && !finished;
</script>


{#if currentStep !== ''}
	<div class="w-full">
		<div class="relative flex items-center justify-between">
			<div class="absolute left-0 right-0 top-1/2 h-0.5 -translate-y-1/2 overflow-hidden bg-white/15">
				<div class="animate-indeterminate absolute inset-0 h-full w-full bg-gradient-to-r from-transparent via-white/20 to-transparent"></div>
			</div>
			<div
				class="absolute left-0 top-1/2 h-0.5 -translate-y-1/2 bg-primary transition-[width] duration-300"
				style="width: {progressPercent}%"
			></div>

			{#each steps as step, index (step.id)}
				{@const completed = isCompleted(index)}
				{@const active = isActive(index)}
				<div class="relative z-10 flex flex-col items-center gap-2">
					<div
						class="flex h-5 w-5 items-center justify-center rounded-full border-2 {completed ? 'border-primary bg-primary' : active ? 'border-primary bg-primary-bg' : 'border-white/30 bg-transparent'}"
					>
						{#if completed}
							<svg
								class="h-3 w-3 text-[#0a0a12]"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="4"
								stroke-linecap="round"
								stroke-linejoin="round"
							>
								<polyline points="20 6 9 17 4 12"></polyline>
							</svg>
						{:else if active}
							<div class="h-2 w-2 rounded-full bg-primary"></div>
						{/if}
					</div>
					<span
						class="text-center text-xs font-medium {active ? 'text-primary' : completed ? 'text-white' : 'text-[#b7b7c8]'}"
					>
						{step.label}
					</span>
				</div>
			{/each}
		</div>
	</div>
{/if}

<style>
	.animate-indeterminate {
		animation: indeterminate 0.8s ease-in-out infinite;
	}

	@keyframes indeterminate {
		0% {
			transform: translateX(-100%);
		}
		100% {
			transform: translateX(250%);
		}
	}
</style>
