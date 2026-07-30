<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';

	let { onFileDrop, accept = '*/*' }: { onFileDrop: (files: File[]) => void; accept?: string } = $props();
	let isDragging = $state(false);
	let inputElement: HTMLInputElement;

	const handleDragOver = (e: DragEvent) => {
		e.preventDefault();
		isDragging = true;
	};

	const handleDragLeave = (e: DragEvent) => {
		e.preventDefault();
		isDragging = false;
	};

	const handleDrop = (e: DragEvent) => {
		e.preventDefault();
		isDragging = false;
		const files = Array.from(e.dataTransfer?.files || []);
		if (files.length > 0) {
			onFileDrop(files);
		}
	};

	const handleClick = () => {
		inputElement.click();
	};

	const handleKeyDown = (e: KeyboardEvent) => {
		if (e.key === 'Enter' || e.key === ' ') {
			e.preventDefault();
			inputElement.click();
		}
	};

	const handleInputChange = (e: Event) => {
		const target = e.target as HTMLInputElement;
		const files = Array.from(target.files || []);
		if (files.length > 0) {
			onFileDrop(files);
		}
		target.value = '';
	};
</script>

<div
	class="border-2 border-dashed rounded-xl p-8 text-center cursor-pointer transition-colors {isDragging
		? 'border-primary bg-primary/10'
		: 'border-primary/30 bg-primary/5 hover:border-primary/50 hover:bg-primary/10'}"
	ondragover={handleDragOver}
	ondragleave={handleDragLeave}
	ondrop={handleDrop}
	onclick={handleClick}
	onkeydown={handleKeyDown}
	role="button"
	tabindex="0"
>
	<input
		bind:this={inputElement}
		type="file"
		{accept}
		class="hidden"
		name="file"
		multiple={false}
		onchange={handleInputChange}
	/>
	<div class="flex flex-col items-center gap-3">
		<div class="text-4xl">📄</div>
		<div>
			<p class="text-[#e9e9f1] font-medium">{m.analyse_dropzone_drag()}</p>
			<p class="text-sm text-[#b7b7c8]">{m.analyse_dropzone_click()}</p>
		</div>
	</div>
</div>
