<script>
	import { createEventDispatcher } from 'svelte'
	import { clickOutside } from '$lib/utils'

	const dispatch = createEventDispatcher()

	export let label = 'Options'
	export let options = []

	let isOpen = false
	function toggle() {
		isOpen = !isOpen
	}
	function onSelect(item) {
		isOpen = false
		dispatch('select', item)
	}
</script>

<template lang="pug">
<div class="relative inline-block text-left" use:clickOutside={() => { isOpen = false;}}>
  div(on:click="{toggle}")
    slot(name="button")
      button.dd-btn(
        type='button'
        class='hover:bg-gray-50 focus:outline-none focus:ring-1 focus:ring-offset-1 focus:ring-offset-blue-200 focus:ring-blue-200'
      )
        | {label}
        svg(width='20' height='20' fill='currentColor' viewBox='0 0 1792 1792' xmlns='http://www.w3.org/2000/svg')
          path(d='M1408 704q0 26-19 45l-448 448q-19 19-45 19t-45-19l-448-448q-19-19-19-45t19-45 45-19h896q26 0 45 19t19 45z')
  +if('isOpen && options.length')
    .origin-top-right.absolute.right-0.mt-2.w-56.rounded-md.shadow-lg.bg-white.ring-1.ring-black.ring-opacity-5(class='')
      .py-1.divide-y.divide-gray-100(role='menu' aria-orientation='vertical' aria-labelledby='options-menu')
        +each('options as o')
          div.cursor-pointer.flex.items-center.block.px-4.py-2.text-md.text-gray-700(
            on:click="{onSelect(o)}"
            class='hover:bg-gray-100 hover:text-gray-900' role='menuitem'
          )
            +if('o.icon')
              .icon.mr-4
                <svelte:component this={o.icon}/>
            .flex.flex-1
              | {o.label}
</div>
</template>

<style lang="postcss">
	.dd-btn {
		@apply border
	     border-gray-300
	     bg-white
	     flex
	     items-center
	     justify-center
	     w-full
	     rounded-md
	     px-4
	     py-1.5
	     text-sm
	     font-medium
	     text-gray-700;
	}
</style>
