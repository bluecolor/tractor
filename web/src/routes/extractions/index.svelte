<script>
	import Trash from '../../assets/icons/trash.svg';

	import { onMount } from 'svelte';
	import { api } from '$lib/utils';

	let extractions = [];
	onMount(async () => {
		const response = await api('GET', 'extractions');
		extractions = await response.json();
	});

	function onDeleteExtraction(id) {
		let extraction = extractions.find((e) => e.id === id);
		let ok = confirm('Are you sure you want to delete this extraction? ' + extraction.name);
		if (ok) {
			api('DELETE', 'extractions/' + id).then((response) => {
				if (response.ok) {
					extractions = extractions.filter((e) => e.id !== id);
				} else {
					response.text().then((text) => {
						alert('Failed to delete extraction\n' + text);
					});
				}
			});
		}
	}
</script>

<template lang="pug">
  .w-full.h-full.flex.flex-col.pt-4
    .flex.justify-between.items-center
      .title
        | Extractions
      .search.space-x-2.inline-flex.items-center()
        .action
          input.input(type="text" placeholder="Search")
        a.action(href="/extractions/new")
          button.btn Add

    .bg-white.mt-4.p-2.rounded-md
      table.min-w-full
        thead
          tr
            th(scope="col")
              | Name
            th(scope="col")
              | Type
            th.actions
        tbody
          +each('extractions as e')
            tr(class="last:border-b-0  hover:bg-gray-50")
              td
                a(href="/extractions/{e.id}")
                  | {e.name}
              td
                | {e.name}
              td.actions
                span(on:click='{onDeleteExtraction(e.id)}')
                  Trash(class="trash")

</template>
