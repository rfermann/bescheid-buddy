<script>
	import { m } from '$lib/paraglide/messages.js';
	import { submitContact } from '$lib/remote/contact.remote';
	import { isHttpError } from '@sveltejs/kit';

	let { type = 'general', key = 'default' } = $props();

	const formAction = submitContact.for(key);
</script>

<form {...formAction.enhance(async ({ submit }) => {
	submit()
		.then(() => {
			alert('success');
		})
		.catch((err) => {
			console.log(err);
			if (isHttpError(err)) {
				alert(err.body.message);
			}
		});
})} class="space-y-4" oninput={() => formAction.validate()}>
	<div>
		<label for="name" class="block text-sm font-medium mb-1 text-gray-200">
			{m.kontakt_form_name()}
		</label>
		<input
			{...formAction.fields.name.as('text')}
			class="w-full px-3 py-2 rounded-xl border border-accent bg-gray-800/50 text-gray-200 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-400 focus:border-purple-400 backdrop-blur-sm {formAction.fields.name.issues()?.length ? 'border-red-400' : ''}"
			placeholder={m.kontakt_form_name()}
		/>
		{#each formAction.fields.name.issues() as issue}
			<p class="mt-1 text-sm text-red-400">{issue.message}</p>
		{:else }
			<p class="mt-1 text-sm text-red-400">&nbsp;</p>
		{/each}
	</div>

	<div>
		<label for="email" class="block text-sm font-medium mb-1 text-gray-200">
			{m.kontakt_form_email()}
		</label>
		<input
			{...formAction.fields.email.as('email')}
			class="w-full px-3 py-2 rounded-xl border border-accent bg-gray-800/50 text-gray-200 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-400 focus:border-purple-400 backdrop-blur-sm {formAction.fields.email.issues()?.length ? 'border-red-400' : ''}"
			placeholder={m.kontakt_form_email()}
		/>
		{#each formAction.fields.email.issues() as issue}
			<p class="mt-1 text-sm text-red-400">{issue.message}</p>
		{:else }
			<p class="mt-1 text-sm text-red-400">&nbsp;</p>
		{/each}
	</div>

	<div>
		<label for="message" class="block text-sm font-medium mb-1 text-gray-200">
			{m.kontakt_form_message()}
		</label>
		<textarea
			{...formAction.fields.message.as("text")}
			rows="5"
			class="w-full px-3 py-2 rounded-xl border border-accent bg-gray-800/50 text-gray-200 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-400 focus:border-purple-400 backdrop-blur-sm resize-none {formAction.fields.message.issues()?.length ? 'border-red-400' : ''}"
			placeholder={m.kontakt_form_message()}
		></textarea>
		{#each formAction.fields.message.issues() as issue}
			<p class="mt-1 text-sm text-red-400">{issue.message}</p>
		{:else }
			<p class="mt-1 text-sm text-red-400">&nbsp;</p>
		{/each}
	</div>

	<input type="hidden" name="type" value={type} />

	<button
		type="submit"
		disabled={formAction.pending !== 0}
		class="w-full px-4 py-3.5 rounded-2xl border border-purple-400 bg-violet-500/20 text-xl font-extrabold text-purple-200 hover:bg-violet-500/30 disabled:bg-gray-600/20 disabled:border-gray-500 disabled:text-gray-400 transition-all duration-300 shadow-[0_4px_15px_rgba(139,92,246,0.3)] hover:shadow-[0_6px_20px_rgba(139,92,246,0.4)] disabled:shadow-none flex items-center justify-center gap-2"
	>
		{#if formAction.pending}
			<svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor"
							d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			Sending...
		{:else}
			{m.kontakt_form_submit()}
		{/if}
	</button>
</form>
