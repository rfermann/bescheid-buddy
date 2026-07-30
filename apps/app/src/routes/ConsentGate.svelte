<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';
	import LegalConsentModal from '$lib/components/LegalConsentModal.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';

	const CONSENT_DISMISSED_KEY = 'analyse-legal-consent-dismissed';

	let {
		canProceed = true,
		onConfirm
	}: { canProceed?: boolean; onConfirm: () => void } = $props();
	let showConsent = $state(false);

	const handleClick = () => {
		if (!canProceed) return;
		if (localStorage.getItem(CONSENT_DISMISSED_KEY) === 'true') {
			onConfirm();
			return;
		}
		showConsent = true;
	};

	const handleConsentConfirm = (dontShowAgain: boolean) => {
		if (dontShowAgain) {
			localStorage.setItem(CONSENT_DISMISSED_KEY, 'true');
		}
		showConsent = false;
		onConfirm();
	};
</script>

<PrimaryButton
	type="button"
	onclick={handleClick}
	text={m.analyse_analyze_button()}
/>

<LegalConsentModal
	bind:open={showConsent}
	onConfirm={handleConsentConfirm}
	onCancel={() => (showConsent = false)}
/>
