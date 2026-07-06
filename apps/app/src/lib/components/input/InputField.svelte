<script lang="ts">
	import type { HTMLInputTypeAttribute } from 'svelte/elements';

	type Props = {
		label: string,
		name: string,
		type: HTMLInputTypeAttribute,
	}

	let { label, name, type }: Props = $props();
</script>
<div class="flex flex-col gap-1">
	<label class="text-xs font-extrabold mb-2 tracking-widest uppercase" for={name}>{label}</label>
	<input type={type} id={name} name={name} class="border  border-white/22 bg-white/5 rounded-xl" />
</div>