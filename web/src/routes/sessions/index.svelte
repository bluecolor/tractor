<script>
	import PlayIcon from '@icons/play.svg';
	import FilterIcon from '@icons/filter.svg';
	import TrashIcon from '@icons/trash.svg';
	import InfoIcon from '@icons/info.svg';
	import { onMount } from 'svelte';
	import { api, wsendpoint } from '$lib/utils';
	import Pagination from '@components/Pagination.svelte';

	let page = {};
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
	export let filters = {
		statuses: [],
		extraction: ''
	};

	async function onLoad(params) {
		let url = 'sessions?' + new URLSearchParams(params);
		console.log(url);
		const response = await api('GET', 'sessions');
		if (!response.ok) {
			const result = await response.json();
			alert(result.error);
			return;
		}
		page = await response.json();
		sessions = page.items;
	}
	function onClearFilters() {
		filters = {
			extraction: '',
			statuses: []
		};
		onLoad();
	}
	function onStatusFilter(status) {
		const s = status.toLowerCase();
		if (filters.statuses.includes(s)) {
			filters.statuses = filters.statuses.filter((f) => f !== s);
		} else {
			filters.statuses.push(s);
		}
	}
	function subscribe() {
		const url = wsendpoint('session/feeds');
		const client = new WebSocket(url);
		client.addEventListener('open', () => {
			console.log('Connected to session feed');
		});
		client.addEventListener('message', (event) => {
			const feed = JSON.parse(event.data);
			console.log(feed);
		});
		client.addEventListener('close', () => {
			console.log('Disconnected from session feed');
		});
	}

	onMount(async () => {
		const e = await api('GET', 'extractions');
		extractions = (await e.json()).items || [];
		let params = Object.fromEntries(new URLSearchParams(location.search).entries());
		onLoad(params);
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
				value: s.startedAt && formatDate(s.startedAt)
			},
			{
				label: 'Ended at',
				value: s.finishedAt && formatDate(s.finishedAt)
			},
			{
				label: 'Duration',
				value: `${Math.floor(
					(new Date(s.finishedAt).getTime() - new Date(s.startedAt).getTime()) / 1000
				)} seconds`
			},
			{
				label: 'Read/Write count',
				value: `R:${s.readCount} / W:${s.writeCount}`
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
  .w-full.h-full.flex.flex-col.pt-4.pb-4
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
          select(name='extractions', value='{filters.extraction}')
            +each('extractions as e')
              option(value='{e.id}') {e.name}
        .form-item
          label(for="status") Status
          .flex.items-center.h-5.gap-x-3(name="status")
            +each('["Success", "Pending", "Failed", "Cancelled"] as s')
              .flex
                <input on:input="{() => onStatusFilter(s)}" id="{s + '-status'}" aria-describedby="remember" type="checkbox" class="w-4 h-4 bg-gray-50 rounded border border-gray-300 focus:ring-2 focus:ring-blue-200"/>
                .ml-3.text-sm
                  label(for="{s + '-status'}" class="font-medium text-gray-700") {s}
        .actions.flex.justify-start.gap-x-3
          <button class="btn" on:click="{() => onLoad({page:0, ...filters})}"> Apply </button>
          <button class="btn danger" on:click="{() => onClearFilters()}">Clear </button>


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
      .mt-4
        <Pagination page='{page.page}' total='{page.total}' first='{page.first}' last='{page.last}' maxPage='{page.max_page}' visible='{page.visible}' on:paginate='{(p) => onLoad({page: p.detail, ...filters})}'/>


</template>

<style lang="postcss">
	table thead tr th {
		@apply font-normal text-base text-gray-500 pl-4 pr-4 pb-2 pt-2;
	}
	table tbody tr td {
		@apply font-normal text-base text-gray-700 pl-4 pr-4 pb-2 pt-2;
	}
	.detail-item:hover .label {
		@apply text-gray-700;
	}
	.detail-item:hover .value {
		@apply text-gray-700;
	}
</style>
