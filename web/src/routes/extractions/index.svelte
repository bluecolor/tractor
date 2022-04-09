<script>
	import PlayIcon from '@icons/play.svg';
	import MoreIcon from '@icons/more.svg';
	import Dropdown from '@components/Dropdown.svelte';
	import { onMount } from 'svelte';
	import { api } from '$lib/utils';
	import { session } from '$app/stores';

	let extractions = [];
	let options = [
		{
			label: 'Delete',
			value: 'delete'
		},
		{
			label: 'Sessions',
			value: 'sessions'
		}
	];

	onMount(async () => {
		const response = await api('GET', 'extractions');
		const result = await response.json();
		extractions = result.map((r) => {
			if (r.sessions.length > 0) {
				r.status = r.sessions[0].status;
			} else {
				r.status = null;
			}
			return r;
		});
	});

	function onRunExtraction(id) {
		api('POST', `extractions/${id}/run`).then(async (response) => {
			if (response.ok) {
				console.log('Extraction run');
				const session = await response.json();
				extractions = extractions.map((e) => {
					if (e.id === session.extraction.id) {
						e.status = session.status;
					}
					return e;
				});
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
	function onDropdown(e, id) {
		switch (e.detail.value) {
			case 'delete':
				onDeleteExtraction(id);
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
              td(scope="col" align="left")
                +if('e.status === "success"')
                  .flex.justify-center.items-center.font-medium.py-1.px-2.bg-white.rounded-full.text-green-700.bg-green-100.border.border-green-300
                    .text-xs.font-normal.leading-none.max-w-full.flex-initial { e.status }
                  +elseif('e.status != null')
                    .flex.justify-center.items-center.font-medium.py-1.px-2.bg-white.rounded-full.text-gray-700.bg-gray-100.border.border-gray-300
                      .text-xs.font-normal.leading-none.max-w-full.flex-initial { e.status }

              td.actions
                div.flex.justify-end.items-center
                  span(on:click='{onRunExtraction(e.id)}')
                    PlayIcon(class="icon-btn")
                  <Dropdown label="Options" bind:options='{options}' on:select='{(x) => onDropdown(x, e.id)}'>
                    div(slot="button")
                      MoreIcon.icon-btn()
                  </Dropdown>


</template>

<style lang="postcss">
	table thead tr th {
		@apply font-normal text-base text-gray-500 pl-4 pr-4 pb-2 pt-2;
	}
	table tbody tr td {
		@apply font-normal text-base text-gray-700 pl-4 pr-4 pb-2 pt-2;
	}
</style>
