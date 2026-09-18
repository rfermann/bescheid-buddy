<script lang="ts">
	import { onMount } from 'svelte';

	type Props = {
		label: string;
		description: string;
		clearLabel: string;
		value?: string;
		disabled?: boolean;
	};

	let { label, description, clearLabel, value = $bindable(''), disabled = false }: Props = $props();
	let canvas: HTMLCanvasElement;
	let context: CanvasRenderingContext2D | null = null;
	let drawing = false;

	const point = (event: PointerEvent) => {
		const bounds = canvas.getBoundingClientRect();
		return {
			x: (event.clientX - bounds.left) * (canvas.width / bounds.width),
			y: (event.clientY - bounds.top) * (canvas.height / bounds.height)
		};
	};

	const start = (event: PointerEvent) => {
		if (disabled || !context) {
			return;
		}
		drawing = true;
		canvas.setPointerCapture(event.pointerId);
		const position = point(event);
		context.beginPath();
		context.moveTo(position.x, position.y);
	};

	const draw = (event: PointerEvent) => {
		if (!drawing || disabled || !context) {
			return;
		}
		const position = point(event);
		context.lineTo(position.x, position.y);
		context.stroke();
	};

	const finish = (event: PointerEvent) => {
		if (!drawing) {
			return;
		}
		drawing = false;
		if (canvas.hasPointerCapture(event.pointerId)) {
			canvas.releasePointerCapture(event.pointerId);
		}
		value = canvas.toDataURL('image/png');
	};

	const clear = () => {
		if (!context || disabled) {
			return;
		}
		context.clearRect(0, 0, canvas.width, canvas.height);
		value = '';
	};

	onMount(() => {
		context = canvas.getContext('2d');
		if (context) {
			context.lineCap = 'round';
			context.lineJoin = 'round';
			context.lineWidth = 4;
			context.strokeStyle = '#111827';
		}
	});
</script>

<div class="flex flex-col gap-2">
	<label for="papierbrief-signature" class="text-xs font-extrabold tracking-widest uppercase">
		{label}
	</label>
	<p id="papierbrief-signature-description" class="text-sm text-[#b7b7c8]">{description}</p>
	<canvas
		bind:this={canvas}
		id="papierbrief-signature"
		width="900"
		height="240"
		aria-label={label}
		aria-describedby="papierbrief-signature-description"
		onpointerdown={start}
		onpointermove={draw}
		onpointerup={finish}
		onpointercancel={finish}
		class="h-40 w-full touch-none rounded-xl border border-white/22 bg-white disabled:opacity-50"
		class:opacity-50={disabled}
	></canvas>
	<button
		type="button"
		onclick={clear}
		disabled={disabled || !value}
		class="self-start rounded-lg border border-white/20 px-3 py-2 text-sm font-semibold text-white hover:border-white/40 disabled:cursor-not-allowed disabled:opacity-50"
	>
		{clearLabel}
	</button>
</div>
