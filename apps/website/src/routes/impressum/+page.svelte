<script>
	import CardHeader from '$lib/components/ui/card/CardHeader.svelte';
	import Card from '$lib/components/ui/card/Card.svelte';
	import CardContent from '$lib/components/ui/card/CardContent.svelte';
	import ContactForm from '$lib/components/ContactForm.svelte';
	import { m } from '$lib/paraglide/messages';
</script>

<Card>
	<CardHeader title={m.kontakt_title()} />
	<CardContent>
		<div class="space-y-6">
			<section>
				<h2 class="text-xl font-semibold mb-3 text-gray-200">Angaben gemäß § 5 TMG:</h2>
				<p class="text-gray-300">Anbieter: MachKInfach Media & Digital Services</p>
				<p class="text-gray-300">Inh. Shamil Somkhishvili</p>
				<p class="text-gray-300">Brückenstrasse 3</p>
				<p class="text-gray-300">14913 Jüterbog</p>
				<p class="text-gray-300">Deutschland</p>
			</section>

			<section>
				<h2 class="text-xl font-semibold mb-3 text-gray-200">Kontakt:</h2>
				<p class="text-gray-300">+49 15562 678031</p>
				<p class="text-gray-300">
					E-Mail: <a href="mailto:kontakt@machkinfach.de" class="text-purple-400 hover:text-purple-300">kontakt@machkinfach.de</a>
				</p>
			</section>

			<section>
				<h2 class="text-xl font-semibold mb-3 text-gray-200">Verantwortlich für den Inhalt nach § 55 Abs. 2 RStV:</h2>
				<p class="text-gray-300">Shamil Somkhishvili</p>
				<p class="text-gray-300">Brückenstrasse 3</p>
				<p class="text-gray-300">14913 Jüterbog</p>
			</section>

			<section>
				<h2 class="text-xl font-semibold mb-3 text-secondary">Hinweis zur Tätigkeit</h2>
				<p class="text-gray-300">
					Die Plattform BescheidBuddy ist ein digitales Assistenzsystem zur Analyse von Dokumenten unter Einsatz
					künstlicher Intelligenz.
					Es erfolgt keine Rechtsberatung im Sinne des Rechtsdienstleistungsgesetzes (RDG).
				</p>
			</section>

			<section>
				<h2 class="text-xl font-semibold mb-3 text-secondary">Streitschlichtung</h2>
				<p class="text-gray-300 mb-3">
					Die Europäische Kommission stellt eine Plattform zur Online-Streitbeilegung (OS) bereit:
					<a href="https://ec.europa.eu/consumers/odr" target="_blank" rel="noopener"
						 class="text-purple-400 hover:text-purple-300">https://ec.europa.eu/consumers/odr</a>
				</p>
				<p class="text-gray-300 mb-3">Unsere E-Mail-Adresse finden Sie oben im Impressum.</p>
				<p class="text-gray-300">
					Wir sind nicht bereit oder verpflichtet, an Streitbeilegungsverfahren vor einer Verbraucherschlichtungsstelle
					teilzunehmen.
				</p>
			</section>
		</div>
	</CardContent>
</Card>