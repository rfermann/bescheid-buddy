<script lang="ts">
	type Props = {
		label: string,
		name: string,
		rows: number
	}

	let { label, name, rows }: Props = $props();
</script>
<div class="flex flex-col gap-1">
	<label class="text-xs font-extrabold mb-2 tracking-widest uppercase" for={name}>{label}</label>
	<textarea {rows} id={name} name={name} class="border  border-white/22 bg-white/5 rounded-xl">
	</textarea>
</div>