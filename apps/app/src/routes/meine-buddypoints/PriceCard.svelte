<script lang="ts">
	import type { BuddyPackage } from '../../generated';

	type Props = BuddyPackage & {
		badge: string;
		note: string;
		active: boolean;
		onclick: (name: string) => void
	}
	let { name, badge, price_in_cents, points, note, active, onclick, created_at }: Props = $props();
</script>
<button
	onclick={() => onclick(name)}
	class={`flex flex-col gap-2 p-4 border border-solid  rounded-[18px] border-white/12 hover:border-primary/28 text-[#e9e9f1] text-left leading-normal shadow-[0_18px_34px_0_rgb(0_0_0/0.32),0_0_28px_0_rgb(0_255_238/0.08)] transition-[transform,border-color,box-shadow,background] duration-150 cursor-pointer hover:shadow-[0_16px_30px_rgb(0_0_0/0.25)] hover:-translate-y-0.5 ${active && 'border-primary/55 hover:border-primary/55 bg-[linear-gradient(rgba(0,255,238,0.18)_0%,rgba(0,255,238,0.08)_100%)]'} `}
	type="button">
	<span
		class="inline-flex py-1.25 px-2.25 w-fit border border-solid border-white/[0.14] rounded-full bg-white/8 text-[11px] font-extrabold tracking-[0.04em] uppercase"
	>{badge}</span>
	<span class="text-lg font-black">{name}</span>
	<span class="text-primary text-[28px] font-black leading-7">{(price_in_cents / 100).toFixed(2)} EUR</span>
	<span class="text-sm font-extrabold">{points} BuddyPoints</span>
	<span class="text-[#b7b7c8] text-xs leading-[17.4px]">{note}</span>
</button>