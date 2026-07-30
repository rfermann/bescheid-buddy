<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';
	import PageHeading from '$lib/components/PageHeading.svelte';
	import Section from '$lib/components/Section.svelte';
	import StepProgress from '$lib/components/StepProgress.svelte';
	import FileSelector from './FileSelector.svelte';
	import ModelSelector from './ModelSelector.svelte';
	import ConsentGate from './ConsentGate.svelte';
	import Result from './Result.svelte';
	import {
		runAnalysisPipeline,
		type AnalysisModel,
		type AnalysisPipelineResult,
		type AnalysisStep
	} from './analyse.client';

	let selectedFile = $state<File | null>(null);
	let selectedModel = $state<AnalysisModel>('gpt-4o-mini');
	let isUploading = $state(false);
	let currentStep = $state('');
	let uploadError = $state<string | null>('some rubbish');
	let analysisResult = $state<AnalysisPipelineResult | null>(null);

	const steps = [
		{ id: 'preparing', label: m.analyse_step_preparing() },
		{ id: 'uploading', label: m.analyse_step_uploading() },
		{ id: 'extracting', label: m.analyse_step_extracting() },
		{ id: 'analyzing', label: m.analyse_step_analyzing() },
		{ id: 'finished', label: m.analyse_step_finished() }
	];

	const handleFileSelected = (file: File | null) => {
		selectedFile = file;
		if (file) uploadError = null;
	};

	const runAnalysis = async () => {
		const file = selectedFile;
		if (!file) {
			return;
		}
		isUploading = true;
		uploadError = null;
		analysisResult = null;
		try {
			analysisResult = await runAnalysisPipeline({
				file,
				model: selectedModel,
				onStep: (step: AnalysisStep) => (currentStep = step)
			});
		} catch {
			uploadError = m.analyse_error_upload_failed();
			currentStep = '';
		} finally {
			isUploading = false;
		}
	};
</script>

<PageHeading title={m.analyse_page_title()} />

<Section>
	<div class="mb-6 flex flex-col items-start justify-between gap-6 md:flex-row md:items-center">
		<div>
			<h2 class="mb-2 text-2xl font-semibold text-white">{m.analyse_upload_title()}</h2>
			<p class="text-sm text-[#b7b7c8]">{m.analyse_upload_subtitle()}</p>
		</div>
	</div>
	<FileSelector onFileSelected={handleFileSelected} />
	<ModelSelector bind:value={selectedModel} disabled={isUploading} />
	<div class="mt-6 grid grid-cols-1 gap-3">
		<ConsentGate canProceed={!!selectedFile} onConfirm={runAnalysis} />
		{#if uploadError}
			<div
				role="alert"
				class="rounded-lg border border-red-400/30 bg-red-400/10 px-3 py-2 text-sm text-red-300"
			>
				{uploadError}
			</div>
		{/if}
		<StepProgress {steps} {currentStep} finished={currentStep === 'finished'} />
	</div>
	<Result result={analysisResult} />
</Section>
