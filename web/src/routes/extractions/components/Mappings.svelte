<script>
	import { api } from '$lib/utils';
	import Dropdown from '@components/Dropdown.svelte';
	import TrashIcon from '@icons/trash.svg';
	import MenuIcon from '@icons/menu.svg';
	import PlusIcon from '@icons/plus.svg';
	import GreaterThanIcon from '@icons/greater-than.svg';

	export let source = { fields: [] };
	export let target = { fields: [] };

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

	$: source.fields = source.fields?.sort((a, b) => a.order - b.order) || [];
	$: target.fields = target.fields?.sort((a, b) => a.order - b.order) || [];

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
	function mapFields({ sf, tf }) {
		source.fields = (sf ?? source.fields).map((s, i) => {
			return {
				order: i,
				...s
			};
		});
		target.fields = (tf ?? target.fields).map((t, i) => {
			return {
				order: i,
				...t
			};
		});

		while (source.fields.length < target.fields.length) {
			source.fields.push({ order: source.fields.length });
		}
		while (target.fields.length < source.fields.length) {
			target.fields.push({ order: target.fields.length });
		}
		source.fields = [...source.fields];
		target.fields = [...target.fields];
		console.log(source.fields);
	}
	function onFetchSource() {
		api('POST', `connections/${source.connectionId}/dataset`, source).then(async (response) => {
			if (response.ok) {
				let dataset = await response.json();
				mapFields({ sf: dataset.fields, tf: undefined });
			} else {
				let errm = await response.text();
				alert('Failed to load source fields\n' + errm);
			}
		});
	}
	function onFetchTarget() {
		api('POST', `connections/${target.connectionId}/dataset`, target).then(async (response) => {
			if (response.ok) {
				let dataset = await response.json();
				mapFields({ sf: undefined, tf: dataset.fields });
			} else {
				let errm = await response.text();
				alert('Failed to load target fields\n' + errm);
			}
		});
	}
	function onFetch() {
		Promise.all([
			api('POST', `connections/${source.connectionId}/dataset`, source.config),
			api('POST', `connections/${target.connectionId}/dataset`, target.config)
		]).then(async (responses) => {
			if (responses.every((r) => r.ok)) {
				let s = await responses[0].json();
				let t = await responses[1].json();
				mapFields({ sf: s.fields, tf: t.fields });
			} else {
				let errms = await Promise.all(responses.map((r) => r.text()));
				alert('Failed to load fields\n' + errms.join('\n'));
			}
		});
	}
	function onDeleteMapping(i) {
		source.fields.splice(i, 1);
		target.fields.splice(i, 1);
		source.fields = [...source.fields];
		target.fields = [...target.fields];
	}
	function onAddMapping() {
		source.fields.push({});
		target.fields.push({});
	}
	function onClear() {
		source.fields = [];
		target.fields = [];
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
        +each('source.fields as sf, i')
          tr(class="last:border-b-0  hover:bg-blue-50")
            td
              input.input(placeholder="Source column", bind:value='{sf.name}')
            td
              span.text-gray-600
                | {sf.type}
            td
              input.input(placeholder="Target column", bind:value='{target.fields[i].name}')
            td
              select.cursor-pointer()
                option(value="string") string
                option(value="integer") integer
                option(value="float") float
                option(value="boolean") boolean
                option(value="date") date
            td
              div.flex.justify-end.items-center
                span.cursor-pointer(on:click='{onDeleteMapping(i)}')
                  GreaterThanIcon(class="fill-current text-gray-200 hover:text-blue-500")
                span(on:click='{onDeleteMapping(i)}')
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
