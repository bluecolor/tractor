<script>
	import EditIcon from '@icons/edit.svg';
	import { onMount } from 'svelte';
	import { api } from '$lib/utils';
	export let connection, config, target;

	let databases = [];
	let tables = [];
	let editTable = false;
	export let modes = [
		{
			label: 'Append',
			value: 'append'
		},
		{
			label: 'Replace',
			value: 'replace'
		},
		{
			label: 'Merge',
			value: 'merge'
		}
	];
	let mode = undefined;
	onMount(async () => {
		api('POST', 'connections/info', {
			connection: connection,
			info: 'databases'
		}).then(async (response) => {
			if (response.ok) {
				databases = await response.json();
			} else {
				let errm = await response.text();
				alert('Failed to load databases\n' + errm);
			}
		});
	});
	function toggleEditTable() {
		editTable = !editTable;
	}
	function onDatabaseChange() {
		api('POST', 'connections/info', {
			connection: connection,
			info: 'tables',
			options: { database: config.database }
		}).then(async (response) => {
			if (response.ok) {
				tables = await response.json();
				config.table = null;
			} else {
				tables = [];
				config.table = null;
			}
		});
	}
	function onTableChange() {}

	function onModeSelect(m) {
		mode = m.value;
	}
</script>

<template lang="pug">
.form-item
  label(for="database") Database
  select(name='database', bind:value='{config.database}' on:change='{onDatabaseChange}')
    +each('databases as d')
      option(value='{d}' selected='{d === config.database}') {d}

.form-item
  label(for="table") Table
  .flex.justify-between.items-center
    +if('!editTable')
      select(name='table', bind:value='{config.table}' on:change='{onTableChange}')
        +each('tables as t')
          option(value='{t}' selected='{t === config.table}') {t}
    +if('editTable')
      input.input(type='text', name='table', bind:value='{config.table}' placeholder='Table name')
    +if('target')
      .icon-btn(on:click='{toggleEditTable}')
        EditIcon
+if('target')
  .form-item
    label(for="table") Mode
    .flex
      +each('modes as m')
        div(on:click='{onModeSelect(m)}')
          label.inline-flex.items-center.mr-2()
            input.form-radio.h-4.w-4.text-blue-400(type='radio' checked='{mode==m.value}'  class='focus:ring-2 focus:ring-offset-1 focus:border-blue-300 focus:ring focus:ring-offset-0 focus:ring-blue-200')
            span.ml-2.text-gray-400 {m.label}


</template>
