<script>
	import { api } from '$lib/utils';
	import Dropdown from '@components/Dropdown.svelte';
	import Trash from '@icons/trash.svg';
	import MenuIcon from '@icons/menu.svg';
	import PlusIcon from '@icons/plus.svg';
	import SaveIcon from '@icons/save.svg';

	export let sourceConnection, targetConnection, sourceDataset, targetDataset;

	let mappings = [];

	let options = [
		{
			label: 'Fetch fileds',
			value: 'fetch'
		},
		{
			label: 'Fetch source fields',
			value: 'fetchesource'
		},
		{
			label: 'Fetch target fields',
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
		}
	}
	function mapSourceFields(fields) {
		fields.forEach((field) => {
			let mapping = mappings.find((mapping) => mapping.source.name === field.name);
			if (!mapping) {
				mapping = {
					__index__: mappings.length,
					source: field,
					target: {}
				};
				mappings.push(mapping);
			} else {
				mapping.source = field;
			}
		});
		mappings = [...mappings];
	}
	function onFetchSource() {
		console.log(sourceConnection, sourceDataset);
		api('POST', `connections/${sourceConnection.id}/fields`, sourceDataset).then(
			async (response) => {
				if (response.ok) {
					let sourceFields = await response.json();
					console.log(sourceFields);
					mapSourceFields(sourceFields);
				} else {
					let errm = await response.text();
					alert('Failed to load source fields\n' + errm);
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
            | Target column
          th(scope="col" align="left")
            | Type
          th.actions.flex.justify-end.items-center(align="right")
            SaveIcon.icon-btn.mr-3
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
              input.input(placeholder="Target column", bind:value='{m.target.name}')
            td
              select.cursor-pointer()
                option(value="string") string
                option(value="integer") integer
                option(value="float") float
                option(value="boolean") boolean
                option(value="date") date

            td
              div.flex.justify-end.items-center(on:click='{onDeleteMapping(m)}')
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
