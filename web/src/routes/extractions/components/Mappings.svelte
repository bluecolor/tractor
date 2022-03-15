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
					source: field,
					target: null
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
	function onDeleteMapping(m) {
		mappings = mappings.filter(
			(mapping) =>
				(!mapping.source || mapping.source?.name !== m.source?.name) &&
				(!mapping.target || mapping.target?.name !== m.target?.name)
		);
	}
</script>

<template lang="pug">
.mappings
  .mappings
    table.table.min-w-full
      thead
        tr
          th(scope="col")
            | Source column
          th(scope="col")
            | Target column
          th.actions.flex.space-x-2.justify-end(align="right")
            PlusIcon.icon-btn
            Dropdown(label="Options" bind:options='{options}' on:select='{onDropdown}') Reset
              div(slot="button")
                MenuIcon.icon-btn()
      tbody
        +each('mappings as m')
          tr(class="last:border-b-0  hover:bg-gray-50")
            td
              | {m.source?.name}
            td
              | {m.target?.name}
            td.actions
              span(on:click='{onDeleteMapping(m)}')
                Trash(class="trash")
</template>

<style lang="postcss">
	table thead tr th {
		@apply font-normal text-base text-gray-500;
	}
</style>
