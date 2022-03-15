<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/utils';
	export let connection, config;
	let databases = [];
	let tables = [];

	function onDatabaseChange() {
		api('POST', 'connections/connectors/resolve', {
			connection: connection,
			request: 'tables',
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
	onMount(async () => {
		api('POST', 'connections/connectors/resolve', {
			connection: connection,
			request: 'databases'
		}).then(async (response) => {
			if (response.ok) {
				databases = await response.json();
			} else {
				let errm = await response.text();
				alert('Failed to load databases\n' + errm);
			}
		});
	});
</script>

<template lang="pug">
.form-item
  label(for="database") Database
  select(name='database', bind:value='{config.database}' on:change='{onDatabaseChange}')
    +each('databases as d')
      option(value='{d}' selected='{d === config.database}') {d}

.form-item
  label(for="table") Table
  select(name='table', bind:value='{config.table}' on:change='{onTableChange}')
    +each('tables as t')
      option(value='{t}' selected='{t === config.table}') {t}


</template>
