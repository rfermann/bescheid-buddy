<script>
	import CardHeader from '$lib/components/ui/card/CardHeader.svelte';
	import Card from '$lib/components/ui/card/Card.svelte';
	import CardContent from '$lib/components/ui/card/CardContent.svelte';
	import { m } from '$lib/paraglide/messages';
	import { page } from '$app/state';
</script>

<svelte:head>
	<title>MachKinfach - {m.datenschutz_title()}</title>
</svelte:head>

<Card>
	<CardHeader title={m.datenschutz_title()} />
	<CardContent>
		<div class="space-y-6">
			<p class="text-gray-300">{m.datenschutz_subtitle()}</p>
			<p class="text-gray-300">{m.datenschutz_stand()}</p>

			<h3 class="text-lg font-semibold mb-3 text-secondary">1. {m.datenschutz_section1_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section1_text1()}</p>
			<p class="text-gray-300">{m.datenschutz_section1_text2()}</p>
			<p class="text-gray-300">{m.datenschutz_section1_website()}</p>
			<p class="text-gray-300">{m.datenschutz_section1_email()}</p>

			<h3 class="text-lg font-semibold mb-3 text-secondary">2. {m.datenschutz_section2_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section2_text1()}</p>
			<p class="text-gray-300">{m.datenschutz_section2_text2()}</p>

			<h3 class="text-lg font-semibold mb-3 text-secondary">3. {m.datenschutz_section3_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section3_text1()}</p>
			<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
				<li>{m.datenschutz_section3_list1()}</li>
				<li>{m.datenschutz_section3_list2()}</li>
				<li>{m.datenschutz_section3_list3()}</li>
				<li>{m.datenschutz_section3_list4()}</li>
			</ul>

			<h3 class="text-lg font-semibold mb-3 text-secondary">4. {m.datenschutz_section4_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section4_text1()}</p>
			<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
				<li>{m.datenschutz_section4_list1()}</li>
				<li>{m.datenschutz_section4_list2()}</li>
				<li>{m.datenschutz_section4_list3()}</li>
				<li>{m.datenschutz_section4_list4()}</li>
				<li>{m.datenschutz_section4_list5()}</li>
				<li>{m.datenschutz_section4_list6()}</li>
			</ul>

			<h3 class="text-lg font-semibold mb-3 text-secondary">5. {m.datenschutz_section5_title()}</h3>
			<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
				<li>{m.datenschutz_section5_list1()}</li>
				<li>{m.datenschutz_section5_list2()}</li>
				<li>{m.datenschutz_section5_list3()}</li>
			</ul>

			<h3 class="text-lg font-semibold mb-3 text-secondary">6. {m.datenschutz_section6_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section6_text1()}</p>
			<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
				<li>{m.datenschutz_section6_list1()}</li>
				<li>{m.datenschutz_section6_list2()}</li>
			</ul>

			<h3 class="text-lg font-semibold mb-3 text-secondary">7. {m.datenschutz_section7_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section7_text1()}</p>
			
			<div class="ml-4">
				<h4 class="text-md font-medium mb-2 text-gray-200">7.1 {m.datenschutz_section7_1_title()}</h4>
				<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
					<li>{m.datenschutz_section7_1_list1()}</li>
					<li>{m.datenschutz_section7_1_list2()}</li>
					<li>{m.datenschutz_section7_1_list3()}</li>
					<li>{m.datenschutz_section7_1_list4()}</li>
				</ul>
				<p class="text-gray-300 font-medium mt-2">{m.datenschutz_section7_1_rechtsgrundlagen()}</p>
				<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
					<li>{m.datenschutz_section7_1_list5()}</li>
					<li>{m.datenschutz_section7_1_list6()}</li>
					<li>{m.datenschutz_section7_1_list7()}</li>
					<li>{m.datenschutz_section7_1_list8()}</li>
				</ul>

				<h4 class="text-md font-medium mb-2 mt-4 text-gray-200">7.2 {m.datenschutz_section7_2_title()}</h4>
				<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
					<li>{m.datenschutz_section7_2_list1()}</li>
					<li>{m.datenschutz_section7_2_list2()}</li>
					<li>{m.datenschutz_section7_2_list3()}</li>
				</ul>
				<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
					<li>{m.datenschutz_section7_2_list4()}</li>
					<li>{m.datenschutz_section7_2_list5()}</li>
					<li>{m.datenschutz_section7_2_list6()}</li>
					<li>{m.datenschutz_section7_2_list7()}</li>
				</ul>
			</div>

			<h3 class="text-lg font-semibold mb-3 text-secondary">8. {m.datenschutz_section8_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section8_text1()}</p>
			<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
				<li>{m.datenschutz_section8_list1()}</li>
				<li>{m.datenschutz_section8_list2()}</li>
			</ul>

			<h3 class="text-lg font-semibold mb-3 text-secondary">9. {m.datenschutz_section9_title()}</h3>
			<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
				<li>{m.datenschutz_section9_list1()}</li>
				<li>{m.datenschutz_section9_list2()}</li>
				<li>{m.datenschutz_section9_list3()}</li>
			</ul>

			<h3 class="text-lg font-semibold mb-3 text-secondary">10. {m.datenschutz_section10_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section10_text1()}</p>
			<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
				<li>{m.datenschutz_section10_list1()}</li>
				<li>{m.datenschutz_section10_list2()}</li>
				<li>{m.datenschutz_section10_list3()}</li>
				<li>{m.datenschutz_section10_list4()}</li>
			</ul>

			<h3 class="text-lg font-semibold mb-3 text-secondary">11. {m.datenschutz_section11_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section11_text1()}</p>
			<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
				<li>{m.datenschutz_section11_list1()}</li>
				<li>{m.datenschutz_section11_list2()}</li>
				<li>{m.datenschutz_section11_list3()}</li>
			</ul>
			<p class="text-gray-300">{m.datenschutz_section11_text2()}</p>

			<h3 class="text-lg font-semibold mb-3 text-secondary">12. {m.datenschutz_section12_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section12_text1()}</p>
			<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
				<li>{m.datenschutz_section12_list1()}</li>
				<li>{m.datenschutz_section12_list2()}</li>
				<li>{m.datenschutz_section12_list3()}</li>
			</ul>
			<p class="text-gray-300">{m.datenschutz_section12_text2()}</p>

			<h3 class="text-lg font-semibold mb-3 text-secondary">13. {m.datenschutz_section13_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section13_text1()}</p>

			<h3 class="text-lg font-semibold mb-3 text-secondary">14. {m.datenschutz_section14_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section14_text1()}</p>
			<p class="text-gray-300">{m.datenschutz_section14_text2()}</p>
			<p class="text-gray-300">{m.datenschutz_section14_text3()}</p>
			<p class="text-gray-300">{m.datenschutz_section14_text4()}</p>

			<h3 class="text-lg font-semibold mb-3 text-secondary">15. {m.datenschutz_section15_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section15_text1()}</p>
			<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
				<li>{m.datenschutz_section15_list1()}</li>
				<li>{m.datenschutz_section15_list2()}</li>
				<li>{m.datenschutz_section15_list3()}</li>
			</ul>

			<h3 class="text-lg font-semibold mb-3 text-secondary">16. {m.datenschutz_section16_title()}</h3>
			<ul class="list-disc list-inside text-gray-300 space-y-1 ml-4">
				<li>{m.datenschutz_section16_list1()}</li>
				<li>{m.datenschutz_section16_list2()}</li>
				<li>{m.datenschutz_section16_list3()}</li>
				<li>{m.datenschutz_section16_list4()}</li>
				<li>{m.datenschutz_section16_list5()}</li>
				<li>{m.datenschutz_section16_list6()}</li>
				<li>{m.datenschutz_section16_list7()}</li>
			</ul>

			<h3 class="text-lg font-semibold mb-3 text-secondary">17. {m.datenschutz_section17_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section17_text1()}</p>

			<h3 class="text-lg font-semibold mb-3 text-secondary">18. {m.datenschutz_section18_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section18_text1()}</p>

			<h3 class="text-lg font-semibold mb-3 text-secondary">19. {m.datenschutz_section19_title()}</h3>
			<p class="text-gray-300">{m.datenschutz_section19_text1()}</p>
		</div>
	</CardContent>
</Card>
