<script>
	import PlayIcon from '@icons/play.svg';
	import FilterIcon from '@icons/filter.svg';
	import TrashIcon from '@icons/trash.svg';
	import InfoIcon from '@icons/info.svg';
	import { onMount } from 'svelte';
	import { api, wsendpoint } from '$lib/utils';

	let extraction;
	let sessions = [];
	let extractions = [];
	let filtersOpen = false;
	let detailId = null;
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

	function subscribe() {
		const url = wsendpoint('session/feeds');
		const client = new WebSocket(url);
		client.addEventListener('open', () => {
			console.log('Connected to session feed');
		});
		client.addEventListener('message', (event) => {
			const data = JSON.parse(event.data);
			console.log(data);
		});
		client.addEventListener('close', () => {
			console.log('Disconnected from session feed');
		});
	}

	onMount(async () => {
		Promise.all([api('GET', 'extractions'), api('GET', 'sessions')]).then(async ([e, s]) => {
			extractions = await e.json();
			sessions = await s.json();
			console.log(sessions);
		});

		let response = await api('GET', 'sessions');
		sessions = await response.json();
		response = await api('GET', 'extractions');
		extractions = await response.json();
		subscribe();
	});

	function onRunExtraction(id) {
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
	function onDeleteSession(id) {
		let ok = confirm('Are you sure you want to delete this session? ');
		if (ok) {
			api('DELETE', 'sessions/' + id).then((response) => {
				if (response.ok) {
					sessions = sessions.filter((s) => s.id !== id);
				} else {
					response.text().then((text) => {
						alert('Failed to delete session\n' + text);
					});
				}
			});
		}
	}
	function formatDate(d) {
		return new Date(d).toLocaleString();
	}
	function getSessionDetails(s) {
		return [
			{
				label: 'Started at',
				value: formatDate(s.createdAt)
			},
			{
				label: 'Ended at',
				value: formatDate(s.createdAt)
			},
			{
				label: 'Duration',
				value: '100 minutes'
			},
			{
				label: 'Read/Write count',
				value: 'R:1012 / W:1012'
			}
		];
	}
	$: getRowClass = function (sessionId) {
		if (detailId === sessionId) {
			return 'last:border-b-0 bg-yellow-50';
		}
		return 'last:border-b-0 hover:bg-gray-50';
	};
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
          +each('sessions as s')
            tr(class="{getRowClass(s.id)}")
              td
                a(href="/extractions/{s.extraction.id}")
                  | {s.extraction.name}
              td
                | {s.extraction.sourceDataset.name}<span class="text-gray-400">@{s.extraction.sourceDataset.connection.name} </span>
              td
                | {s.extraction.targetDataset.name}<span class="text-gray-400">@{s.extraction.targetDataset.connection.name} </span>
              td
                +if('s.status === "success"')
                  .flex.justify-center.items-center.font-medium.py-1.px-2.bg-white.rounded-full.text-green-700.bg-green-100.border.border-green-300
                    .text-xs.font-normal.leading-none.max-w-full.flex-initial { s.status }
                  +else
                    .flex.justify-center.items-center.font-medium.py-1.px-2.bg-white.rounded-full.text-gray-700.bg-gray-100.border.border-gray-300
                      .text-xs.font-normal.leading-none.max-w-full.flex-initial { s.status }
              td.actions
                div.flex.justify-end.items-center
                  span(on:click='{onDeleteSession(s.id)}')
                    TrashIcon(class="trash")
                  <span on:click='{() => detailId = detailId && detailId===s.id ? null : s.id }'>
                    InfoIcon(class="icon-btn")
                  </span>
            +if('detailId === s.id')
              <tr class="{getRowClass(s.id)}">
                td(colspan="5")
                  .flex.flex-col.pl-4
                    +each("getSessionDetails(s) as d")
                      .detail-item.flex.border-b.border-gray-200(class="last:border-b-0")
                        .label.text-gray-400.border-r.border-gray-200(class="w-1/3") {d.label}
                        .value.text-gray-500.pl-2 {d.value}

              </tr>

</template>

<style lang="postcss">
	table thead tr th {
		@apply font-normal text-base text-gray-500 pl-4 pr-4 pb-2 pt-2;
	}
	table tbody tr td {
		@apply font-normal text-base text-gray-700 pl-4 pr-4 pb-2 pt-2;
	}
	.detail-item:hover .label {
		@apply text-red-400;
	}
	.detail-item:hover .value {
		@apply text-red-500;
	}
</style>
