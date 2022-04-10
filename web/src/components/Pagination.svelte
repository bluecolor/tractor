<script>
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();
	export let page = 0;
	export let total = 0;
	export let first = true;
	export let last = false;
	export let visible = 0;
	export let maxPage = 0;

	function onPrev() {
		if (page > 0) {
			dispatch('paginate', page - 1);
		}
	}
	function onNext() {
		if (page < maxPage) {
			dispatch('paginate', page + 1);
		}
	}
</script>

<template lang="pug">
.flex.justify-end.items-center
  .inline-flex.mr-2
    span.text-gray-300
      | page {page} showing {visible} total {total}
  .inline-flex.pagination(class='xs:mt-0')
    button.prev(disabled='{first}' on:click='{onPrev}')
      | Prev
    button.next(disabled='{last}' on:click='{onNext}')
      | Next
</template>
