<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import * as m from '$lib/paraglide/messages.js';
	import { getLocale } from '$lib/paraglide/runtime';
	import PageHeading from '$lib/components/PageHeading.svelte';
	import Section from '$lib/components/Section.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import InputField from '$lib/components/input/InputField.svelte';
	import Textarea from '$lib/components/input/Textarea.svelte';
	import SignatureCanvas from './SignatureCanvas.svelte';
	import ModelSelector from '../ModelSelector.svelte';
	import { getModelsQuery } from '../analyse.remote';
	import { getAnalysisQuery } from '../meine-analysen/analyses.remote';
	import { getBuddyPointsQuery } from '../meine-buddypoints/points.remote';
	import {
		confirmPaperLetterUploadCommand,
		createPaperLetterDeliveryCommand,
		createGenerationCommand,
		createPDFCommand,
		createPaperLetterUploadCommand,
		type CreateGenerationRequest,
		type CreatePDFRequest
	} from './papierbrief.remote';
	import {
		GenerationCancelledError,
		GenerationTimedOutError,
		DeliveryCancelledError,
		DeliveryTimedOutError,
		pollDelivery,
		pollGeneration
	} from './papierbrief.client';
	import {
		maxPaperLetterPDFBytes,
		paperLetterPDFWithinLimit,
		putPaperLetterPDF
	} from './papierbrief.upload.client';
	import type { Address, DeliveryBody, GenerationResultBody, Metadata } from '../../generated';

	const categories = [
		'Jobcenter',
		'Krankenkasse',
		'Finanzamt',
		'Ausländerbehörde',
		'Rentenversicherung',
		'Arbeitsagentur',
		'Inkasso',
		'Gericht',
		'Schule/Kita',
		'Polizei',
		'Versicherung',
		'Bank',
		'Telekom/Internet',
		'Energie',
		'Other',
		'Universal'
	] as const;
	const subcategories = [
		'widerspruch',
		'fristverlaengerung',
		'unterlagen_nachreichen',
		'zahlung_klaeren',
		'termin_verschieben',
		'kuendigung',
		'freie_formulierung'
	] as const;
	type Category = (typeof categories)[number];
	type Subcategory = (typeof subcategories)[number];

	const categoryLabels: Record<Category, () => string> = {
		Jobcenter: m.papierbrief_category_jobcenter,
		Krankenkasse: m.papierbrief_category_health_insurance,
		Finanzamt: m.papierbrief_category_tax_office,
		Ausländerbehörde: m.papierbrief_category_immigration_office,
		Rentenversicherung: m.papierbrief_category_pension_insurance,
		Arbeitsagentur: m.papierbrief_category_employment_agency,
		Inkasso: m.papierbrief_category_debt_collection,
		Gericht: m.papierbrief_category_court,
		'Schule/Kita': m.papierbrief_category_school,
		Polizei: m.papierbrief_category_police,
		Versicherung: m.papierbrief_category_insurance,
		Bank: m.papierbrief_category_bank,
		'Telekom/Internet': m.papierbrief_category_telecom,
		Energie: m.papierbrief_category_energy,
		Other: m.papierbrief_category_other,
		Universal: m.papierbrief_category_universal
	};
	const subcategoryLabels: Record<Subcategory, () => string> = {
		widerspruch: m.papierbrief_subcategory_objection,
		fristverlaengerung: m.papierbrief_subcategory_extension,
		unterlagen_nachreichen: m.papierbrief_subcategory_documents,
		zahlung_klaeren: m.papierbrief_subcategory_payment,
		termin_verschieben: m.papierbrief_subcategory_reschedule,
		kuendigung: m.papierbrief_subcategory_termination,
		freie_formulierung: m.papierbrief_subcategory_free
	};

	const models = $derived(await getModelsQuery({}));
	const pointsQuery = getBuddyPointsQuery({});
	const buddyPoints = $derived(await pointsQuery);
	const locale = getLocale() as CreateGenerationRequest['locale'];
	const requestedCaseId = page.url.searchParams.get('analyse');
	const requestedDocumentId = page.url.searchParams.get('dokument');

	let selectedModel = $state('');
	let senderName = $state('');
	let senderAddress = $state('');
	let recipientName = $state('');
	let recipientAddress = $state('');
	let category = $state<Category>('Other');
	let subcategory = $state<Subcategory>('freie_formulierung');
	let subject = $state('');
	let reference = $state('');
	let instruction = $state('');
	let germanLetter = $state('');
	let improvementFeedback = $state('');
	let sourceCaseId = $state<string | undefined>();
	let sourceDocumentId = $state<string | undefined>();
	let sourceState = $state<'manual' | 'loading' | 'loaded' | 'unavailable'>(
		requestedCaseId || requestedDocumentId ? 'loading' : 'manual'
	);
	let sourceName = $state('');
	let result = $state<GenerationResultBody | null>(null);
	let isSubmitting = $state(false);
	let isPolling = $state(false);
	let activeMode = $state<'generate' | 'improve' | null>(null);
	let liveStatus = $state('');
	let formError = $state('');
	let fieldErrors = $state<Record<string, string>>({});
	let retryRequest = $state<{ fingerprint: string; clientRequestId: string } | null>(null);
	let pdfFilename = $state('papierbrief.pdf');
	let signaturePNG = $state('');
	let isPreparingPDF = $state(false);
	let pdfError = $state('');
	let preparedPDFBlob = $state<Blob | null>(null);
	let preparedPDFURL = $state('');
	let preparedPDFFilename = $state('');
	let preparedPDFPageCount = $state(0);
	let preparedPDFFingerprint = $state('');
	let uploadPDFBlob = $state<Blob | null>(null);
	let uploadPDFFilename = $state('');
	let uploadPDFSource = $state<'generated' | 'ready' | null>(null);
	let isUploadingPDF = $state(false);
	let uploadStage = $state<'uploading' | 'confirming' | null>(null);
	let uploadError = $state('');
	let confirmedUploadID = $state('');
	let confirmedUploadPages = $state(0);
	let confirmedUploadSize = $state(0);
	let confirmedDeliveryPrice = $state(0);
	let delivery = $state<DeliveryBody | null>(null);
	let deliveryError = $state('');
	let deliveryRetryRequestID = $state('');
	let isCreatingDelivery = $state(false);
	let isPollingDelivery = $state(false);
	const controlsDisabled = $derived(
		isSubmitting ||
			isPolling ||
			isPreparingPDF ||
			isUploadingPDF ||
			isCreatingDelivery ||
			isPollingDelivery
	);
	let pollVersion = 0;
	let deliveryPollVersion = 0;
	let destroyed = false;

	$effect(() => {
		if (!selectedModel && models.length > 0) {
			selectedModel = models[0].id;
		}
	});

	const isUUID = (value: string | null): value is string =>
		!!value &&
		/^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(value);

	const clean = (value?: string) => value?.trim() ?? '';
	const postalCityPattern = /^\d{5}\s+\S.+$/m;
	const streetOrPostOfficeBoxPattern =
		/(^|\n)\s*(?:postfach\s+\d[\d\s-]*|.+\d+[a-z]?(?:[-/]\d+[a-z]?)?)\s*($|\n)/im;
	const validPDFAddress = (value: string) =>
		postalCityPattern.test(value.replaceAll(',', '\n')) &&
		streetOrPostOfficeBoxPattern.test(value.replaceAll(',', '\n'));
	const pdfFingerprint = () =>
		JSON.stringify([
			germanLetter,
			subject,
			senderName,
			senderAddress,
			recipientName,
			recipientAddress,
			reference,
			pdfFilename,
			signaturePNG
		]);
	const clearUploadConfirmation = () => {
		confirmedUploadID = '';
		confirmedUploadPages = 0;
		confirmedUploadSize = 0;
		confirmedDeliveryPrice = 0;
		delivery = null;
		deliveryError = '';
		deliveryRetryRequestID = '';
		deliveryPollVersion += 1;
	};
	const clearPreparedPDF = () => {
		if (preparedPDFURL) {
			URL.revokeObjectURL(preparedPDFURL);
		}
		preparedPDFBlob = null;
		preparedPDFURL = '';
		preparedPDFFilename = '';
		preparedPDFPageCount = 0;
		preparedPDFFingerprint = '';
		if (uploadPDFSource === 'generated') {
			uploadPDFBlob = null;
			uploadPDFFilename = '';
			uploadPDFSource = null;
			clearUploadConfirmation();
		}
	};

	$effect(() => {
		const current = pdfFingerprint();
		if (preparedPDFFingerprint && preparedPDFFingerprint !== current) {
			clearPreparedPDF();
		}
	});
	const formatAddress = (address?: Address) => {
		if (!address) {
			return '';
		}
		const name = clean(address.name).toLocaleLowerCase();
		const structured = [
			[clean(address.street), clean(address.house_number)].filter(Boolean).join(' '),
			[clean(address.postal_code), clean(address.city)].filter(Boolean).join(' '),
			clean(address.state),
			clean(address.country)
		].filter(Boolean);
		const lines = clean(address.address)
			? clean(address.address).split(/\n|,/).map(clean).filter(Boolean)
			: structured;
		const seen: string[] = [];
		return lines
			.filter((line) => {
				const key = line.toLocaleLowerCase();
				if (!key || key === name || seen.includes(key)) {
					return false;
				}
				seen.push(key);
				return true;
			})
			.join('\n');
	};

	const mappedCategory = (value?: string): Category => {
		const match = categories.find(
			(candidate) => candidate.toLocaleLowerCase() === clean(value).toLocaleLowerCase()
		);
		return match ?? 'Other';
	};
	const mappedSubcategory = (value?: string): Subcategory => {
		const match = subcategories.find(
			(candidate) => candidate.toLocaleLowerCase() === clean(value).toLocaleLowerCase()
		);
		return match ?? 'freie_formulierung';
	};

	const applySource = (metadata: Metadata) => {
		const outgoingSender = metadata.receiver;
		const outgoingRecipient = metadata.sender;
		senderName = clean(outgoingSender?.name);
		senderAddress = formatAddress(outgoingSender);
		recipientName = clean(outgoingRecipient?.name);
		recipientAddress = formatAddress(outgoingRecipient);
		category = mappedCategory(outgoingRecipient?.category);
		subcategory = mappedSubcategory(metadata.letter_type);
		subject = clean(metadata.subject?.original || metadata.subject?.localized);
		reference = clean(metadata.reference_number);
	};

	const loadSource = async () => {
		if (!isUUID(requestedCaseId) || !isUUID(requestedDocumentId)) {
			if (requestedCaseId || requestedDocumentId) {
				sourceState = 'unavailable';
			}
			return;
		}
		try {
			const detail = await getAnalysisQuery({ id: requestedCaseId });
			const document = detail.documents.find(
				(candidate) =>
					candidate.documentId === requestedDocumentId && candidate.analysis.status === 'completed'
			);
			if (!document?.analysis.metadata) {
				sourceState = 'unavailable';
				return;
			}
			applySource(document.analysis.metadata);
			sourceCaseId = requestedCaseId;
			sourceDocumentId = requestedDocumentId;
			sourceName = document.fileName;
			sourceState = 'loaded';
		} catch {
			sourceState = 'unavailable';
		}
	};

	onMount(() => {
		void loadSource();
		return () => {
			destroyed = true;
			pollVersion += 1;
			deliveryPollVersion += 1;
			if (preparedPDFURL) {
				URL.revokeObjectURL(preparedPDFURL);
			}
		};
	});

	const validate = (mode: 'generate' | 'improve') => {
		const errors: Record<string, string> = {};
		if (!selectedModel) {
			errors.model = m.papierbrief_validation_model();
		}
		if (mode === 'improve') {
			if (!clean(germanLetter)) {
				errors.currentLetter = m.papierbrief_validation_current_letter();
			}
			if (!clean(improvementFeedback)) {
				errors.feedback = m.papierbrief_validation_feedback();
			}
			fieldErrors = errors;
			return Object.keys(errors).length === 0;
		}
		if (!clean(senderName)) {
			errors.senderName = m.papierbrief_validation_sender_name();
		}
		if (!clean(senderAddress)) {
			errors.senderAddress = m.papierbrief_validation_sender_address();
		}
		if (!clean(recipientName)) {
			errors.recipientName = m.papierbrief_validation_recipient_name();
		}
		if (!clean(recipientAddress))
			errors.recipientAddress = m.papierbrief_validation_recipient_address();
		if (
			!clean(subject) &&
			!clean(reference) &&
			!clean(instruction) &&
			subcategory === 'freie_formulierung'
		) {
			errors.context = m.papierbrief_validation_context();
		}
		fieldErrors = errors;
		return Object.keys(errors).length === 0;
	};

	const errorMessage = (code: string) => {
		switch (code) {
			case 'unauthorized':
				return m.papierbrief_error_unauthorized();
			case 'insufficient_points':
				return m.papierbrief_error_insufficient_points();
			case 'invalid_source':
				return m.papierbrief_error_invalid_source();
			case 'source_not_ready':
				return m.papierbrief_error_source_not_ready();
			case 'invalid_input':
				return m.papierbrief_error_invalid_input();
			case 'invalid_generation_input':
				return m.papierbrief_error_invalid_input();
			case 'unsupported_model':
				return m.papierbrief_error_model();
			case 'invalid_generation_output':
				return m.papierbrief_error_output();
			default:
				return m.papierbrief_error_generic();
		}
	};

	const currentRequest = (
		mode: 'generate' | 'improve'
	): Omit<CreateGenerationRequest, 'clientRequestId'> => ({
		mode,
		model: selectedModel,
		locale,
		sourceCaseId,
		sourceDocumentId,
		input: {
			sender: { name: senderName, address: senderAddress },
			recipient: { name: recipientName, address: recipientAddress },
			category,
			subcategory,
			subject,
			reference,
			instruction,
			currentGermanLetter: mode === 'improve' ? germanLetter : undefined,
			feedback: mode === 'improve' ? improvementFeedback : undefined
		}
	});

	const submitGeneration = async (mode: 'generate' | 'improve') => {
		if (controlsDisabled || !validate(mode)) {
			if (Object.keys(fieldErrors).length > 0) {
				liveStatus = m.papierbrief_validation_summary();
			}
			return;
		}

		const request = currentRequest(mode);
		const fingerprint = JSON.stringify(request);
		const clientRequestId =
			retryRequest?.fingerprint === fingerprint
				? retryRequest.clientRequestId
				: crypto.randomUUID();
		const version = ++pollVersion;
		let requestAccepted = false;
		isSubmitting = true;
		activeMode = mode;
		formError = '';
		if (mode === 'generate') {
			result = null;
			germanLetter = '';
		}
		liveStatus =
			mode === 'improve'
				? m.papierbrief_improve_status_submitting()
				: m.papierbrief_status_submitting();

		try {
			const response = await createGenerationCommand({ ...request, clientRequestId });
			if (!response.ok) {
				retryRequest = { fingerprint, clientRequestId };
				if (response.error === 'invalid_source') {
					sourceCaseId = undefined;
					sourceDocumentId = undefined;
					sourceState = 'unavailable';
				}
				formError = errorMessage(response.error);
				liveStatus = formError;
				toast.error(formError);
				return;
			}

			requestAccepted = true;
			retryRequest = null;
			isSubmitting = false;
			isPolling = true;
			liveStatus =
				mode === 'improve' ? m.papierbrief_improve_status_queued() : m.papierbrief_status_queued();
			await pointsQuery.refresh();

			const completed = await pollGeneration(
				response.generation.generationId,
				() => destroyed || version !== pollVersion,
				(value) => {
					liveStatus =
						value.status === 'processing'
							? mode === 'improve'
								? m.papierbrief_improve_status_processing()
								: m.papierbrief_status_processing()
							: mode === 'improve'
								? m.papierbrief_improve_status_queued()
								: m.papierbrief_status_queued();
				}
			);
			if (completed.status === 'completed' && completed.result) {
				result = completed.result;
				germanLetter = completed.result.germanLetter;
				if (mode === 'improve') {
					improvementFeedback = '';
				}
				liveStatus =
					mode === 'improve'
						? m.papierbrief_improve_status_completed()
						: m.papierbrief_status_completed();
				toast.success(liveStatus);
			} else {
				const reason = errorMessage(completed.errorMessage ?? 'generation_failed');
				formError = completed.refunded ? `${reason} ${m.papierbrief_refund_completed()}` : reason;
				liveStatus = formError;
				toast.error(formError);
				if (completed.refunded) {
					await pointsQuery.refresh();
				}
			}
		} catch (error) {
			if (error instanceof GenerationCancelledError) {
				return;
			}
			if (!requestAccepted) {
				retryRequest = { fingerprint, clientRequestId };
			}
			formError =
				error instanceof GenerationTimedOutError
					? m.papierbrief_status_timed_out()
					: m.papierbrief_error_generic();
			liveStatus = formError;
			toast.error(formError);
		} finally {
			if (version === pollVersion) {
				isSubmitting = false;
				isPolling = false;
				activeMode = null;
			}
		}
	};

	const generate = () => submitGeneration('generate');
	const improve = () => submitGeneration('improve');

	const validatePDF = () => {
		const errors: Record<string, string> = {};
		if (!clean(germanLetter)) {
			errors.currentLetter = m.papierbrief_validation_current_letter();
		}
		if (!clean(subject)) {
			errors.pdfSubject = m.papierbrief_pdf_validation_subject();
		}
		if (!clean(senderName)) {
			errors.senderName = m.papierbrief_validation_sender_name();
		}
		if (!validPDFAddress(senderAddress))
			errors.senderAddress = m.papierbrief_pdf_validation_address();
		if (!clean(recipientName)) {
			errors.recipientName = m.papierbrief_validation_recipient_name();
		}
		if (!validPDFAddress(recipientAddress))
			errors.recipientAddress = m.papierbrief_pdf_validation_address();
		if (!clean(pdfFilename)) {
			errors.pdfFilename = m.papierbrief_pdf_validation_filename();
		}
		fieldErrors = errors;
		return Object.keys(errors).length === 0;
	};

	const preparePDF = async () => {
		if (controlsDisabled || !validatePDF()) {
			if (Object.keys(fieldErrors).length > 0) {
				pdfError = m.papierbrief_pdf_validation_summary();
			}
			return;
		}

		const request: CreatePDFRequest = {
			germanText: germanLetter,
			subject,
			sender: { name: senderName, address: senderAddress },
			recipient: { name: recipientName, address: recipientAddress },
			reference,
			filename: pdfFilename,
			signaturePngBase64: signaturePNG || undefined
		};
		const fingerprint = pdfFingerprint();
		isPreparingPDF = true;
		pdfError = '';
		try {
			const response = await createPDFCommand(request);
			if (!response.ok) {
				pdfError =
					response.error === 'unauthorized'
						? m.papierbrief_pdf_error_unauthorized()
						: response.error === 'invalid_input'
							? m.papierbrief_pdf_error_invalid_input()
							: m.papierbrief_pdf_error_generic();
				toast.error(pdfError);
				return;
			}

			const blob = new Blob([response.bytes.slice()], { type: 'application/pdf' });
			const nextURL = URL.createObjectURL(blob);
			clearPreparedPDF();
			preparedPDFBlob = blob;
			preparedPDFURL = nextURL;
			preparedPDFFilename = response.filename;
			preparedPDFPageCount = response.pageCount;
			preparedPDFFingerprint = fingerprint;
			uploadPDFBlob = blob;
			uploadPDFFilename = response.filename;
			uploadPDFSource = 'generated';
			uploadError = '';
			clearUploadConfirmation();
			toast.success(m.papierbrief_pdf_status_ready());
		} catch {
			pdfError = m.papierbrief_pdf_error_generic();
			toast.error(pdfError);
		} finally {
			isPreparingPDF = false;
		}
	};

	const downloadPDF = () => {
		if (!preparedPDFBlob || !preparedPDFURL) {
			return;
		}
		const link = document.createElement('a');
		link.href = preparedPDFURL;
		link.download = preparedPDFFilename;
		link.click();
	};

	const selectReadyPDF = (event: Event) => {
		const input = event.currentTarget as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;
		clearUploadConfirmation();
		if (!paperLetterPDFWithinLimit(file)) {
			uploadPDFBlob = null;
			uploadPDFFilename = '';
			uploadPDFSource = null;
			uploadError = m.papierbrief_upload_error_too_large();
			return;
		}
		uploadPDFBlob = file;
		uploadPDFFilename = file.name;
		uploadPDFSource = 'ready';
		uploadError = '';
	};

	const uploadErrorMessage = (code: string) => {
		switch (code) {
			case 'unauthorized':
				return m.papierbrief_pdf_error_unauthorized();
			case 'missing_object':
				return m.papierbrief_upload_error_missing();
			case 'too_large':
				return m.papierbrief_upload_error_too_large();
			case 'invalid_pdf':
			case 'invalid_filename':
				return m.papierbrief_pdf_error_invalid_input();
			default:
				return m.papierbrief_pdf_error_generic();
		}
	};

	const uploadSelectedPDF = async () => {
		if (controlsDisabled || !uploadPDFBlob) {
			uploadError = m.papierbrief_upload_error_select();
			return;
		}
		isUploadingPDF = true;
		uploadError = '';
		try {
			uploadStage = 'uploading';
			const created = await createPaperLetterUploadCommand({ filename: uploadPDFFilename });
			if (!created.ok) {
				uploadError = uploadErrorMessage(created.error);
				return;
			}
			if (!(await putPaperLetterPDF(created.uploadUrl, uploadPDFBlob))) {
				uploadError = m.papierbrief_upload_error_put();
				return;
			}
			uploadStage = 'confirming';
			const confirmed = await confirmPaperLetterUploadCommand({ uploadId: created.uploadId });
			if (!confirmed.ok) {
				uploadError = uploadErrorMessage(confirmed.error);
				return;
			}
			confirmedUploadID = confirmed.uploadId;
			confirmedUploadPages = confirmed.pageCount;
			confirmedUploadSize = confirmed.sizeBytes;
			confirmedDeliveryPrice = confirmed.deliveryPrice;
			toast.success(m.papierbrief_upload_confirmed());
		} catch {
			uploadError = m.papierbrief_pdf_error_generic();
		} finally {
			uploadStage = null;
			isUploadingPDF = false;
			if (uploadError) toast.error(uploadError);
		}
	};

	const deliveryErrorMessage = (code: string) => {
		switch (code) {
			case 'unauthorized':
				return m.papierbrief_pdf_error_unauthorized();
			case 'insufficient_points':
				return m.papierbrief_error_insufficient_points();
			case 'upload_not_found':
				return m.papierbrief_upload_error_missing();
			case 'upload_not_confirmed':
			case 'invalid_input':
				return m.papierbrief_pdf_error_invalid_input();
			default:
				return m.papierbrief_delivery_error_generic();
		}
	};

	const submitDelivery = async () => {
		if (controlsDisabled || !confirmedUploadID) return;
		const clientRequestId = deliveryRetryRequestID || crypto.randomUUID();
		const version = ++deliveryPollVersion;
		let accepted = false;
		isCreatingDelivery = true;
		deliveryError = '';
		try {
			const response = await createPaperLetterDeliveryCommand({
				uploadId: confirmedUploadID,
				clientRequestId
			});
			if (!response.ok) {
				deliveryRetryRequestID = response.error === 'request_failed' ? clientRequestId : '';
				deliveryError = deliveryErrorMessage(response.error);
				toast.error(deliveryError);
				return;
			}
			accepted = true;
			deliveryRetryRequestID = '';
			delivery = response.delivery;
			isCreatingDelivery = false;
			isPollingDelivery = true;
			await pointsQuery.refresh();
			const completed = await pollDelivery(
				response.delivery.deliveryId,
				() => destroyed || version !== deliveryPollVersion,
				(value) => (delivery = value)
			);
			delivery = completed;
			if (completed.status === 'failed') {
				await pointsQuery.refresh();
			}
		} catch (error) {
			if (error instanceof DeliveryCancelledError) return;
			if (!accepted) deliveryRetryRequestID = clientRequestId;
			deliveryError =
				error instanceof DeliveryTimedOutError
					? m.papierbrief_delivery_timed_out()
					: m.papierbrief_delivery_error_generic();
			toast.error(deliveryError);
		} finally {
			if (version === deliveryPollVersion) {
				isCreatingDelivery = false;
				isPollingDelivery = false;
			}
		}
	};

	const formattedUploadSize = $derived(
		new Intl.NumberFormat(locale, { maximumFractionDigits: 2 }).format(
			confirmedUploadSize / 1024 / 1024
		)
	);
