<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';
	import FileDropZone from '$lib/components/FileDropZone.svelte';

	const MAX_UPLOAD_MB = 15;
	const MAX_UPLOAD_BYTES = MAX_UPLOAD_MB * 1024 * 1024;

	let { onFileSelected }: { onFileSelected: (file: File | null) => void } = $props();
	let selectedFile = $state<File | null>(null);
	let error = $state<string | null>(null);

	const handleFileDrop = (files: File[]) => {
		const file = files[0];
		if (!file) return;

		if (file.size > MAX_UPLOAD_BYTES) {
			selectedFile = null;
			error = m.analyse_error_file_too_large({ maxMb: MAX_UPLOAD_MB });
			onFileSelected(null);
			return;
		}

		error = null;
		selectedFile = file;
		onFileSelected(file);
	};
</script>

<FileDropZone onFileDrop={handleFileDrop} accept="image/*,.pdf" />

{#if error}
	<span class="text-red-400">{error}</span>
{/if}

{#if selectedFile}
	<span class="text-[#b7b7c8]">{m.analyse_selected_file()} {selectedFile.name}</span>
{/if}
