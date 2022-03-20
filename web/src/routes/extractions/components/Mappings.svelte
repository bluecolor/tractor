<script>
	import { api } from '$lib/utils';
	import Dropdown from '@components/Dropdown.svelte';
	import TrashIcon from '@icons/trash.svg';
	import MenuIcon from '@icons/menu.svg';
	import PlusIcon from '@icons/plus.svg';
	import GreaterThanIcon from '@icons/greater-than.svg';

	export let sourceConnection, targetConnection, sourceDataset, targetDataset;

	export let mappings = [];

	let options = [
		{
			label: 'Fetch',
			value: 'fetch'
		},
		{
			label: 'Fetch source',
			value: 'fetchesource'
		},
		{
			label: 'Fetch target',
			value: 'fetchetarget'
		},
		{
			label: 'Clear',
			value: 'clear'
		}
	];
	function onDropdown(e) {
		const { value } = e.detail;
		switch (value) {
			case 'fetchesource':
				onFetchSource();
				break;
			case 'fetchetarget':
				onFetchTarget();
				break;
			case 'fetch':
				onFetch();
				break;
			case 'clear':
				onClear();
				break;
		}
	}
	function mapFields({ sources, targets }) {
		sources = sources ?? mappings.map((m) => m.source);
		targets = targets ?? mappings.map((m) => m.target);
		while (sources.length < targets.length) {
			sources.push({});
		}
		while (targets.length < sources.length) {
			targets.push({});
		}
		mappings = sources.map((source, i) => {
			let m = {
				source: source
			};
			let ti = targets.findIndex((t) => t.name === source.name);
			if (ti > -1) {
				m.target = targets[ti];
				targets.splice(ti, 1);
				return m;
			}
			m.target = targets.shift();
			return m;
		});
		mappings = [...mappings];
	}
	function onFetchSource() {
		api('POST', `connections/${sourceConnection.id}/dataset`, sourceDataset).then(
			async (response) => {
				if (response.ok) {
					let source = await response.json();
					mapFields({ sources: source.fields, targets: undefined });
				} else {
					let errm = await response.text();
					alert('Failed to load source fields\n' + errm);
				}
			}
		);
	}
	function onFetchTarget() {
		api('POST', `connections/${targetConnection.id}/dataset`, targetDataset).then(
			async (response) => {
				if (response.ok) {
					let target = await response.json();
					mapFields({ sources: undefined, targets: target.fields });
				} else {
					let errm = await response.text();
					alert('Failed to load target fields\n' + errm);
				}
			}
		);
	}
	function onFetch() {
		Promise.all([
			api('POST', `connections/${sourceConnection.id}/dataset`, sourceDataset),
			api('POST', `connections/${targetConnection.id}/dataset`, targetDataset)
		]).then(async (responses) => {
			if (responses.every((r) => r.ok)) {
				let sources = await responses[0].json();
				let targets = await responses[1].json();
				mapFields({ sources: sources.fields, targets: targets.fields });
			} else {
				let errms = await Promise.all(responses.map((r) => r.text()));
				alert('Failed to load fields\n' + errms.join('\n'));
			}
		});
	}
	function onDeleteMapping(i) {
		mappings.splice(i, 1);
		mappings = [...mappings];
	}
	function onAddMapping() {
		mappings = [
			...mappings,
			{
				source: {},
				target: {}
			}
		];
	}
	function onClear() {
		mappings = [];
	}
</script>

<template lang="pug">
.mappings
  .mappings
    table.table.min-w-full
      thead
        tr
          th(scope="col" align="left")
            | Source column
          th(scope="col" align="left")
            | Source type
          th(scope="col" align="left")
            | Target column
          th(scope="col" align="left")
            | Target type
          th.actions.flex.justify-end.items-center(align="right")
            .action.icon-btn.mr-3(on:click='{onAddMapping}')
              PlusIcon()
            Dropdown(label="Options" bind:options='{options}' on:select='{onDropdown}') Reset
              div(slot="button")
                MenuIcon.icon-btn()
      tbody
        +each('mappings as m, i')
          tr(class="last:border-b-0  hover:bg-blue-50")
            td
              input.input(placeholder="Source column", bind:value='{m.source.name}')
            td
              span.text-gray-600
                | {m.source.type}
            td
              input.input(placeholder="Target column", bind:value='{m.target.name}')
            td
              select.cursor-pointer()
                option(value="string") string
                option(value="integer") integer
                option(value="float") float
                option(value="boolean") boolean
                option(value="date") date
            td
              div.flex.justify-end.items-center
                span.cursor-pointer(on:click='{onDeleteMapping(m, i)}')
                  GreaterThanIcon(class="fill-current text-gray-200 hover:text-blue-500")
                span(on:click='{onDeleteMapping(m, i)}')
                  TrashIcon(class="trash")

</template>

<style lang="postcss">
	table thead tr th {
		@apply font-normal text-base text-gray-500 pl-4 pr-4 pb-2 pt-2;
	}
	table tbody tr td {
		@apply font-normal text-base text-gray-700 pl-4 pr-4 pb-2 pt-2;
	}
</style>
