<script lang="ts">
	import { onMount } from 'svelte';

	let canvas: HTMLCanvasElement | null = $state(null);

	let pointer = $state({ x: 0, y: 0 });
	let size = $state({ w: 0, h: 0, dpr: 1 });

	const updatePointerFromClientXY = (clientX: number, clientY: number) => {
		pointer.x = clientX;
		pointer.y = clientY;
	};

	const resize = () => {
		if (!canvas) {
			return;
		}

		const dpr = Math.max(1, window.devicePixelRatio || 1);
		const w = window.innerWidth;
		const h = window.innerHeight;

		size.w = w;
		size.h = h;
		size.dpr = dpr;

		canvas.width = Math.floor(w * dpr);
		canvas.height = Math.floor(h * dpr);

		const ctx = canvas.getContext('2d');
		if (!ctx) return;
		ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

		if (pointer.x === 0 && pointer.y === 0) {
			pointer.x = w / 2;
			pointer.y = h / 2;
		}
	};

	const onMouseMove = (e: MouseEvent) => {
		updatePointerFromClientXY(e.clientX, e.clientY);
	};

	const onTouchMove = (e: TouchEvent) => {
		const t = e.touches[0];
		if (!t) return;
		updatePointerFromClientXY(t.clientX, t.clientY);
	};

	onMount(() => {
		resize();
	});

	let raf = 0;

	const draw = () => {
		if (!canvas) return;
		const ctx = canvas.getContext('2d');
		if (!ctx) return;

		ctx.clearRect(0, 0, size.w, size.h);

		const r = Math.min(size.w, size.h) * 0.25;
		const g = ctx.createRadialGradient(pointer.x, pointer.y, 0, pointer.x, pointer.y, r);
		g.addColorStop(0, 'rgba(230,220,255,0.35)');
		g.addColorStop(1, 'rgba(0,0,0,0)');

		ctx.fillStyle = g;
		ctx.beginPath();
		ctx.arc(pointer.x, pointer.y, r, 0, Math.PI * 2);
		ctx.fill();

		raf = window.requestAnimationFrame(draw);
	};

	$effect(() => {
		if (canvas) {
			raf = window.requestAnimationFrame(draw);
			return () => {
				window.cancelAnimationFrame(raf);
			};
		}
	});
</script>

<svelte:window onresize={resize} onmousemove={onMouseMove} ontouchmove={onTouchMove} />

<canvas bind:this={canvas} class="neon-canvas" aria-hidden="true"></canvas>

<style>
    .neon-canvas {
        position: fixed;
        inset: 0;
        width: 100%;
        height: 100%;
        z-index: -1;
        pointer-events: none;
        background: radial-gradient(circle at top, #2a123e 0, #050509 55%);
    }
</style>