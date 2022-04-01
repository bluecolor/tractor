<script>
	import PlayIcon from '@icons/play.svg';
	import MoreIcon from '@icons/more.svg';

	import { onMount } from 'svelte';
	import { api } from '$lib/utils';

	let extractions = [];
	onMount(async () => {
		const response = await api('GET', 'extractions');
		extractions = await response.json();
	});

	function onRunExtraction(id) {
		console.log(id);
		api('POST', `extractions/${id}/run`).then((response) => {
			console.log(response);
			if (response.ok) {
				console.log('Extraction run');
			} else {
				response.text().then((text) => {
					alert('Failed to run extraction\n' + text);
				});
			}
		});
	}
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
            th(scope="col" align="left")
              | Name
            th(scope="col" align="left")
              | Source
            th(scope="col" align="left")
              | Target
            th(scope="col" align="left")
              | Status
            th.actions
        tbody
          +each('extractions as e')
            tr(class="last:border-b-0  hover:bg-gray-50")
              td
                a(href="/extractions/{e.id}")
                  | {e.name}
              td
                | {e.sourceDataset.name}<span class="text-gray-400">@{e.sourceDataset.connection.name} </span>
              td
                | {e.targetDataset.name}<span class="text-gray-400">@{e.targetDataset.connection.name} </span>
              td
                .flex.justify-center.items-center.m-1.font-medium.py-1.px-2.bg-white.rounded-full.text-gray-700.bg-gray-100.border.border-gray-300
                  .text-xs.font-normal.leading-none.max-w-full.flex-initial Idle

              td.actions
                div.flex.justify-end.items-center
                  span(on:click='{onRunExtraction(e.id)}')
                    PlayIcon(class="icon-btn")
                  span(on:click='{onDeleteExtraction(e.id)}')
                    MoreIcon(class="trash")


</template>

<style lang="postcss">
	table thead tr th {
		@apply font-normal text-base text-gray-500 pl-4 pr-4 pb-2 pt-2;
	}
	table tbody tr td {
		@apply font-normal text-base text-gray-700 pl-4 pr-4 pb-2 pt-2;
	}
</style>
