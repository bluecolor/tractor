<script>
	import Mappings from './components/Mappings.svelte';
	import MySQLDataset from './components/MySQLDataset.svelte';
	import { onMount } from 'svelte';
	import { api } from '$lib/utils';
	let connections = [];
	let sourceConnectionId = null;
	let targetConnectionId = null;
	let sourceComponent = null;
	let targetComponent = null;
	let sourceDatasetConfig = {};
	let targetDatasetConfig = {};

	function getConnection(id) {
		return connections.find((connection) => connection.id === id);
	}
	function getComponent(connectionId) {
		const components = {
			mysql: MySQLDataset
		};
		let connection = connections.find((c) => c.id == connectionId);
		let connectionTypeCode = connection.connectionType.code;
		const k = Object.keys(components).find(
			(key) => key.toLowerCase() === connectionTypeCode.toLowerCase()
		);
		return components[k];
	}
	onMount(async () => {
		api('GET', 'connections').then(async (response) => {
			if (response.ok) {
				connections = await response.json();
			} else {
				let errm = await response.text();
				alert('Failed to load connections\n' + errm);
			}
		});
	});
</script>

<template lang="pug">
  .w-full.flex.flex-col.pt-4.mb-4
    .flex.justify-between.items-center
      .title
        | New Extraction
    .bg-white.mt-4.rounded-md.flex.flex-col
      .toolbar.flex.justify-end.pt-4.pr-4
        button.btn Save
      .grid.grid-cols-2.gap-4.w-full.p-4
        .source
          .form-item
            label(for='source-connection') Source connection
            select(id='source-connection' name='sourceConnection', bind:value='{sourceConnectionId}')
              +each('connections as c')
                option(value='{c.id}' selected='{sourceConnectionId === c.id}') {c.name}
          .form-item
            +if('sourceConnectionId')
                <svelte:component bind:config={sourceDatasetConfig} this={getComponent(sourceConnectionId)} connection={connections.find((c) => c.id == sourceConnectionId)} />

        .target
          .form-item
            label(for='target-connection') Target connection
            select(id='target-connection' name='targetConnection' bind:value='{targetConnectionId}')
              +each('connections as c')
                option(value='{c.id}' selected='{targetConnectionId === c.id}') {c.name}
          .form-item
            +if('targetConnectionId')
                <svelte:component target={true} bind:config={targetDatasetConfig} this={getComponent(targetConnectionId)} connection={connections.find((c) => c.id == targetConnectionId)} />

      +if('sourceDatasetConfig && targetDatasetConfig && sourceConnectionId && targetConnectionId')
        Mappings(
          sourceDataset='{sourceDatasetConfig}'
          targetDataset='{targetDatasetConfig}'
          sourceConnection='{getConnection(sourceConnectionId)}',
          targetConnection='{getConnection(targetConnectionId)}'
        )
</template>
