<script lang="ts">
	import { page } from '$app/state';
	import { m } from '$lib/paraglide/messages.js';
</script>

<svelte:head>
	<title>{m.datenschutz_title()}</title>
</svelte:head>

<div class="max-w-7xl mx-auto px-5 py-6 pb-16">
	<div class="rounded-2xl bg-white/6 border border-white/10 shadow-2xl p-5 mb-5">
		<h2 id="privacyTitle" class="mb-1.5 text-xl text-white">{m.datenschutz_title()}</h2>
		<p id="privacySubtitle" class="mb-10 text-[#b7b7c8] leading-relaxed">{m.datenschutz_subtitle()}</p>
		<div class="rounded-xl border border-white/10 bg-black/20 p-4 leading-relaxed">
			<ol class="mb-3 ml-5 text-[rgba(233,233,241,0.92)]">
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section1_title()}</h3>
					<p class="mb-2.5">{m.datenschutz_section1_content()}</p>
					<p class="mb-2.5 whitespace-pre-line">{m.datenschutz_section1_contact()}</p>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section2_title()}</h3>
					<p class="mb-2.5 whitespace-pre-line">{m.datenschutz_section2_content()}</p>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section3_title()}</h3>
					<ul class="mb-3 ml-5">
						<li>{m.datenschutz_section3_item1()}</li>
						<li>{m.datenschutz_section3_item2()}</li>
						<li>{m.datenschutz_section3_item3()}</li>
						<li>{m.datenschutz_section3_item4()}</li>
					</ul>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section4_title()}</h3>
					<p class="mb-2.5">{m.datenschutz_section4_intro()}</p>
					<ul class="mb-3 ml-5">
						<li>{m.datenschutz_section4_item1()}</li>
						<li>{m.datenschutz_section4_item2()}</li>
						<li>{m.datenschutz_section4_item3()}</li>
						<li>{m.datenschutz_section4_item4()}</li>
						<li>{m.datenschutz_section4_item5()}</li>
						<li>{m.datenschutz_section4_item6()}</li>
					</ul>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section5_title()}</h3>
					<ul class="mb-3 ml-5">
						<li>{m.datenschutz_section5_item1()}</li>
						<li>{m.datenschutz_section5_item2()}</li>
						<li>{m.datenschutz_section5_item3()}</li>
					</ul>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section6_title()}</h3>
					<p class="mb-2.5">{m.datenschutz_section6_content()}</p>
					<ul class="mb-3 ml-5">
						<li>{m.datenschutz_section6_item1()}</li>
						<li>{m.datenschutz_section6_item2()}</li>
					</ul>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section7_title()}</h3>
					<p class="mb-2.5">{m.datenschutz_section7_intro()}</p>
					<h4 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section7_1_title()}</h4>
					<ul class="mb-3 ml-5">
						<li>{m.datenschutz_section7_1_item1()}</li>
						<li>{m.datenschutz_section7_1_item2()}</li>
						<li>{m.datenschutz_section7_1_item3()}</li>
					</ul>
					<p class="mb-2.5 whitespace-pre-line">{m.datenschutz_section7_1_legal()}</p>
					<ul class="mb-3 ml-5">
						<li>{m.datenschutz_section7_1_legal1()}</li>
						<li>{m.datenschutz_section7_1_legal2()}</li>
						<li>{m.datenschutz_section7_1_legal3()}</li>
						<li>{m.datenschutz_section7_1_legal4()}</li>
					</ul>
					<h4 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section7_2_title()}</h4>
					<ul class="mb-3 ml-5">
						<li>{m.datenschutz_section7_2_item1()}</li>
						<li>{m.datenschutz_section7_2_item2()}</li>
					</ul>
					<p class="mb-2.5 whitespace-pre-line">{m.datenschutz_section7_2_legal()}</p>
					<ul class="mb-3 ml-5">
						<li>{m.datenschutz_section7_2_legal1()}</li>
						<li>{m.datenschutz_section7_2_legal2()}</li>
						<li>{m.datenschutz_section7_2_legal3()}</li>
					</ul>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section8_title()}</h3>
					<p class="mb-2.5 whitespace-pre-line">{m.datenschutz_section8_content()}</p>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section9_title()}</h3>
					<ul class="mb-3 ml-5">
						<li>{m.datenschutz_section9_item1()}</li>
						<li>{m.datenschutz_section9_item2()}</li>
						<li>{m.datenschutz_section9_item3()}</li>
					</ul>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section10_title()}</h3>
					<p class="mb-2.5">{m.datenschutz_section10_intro()}</p>
					<ul class="mb-3 ml-5">
						<li>{m.datenschutz_section10_item1()}</li>
						<li>{m.datenschutz_section10_item2()}</li>
						<li>{m.datenschutz_section10_item3()}</li>
						<li>{m.datenschutz_section10_item4()}</li>
					</ul>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section11_title()}</h3>
					<p class="mb-2.5">{m.datenschutz_section11_intro()}</p>
					<ul class="mb-3 ml-5">
						<li>{m.datenschutz_section11_item1()}</li>
						<li>{m.datenschutz_section11_item2()}</li>
						<li>{m.datenschutz_section11_item3()}</li>
					</ul>
					<p class="mb-2.5">{m.datenschutz_section11_warning()}</p>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section12_title()}</h3>
					<p class="mb-2.5">{m.datenschutz_section12_intro()}</p>
					<ul class="mb-3 ml-5">
						<li>{m.datenschutz_section12_item1()}</li>
						<li>{m.datenschutz_section12_item2()}</li>
						<li>{m.datenschutz_section12_item3()}</li>
					</ul>
					<p class="mb-2.5">{m.datenschutz_section12_legal()}</p>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section13_title()}</h3>
					<p class="mb-2.5">{m.datenschutz_section13_content()}</p>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section14_title()}</h3>
					<p class="mb-2.5">{m.datenschutz_section14_item1()}</p>
					<p class="mb-2.5">{m.datenschutz_section14_item2()}</p>
					<p class="mb-2.5">{m.datenschutz_section14_item3()}</p>
					<p class="mb-2.5">{m.datenschutz_section14_item4()}</p>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section15_title()}</h3>
					<p class="mb-2.5 whitespace-pre-line">{m.datenschutz_section15_content()}</p>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section16_title()}</h3>
					<ul class="mb-3 ml-5">
						<li>{m.datenschutz_section16_item1()}</li>
						<li>{m.datenschutz_section16_item2()}</li>
						<li>{m.datenschutz_section16_item3()}</li>
						<li>{m.datenschutz_section16_item4()}</li>
						<li>{m.datenschutz_section16_item5()}</li>
						<li>{m.datenschutz_section16_item6()}</li>
						<li>{m.datenschutz_section16_item7()}</li>
					</ul>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section17_title()}</h3>
					<p class="mb-2.5">{m.datenschutz_section17_content()}</p>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section18_title()}</h3>
					<p class="mb-2.5">{m.datenschutz_section18_content()}</p>
				</li>
				<li class="mb-2.5">
					<h3 class="mt-3.5 mb-2 text-base text-white">{m.datenschutz_section19_title()}</h3>
					<p class="mb-2.5 whitespace-pre-line">{m.datenschutz_section19_content()}</p>
				</li>
			</ol>
		</div>
	</div>
</div>

<style>
	ol {
		counter-reset: section;
	}
	ol > li {
		counter-increment: section;
	}
	ol > li > h3::before {
		content: counter(section) ". ";
		font-weight: 600;
	}
</style>
