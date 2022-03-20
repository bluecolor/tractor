<script>
	import ChevronRightIcon from '@icons/chevron-right.svg';
	import ChevronDownIcon from '@icons/chevron-down.svg';
	export let open = true;
	export let title = 'Options';
	function onToggle() {
		open = !open;
	}
</script>

<template lang="pug">

.px-4.cursor-pointer.py-2.flex.justify-between.border-y.border-gray-100(on:click='{onToggle}' class="hover:bg-yellow-50")
  .w-full.text-gray-400.flex.items-center(class="hover:text-gray-500" )
    .chevron.mr-2
      +if('!open')
        ChevronRightIcon(widht="22" height="22" class="fill-current")
      +if('open')
        ChevronDownIcon(widht="22" height="22" class="fill-current")
    .text
      | {title}

div(class:hidden="{!open}")
  slot


</template>
