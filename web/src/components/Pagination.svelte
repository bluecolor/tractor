<script>
	import { createEventDispatcher } from 'svelte'
	const dispatch = createEventDispatcher()
	export let detail = false
	export let page = 0
	export let total = 0
	export let first = true
	export let last = false
	export let visible = 0
	export let maxPage = 0

	function onPrev() {
		if (page > 0) {
			dispatch('paginate', page - 1)
		}
	}
	function onNext() {
		if (page < maxPage) {
			dispatch('paginate', page + 1)
		}
	}
</script>

<template lang="pug">
.flex.justify-end.items-center
  +if('detail')
    .inline-flex.mr-2
      span.text-gray-300
        | page {page} showing {visible} total {total}
  .inline-flex.pagination(class='xs:mt-0')
    button.prev(disabled='{first}' on:click='{onPrev}' class="{first ? 'disabled' : ''}")
      | Prev
    button.next(disabled='{last}' on:click='{onNext}' class="{last ? 'disabled' : ''}")
      | Next
</template>

<style lang="postcss">
	.pagination button {
		@apply focus:ring-blue-200
      focus:ring-offset-blue-200
      focus:ring-2
      focus:ring-offset-1
      focus:outline-none
      focus:border-blue-200;
	}

	.pagination button.disabled {
		cursor: not-allowed;
		pointer-events: none;
		opacity: 0.5;
	}
	.pagination button.prev {
		@apply py-2
      px-4
      text-sm
      font-medium
      text-gray-800
      border
      border-gray-300
      rounded-l-sm;
	}
	.pagination button.next {
		@apply py-2
      px-4
      text-sm
      font-medium
      text-gray-800
      border
      border-l-0
      border-gray-300
      rounded-r-sm;
	}
</style>
