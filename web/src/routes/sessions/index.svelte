<script>
	import PlayIcon from '@icons/play.svg';
	import MoreIcon from '@icons/more.svg';
	import FilterIcon from '@icons/filter.svg';
	import Dropdown from '@components/Dropdown.svelte';
	import { onMount } from 'svelte';
	import { api } from '$lib/utils';

	let extraction;
	let extractions = [];
	let filtersOpen = false;
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
				r.status = 'pending';
			}
			return r;
		});
	});

	function onRunExtraction(id) {
		console.log(id);
		api('POST', `extractions/${id}/run`).then((response) => {
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
        | Sessions
      .search.space-x-2.inline-flex.items-center()
        <span class="action" on:click='{() => filtersOpen = !filtersOpen}'>
          FilterIcon(class="icon-btn")
        </span>
    +if('filtersOpen')
      .bg-white.mt-4.p-2.rounded-md
        .form-item
          label(for="extractions") Extraction
          select(name='extractions', bind:value='{extraction}')
            +each('extractions as e')
              option(value='{e}') {e.name}
        .form-item
          label(for="status") Status
          .flex.items-center.h-5.gap-x-3(name="status")
            +each('["Success", "Pending", "Failed", "Cancelled"] as s')
              .flex
                input(id="{s + '-status'}" aria-describedby="remember" type="checkbox" class="w-4 h-4 bg-gray-50 rounded border border-gray-300 focus:ring-2 focus:ring-blue-200")
                .ml-3.text-sm
                  label(for="{s + '-status'}" class="font-medium text-gray-700") {s}
        .actions.flex.justify-start.gap-x-3
          button.btn Apply
          button.btn.danger Clear

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
                +if('e.status === "success"')
                  .flex.justify-center.items-center.m-1.font-medium.py-1.px-2.bg-white.rounded-full.text-green-700.bg-green-100.border.border-green-300
                    .text-xs.font-normal.leading-none.max-w-full.flex-initial { e.status }
                  +else
                    .flex.justify-center.items-center.m-1.font-medium.py-1.px-2.bg-white.rounded-full.text-gray-700.bg-gray-100.border.border-gray-300
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