</script>

<svelte:head><title>{m.papierbrief_page_title()} - BescheidBuddy</title></svelte:head>

<div dir={locale === 'ar' ? 'rtl' : 'ltr'}>
	<PageHeading title={m.papierbrief_page_title()} />

	<Section>
		<p class="text-[#b7b7c8]">{m.papierbrief_intro()}</p>
		<div class="mt-4 rounded-xl border border-white/10 bg-white/5 p-4" aria-live="polite">
			<h2 class="font-semibold text-white">{m.papierbrief_source_title()}</h2>
			{#if sourceState === 'loading'}
				<p class="mt-1 text-sm text-[#b7b7c8]">{m.papierbrief_source_loading()}</p>
			{:else if sourceState === 'loaded'}
				<p class="mt-1 text-sm text-emerald-300">
					{m.papierbrief_source_loaded({ name: sourceName })}
				</p>
			{:else if sourceState === 'unavailable'}
				<p class="mt-1 text-sm text-amber-200">{m.papierbrief_source_unavailable()}</p>
			{:else}
				<p class="mt-1 text-sm text-[#b7b7c8]">{m.papierbrief_source_manual()}</p>
			{/if}
		</div>
	</Section>

	<Section>
		<h2 class="mb-5 text-2xl font-semibold text-white">{m.papierbrief_form_title()}</h2>
		<div class="grid gap-6 lg:grid-cols-2">
			<fieldset class="grid gap-4 rounded-xl border border-white/10 p-4">
				<legend class="px-2 font-semibold text-white">{m.papierbrief_sender_title()}</legend>
				<InputField
					name="sender-name"
					type="text"
					label={m.papierbrief_name_label()}
					bind:value={senderName}
					maxlength={200}
					autocomplete="name"
					error={fieldErrors.senderName}
					disabled={controlsDisabled}
				/>
				<Textarea
					name="sender-address"
					rows={4}
					label={m.papierbrief_address_label()}
					bind:value={senderAddress}
					maxlength={1000}
					error={fieldErrors.senderAddress}
					disabled={controlsDisabled}
				/>
			</fieldset>
			<fieldset class="grid gap-4 rounded-xl border border-white/10 p-4">
				<legend class="px-2 font-semibold text-white">{m.papierbrief_recipient_title()}</legend>
				<InputField
					name="recipient-name"
					type="text"
					label={m.papierbrief_name_label()}
					bind:value={recipientName}
					maxlength={200}
					autocomplete="organization"
					error={fieldErrors.recipientName}
					disabled={controlsDisabled}
				/>
				<Textarea
					name="recipient-address"
					rows={4}
					label={m.papierbrief_address_label()}
					bind:value={recipientAddress}
					maxlength={1000}
					error={fieldErrors.recipientAddress}
					disabled={controlsDisabled}
				/>
			</fieldset>
		</div>

		<div class="mt-6 grid gap-5 md:grid-cols-2">
			<div class="flex flex-col gap-2">
				<label for="category" class="text-xs font-extrabold tracking-widest uppercase"
					>{m.papierbrief_category_label()}</label
				>
				<select
					id="category"
					bind:value={category}
					disabled={controlsDisabled}
					class="rounded-xl border border-white/22 bg-white/5 disabled:opacity-50"
				>
					{#each categories as value (value)}<option class="bg-[#1a1a24]" {value}
							>{categoryLabels[value]()}</option
						>{/each}
				</select>
			</div>
			<div class="flex flex-col gap-2">
				<label for="subcategory" class="text-xs font-extrabold tracking-widest uppercase"
					>{m.papierbrief_subcategory_label()}</label
				>
				<select
					id="subcategory"
					bind:value={subcategory}
					disabled={controlsDisabled}
					class="rounded-xl border border-white/22 bg-white/5 disabled:opacity-50"
				>
					{#each subcategories as value (value)}<option class="bg-[#1a1a24]" {value}
							>{subcategoryLabels[value]()}</option
						>{/each}
				</select>
			</div>
			<InputField
				name="subject"
				type="text"
				label={m.papierbrief_subject_label()}
				bind:value={subject}
				maxlength={500}
				error={fieldErrors.pdfSubject}
				disabled={controlsDisabled}
			/>
			<InputField
				name="reference"
				type="text"
				label={m.papierbrief_reference_label()}
				bind:value={reference}
				maxlength={200}
				disabled={controlsDisabled}
			/>
		</div>

		<div class="mt-5">
			<Textarea
				name="instruction"
				rows={5}
				label={m.papierbrief_instruction_label()}
				bind:value={instruction}
				maxlength={2500}
				error={fieldErrors.context}
				disabled={controlsDisabled}
			/>
			<p class="mt-2 text-sm text-[#b7b7c8]">{m.papierbrief_instruction_help()}</p>
		</div>

		<ModelSelector
			id="papierbrief-model"
			label={m.papierbrief_model_label()}
			{models}
			bind:value={selectedModel}
			disabled={controlsDisabled}
		/>
		{#if fieldErrors.model}
			<p class="mt-2 text-sm text-red-300">{fieldErrors.model}</p>
		{/if}

		<div class="mt-7 flex flex-wrap items-center justify-between gap-4">
			<div>
				<p class="font-semibold text-white">{m.papierbrief_price({ points: 10 })}</p>
				<p class="text-sm text-[#b7b7c8]">
					{m.papierbrief_balance({
						points: buddyPoints.monthly_points + buddyPoints.lifetime_points
					})}
				</p>
			</div>
			<PrimaryButton
				type="button"
				onclick={generate}
				disabled={controlsDisabled}
				text={activeMode === 'generate'
					? m.papierbrief_generating_button()
					: m.papierbrief_generate_button()}
			/>
		</div>
	</Section>

	<p class="sr-only" aria-live="assertive">{liveStatus}</p>
	{#if liveStatus && (isSubmitting || isPolling)}
		<p
			role="status"
			class="mb-6 rounded-lg border border-primary/25 bg-primary/10 px-3 py-2 text-sm text-primary"
		>
			{liveStatus}
		</p>
	{/if}
	{#if formError}
		<p
			role="alert"
			class="mb-6 rounded-lg border border-red-400/30 bg-red-400/10 px-3 py-2 text-sm text-red-300"
		>
			{formError}
		</p>
	{/if}

	{#if result}
		<Section>
			<h2 class="text-2xl font-semibold text-white">{m.papierbrief_result_title()}</h2>
			{#if result.subject}
				<p class="mt-4 text-sm text-[#b7b7c8]">
					<strong class="text-white">{m.papierbrief_result_subject()}:</strong>
					{result.subject}
				</p>
			{/if}
			<div class="mt-5 grid gap-5 lg:grid-cols-2">
				<div class="rounded-xl border border-white/10 bg-white/5 p-4">
					<h3 class="font-semibold text-white">{m.papierbrief_result_german()}</h3>
					<div class="mt-3" dir="ltr">
						<Textarea
							name="german-letter"
							rows={18}
							label={m.papierbrief_editor_label()}
							bind:value={germanLetter}
							maxlength={24000}
							error={fieldErrors.currentLetter}
							disabled={controlsDisabled}
						/>
						<p class="mt-2 text-sm text-[#b7b7c8]">{m.papierbrief_editor_help()}</p>
					</div>
				</div>
				<div class="rounded-xl border border-white/10 bg-white/5 p-4">
					<h3 class="font-semibold text-white">{m.papierbrief_result_translation()}</h3>
					<p class="mt-3 text-sm leading-6 whitespace-pre-wrap text-[#e5e5ef]">
						{result.userLanguageLetter}
					</p>
				</div>
			</div>
			<div class="mt-5 rounded-xl border border-white/10 bg-white/5 p-4">
				<h3 class="font-semibold text-white">{m.papierbrief_result_explanation()}</h3>
				<p class="mt-2 text-sm text-[#b7b7c8]">{result.explanation}</p>
			</div>

			<div class="mt-5 rounded-xl border border-white/10 bg-white/5 p-4">
				<h3 class="text-xl font-semibold text-white">{m.papierbrief_improve_title()}</h3>
				<div class="mt-4">
					<Textarea
						name="improvement-feedback"
						rows={5}
						label={m.papierbrief_feedback_label()}
						bind:value={improvementFeedback}
						maxlength={2500}
						error={fieldErrors.feedback}
						disabled={controlsDisabled}
					/>
					<p class="mt-2 text-sm text-[#b7b7c8]">{m.papierbrief_feedback_help()}</p>
				</div>
				<div class="mt-6 flex flex-wrap items-center justify-between gap-4">
					<div>
						<p class="font-semibold text-white">{m.papierbrief_improve_price({ points: 10 })}</p>
						<p class="text-sm text-[#b7b7c8]">
							{m.papierbrief_balance({
								points: buddyPoints.monthly_points + buddyPoints.lifetime_points
							})}
						</p>
					</div>
					<PrimaryButton
						type="button"
						onclick={improve}
						disabled={controlsDisabled}
						text={activeMode === 'improve'
							? m.papierbrief_improving_button()
							: m.papierbrief_improve_button()}
					/>
				</div>
			</div>

			<div class="mt-5 rounded-xl border border-white/10 bg-white/5 p-4">
				<h3 class="text-xl font-semibold text-white">{m.papierbrief_pdf_title()}</h3>
				<p class="mt-2 text-sm text-[#b7b7c8]">{m.papierbrief_pdf_description()}</p>
				<p class="mt-1 text-sm font-semibold text-emerald-300">{m.papierbrief_pdf_free()}</p>

				<div class="mt-5">
					<InputField
						name="pdf-filename"
						type="text"
						label={m.papierbrief_pdf_filename_label()}
						bind:value={pdfFilename}
						maxlength={200}
						error={fieldErrors.pdfFilename}
						disabled={controlsDisabled}
					/>
				</div>

				<div class="mt-5">
					<SignatureCanvas
						label={m.papierbrief_signature_label()}
						description={m.papierbrief_signature_description()}
						clearLabel={m.papierbrief_signature_clear()}
						bind:value={signaturePNG}
						disabled={controlsDisabled}
					/>
				</div>

				<div class="mt-6 flex flex-wrap items-center gap-4">
					<PrimaryButton
						type="button"
						onclick={preparePDF}
						disabled={controlsDisabled}
						text={isPreparingPDF
							? m.papierbrief_pdf_preparing_button()
							: m.papierbrief_pdf_prepare_button()}
					/>
					{#if preparedPDFURL}
						<PrimaryButton
							type="button"
							onclick={downloadPDF}
							text={m.papierbrief_pdf_download_button()}
						/>
					{/if}
				</div>

				<p class="sr-only" aria-live="assertive">
					{isPreparingPDF ? m.papierbrief_pdf_status_preparing() : pdfError}
				</p>
				{#if isPreparingPDF}
					<p role="status" class="mt-4 text-sm text-primary">
						{m.papierbrief_pdf_status_preparing()}
					</p>
				{/if}
				{#if pdfError}
					<p role="alert" class="mt-4 text-sm text-red-300">{pdfError}</p>
				{/if}
				{#if preparedPDFURL}
					<p class="mt-4 text-sm text-emerald-300" role="status">
						{m.papierbrief_pdf_status_ready_pages({ pages: preparedPDFPageCount })}
					</p>
				{/if}
			</div>
		</Section>
	{/if}

	<Section>
		<h2 class="text-2xl font-semibold text-white">{m.papierbrief_upload_title()}</h2>
		<p class="mt-2 text-sm text-[#b7b7c8]">{m.papierbrief_upload_description()}</p>
		<p class="mt-1 text-sm text-[#b7b7c8]">
			{m.papierbrief_upload_limit({ size: maxPaperLetterPDFBytes / 1024 / 1024 })}
		</p>
		<p class="mt-1 text-sm text-[#b7b7c8]">{m.papierbrief_delivery_pricing()}</p>

		<div class="mt-5">
			<label for="papierbrief-ready-pdf" class="text-xs font-extrabold tracking-widest uppercase">
				{m.papierbrief_upload_file_label()}
			</label>
			<input
				id="papierbrief-ready-pdf"
				type="file"
				accept="application/pdf,.pdf"
				disabled={controlsDisabled}
				onchange={selectReadyPDF}
				class="mt-2 block w-full rounded-xl border border-white/22 bg-white/5 p-3 text-sm disabled:opacity-50"
			/>
		</div>

		{#if uploadPDFBlob}
			<p class="mt-4 text-sm text-white">
				{m.papierbrief_upload_selected({ name: uploadPDFFilename })}
			</p>
		{/if}
		<div class="mt-5">
			<PrimaryButton
				type="button"
				onclick={uploadSelectedPDF}
				disabled={controlsDisabled || !uploadPDFBlob || !!confirmedUploadID}
				text={isUploadingPDF
					? m.papierbrief_upload_uploading_button()
					: m.papierbrief_upload_button()}
			/>
		</div>

		<p class="sr-only" aria-live="assertive">
			{uploadStage === 'confirming'
				? m.papierbrief_upload_confirming()
				: uploadStage === 'uploading'
					? m.papierbrief_upload_uploading()
					: uploadError}
		</p>
		{#if uploadStage}
			<p role="status" class="mt-4 text-sm text-primary">
				{uploadStage === 'confirming'
					? m.papierbrief_upload_confirming()
					: m.papierbrief_upload_uploading()}
			</p>
		{/if}
		{#if uploadError}
			<p role="alert" class="mt-4 text-sm text-red-300">{uploadError}</p>
		{/if}
		{#if confirmedUploadID}
			<p role="status" class="mt-4 text-sm text-emerald-300">
				{m.papierbrief_upload_confirmed_details({
					pages: confirmedUploadPages,
					size: formattedUploadSize,
					points: confirmedDeliveryPrice
				})}
			</p>
			{#if !delivery}
				<p class="mt-4 text-sm text-[#b7b7c8]">
					{m.papierbrief_delivery_confirmation({ points: confirmedDeliveryPrice })}
				</p>
				<div class="mt-4">
					<PrimaryButton
						type="button"
						onclick={submitDelivery}
						disabled={controlsDisabled}
						text={isCreatingDelivery
							? m.papierbrief_delivery_creating_button()
							: m.papierbrief_delivery_confirm_button({ points: confirmedDeliveryPrice })}
					/>
				</div>
			{:else if delivery.status === 'sent'}
				<p role="status" class="mt-4 text-sm text-emerald-300">
					{m.papierbrief_delivery_sent()}
				</p>
			{:else if delivery.status === 'failed'}
				<p role="alert" class="mt-4 text-sm text-red-300">{m.papierbrief_delivery_failed()}</p>
			{:else if delivery.status === 'manual_review'}
				<p role="status" class="mt-4 text-sm text-amber-200">
					{m.papierbrief_delivery_manual_review()}
				</p>
			{:else}
				<p role="status" class="mt-4 text-sm text-primary">
					{delivery.status === 'processing'
						? m.papierbrief_delivery_processing()
						: m.papierbrief_delivery_queued()}
				</p>
			{/if}
		{/if}
		{#if deliveryError}
			<p role="alert" class="mt-4 text-sm text-red-300">{deliveryError}</p>
		{/if}
	</Section>
</div>
