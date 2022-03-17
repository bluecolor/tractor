<script>
	import { api } from '$lib/utils';
	import Dropdown from '@components/Dropdown.svelte';
	import Trash from '@icons/trash.svg';
	import MenuIcon from '@icons/menu.svg';
	import PlusIcon from '@icons/plus.svg';

	export let sourceConnection, targetConnection, sourceDataset, targetDataset;

	let mappings = [];

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
		},
		{
			label: 'Config',
			value: 'config'
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
				__index__: i,
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
	function onDeleteMapping({ __index__ }) {
		mappings.splice(__index__, 1);
		mappings = [...mappings];
	}
	function onAddMapping() {
		mappings = [
			...mappings,
			{
				__index__: mappings.length,
				source: {},
				target: {}
			}
		];
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
            | Target Type
          th.actions.flex.justify-end.items-center(align="right")
            .action.icon-btn.mr-3(on:click='{onAddMapping}')
              PlusIcon()
            Dropdown(label="Options" bind:options='{options}' on:select='{onDropdown}') Reset
              div(slot="button")
                MenuIcon.icon-btn()
      tbody
        +each('mappings as m')
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
                span(on:click='{onDeleteMapping(m)}')
                  Trash(class="trash")
</template>

<style lang="postcss">
	table thead tr th {
		@apply font-normal text-base text-gray-500 pl-4 pr-4 pb-2 pt-2;
	}
	table tbody tr td {
		@apply font-normal text-base text-gray-700 pl-4 pr-4 pb-2 pt-2;
	}
</style>
