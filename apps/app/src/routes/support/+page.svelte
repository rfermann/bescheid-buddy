<script>
	import { m } from '$lib/paraglide/messages';
	import PageHeading from '$lib/components/PageHeading.svelte';
	import Section from '$lib/components/Section.svelte';
	import InputField from '$lib/components/input/InputField.svelte';
	import Textarea from '$lib/components/input/Textarea.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import SupportCasesTable from './SupportCasesTable.svelte';
	import { createSupportCaseForm } from './support.remote.ts';
</script>
<svelte:head>
	<title>{m.support_page_title()} - BescheidBuddy</title>
</svelte:head>

<PageHeading title={m.support_page_title()} />

<Section>
	<div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6 mb-6">
		<div>
			<h2 class="text-2xl font-semibold mb-2 text-[#e9e9f1]">{m.support_contact_heading()}</h2>
			<p class="text-sm text-[#b7b7c8]">{m.support_contact_subtitle()}</p>
		</div>
	</div>
	<form {...createSupportCaseForm}>
		<div class="flex flex-col gap-6">
			<InputField
				type="text"
				label={m.support_subject_label()}
				{...createSupportCaseForm.fields.subject.as("text", "")} />
			<Textarea rows={4} label={m.support_message_label()} {...createSupportCaseForm.fields.message.as("text", "")} />
			<PrimaryButton type="submit" text={m.support_send_button()} />
		</div>
	</form>
</Section>

<SupportCasesTable />