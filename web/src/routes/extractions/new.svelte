<script>
	import Mappings from './components/Mappings.svelte';
	import MySQLDataset from './components/MySQLDataset.svelte';
	import { onMount } from 'svelte';
	import { api } from '$lib/utils';
	import Collapse from '@components/Collapse.svelte';
	let name = '';
	let connections = [];
	let sourceConnectionId = null;
	let targetConnectionId = null;
	let sourceComponent = null;
	let targetComponent = null;
	let sourceDatasetConfig = {
		mode: 'full'
	};
	let targetDatasetConfig = {
		mode: 'append'
	};
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
	function onSave() {
		let payload = {
			name,
			sourceConnectionId,
			targetConnectionId
		};
	}
</script>

<template lang="pug">
  .w-full.flex.flex-col.pt-4.mb-4
    .flex.justify-between.items-center
      .title
        | New Extraction {name}
    .bg-white.mt-4.rounded-md.flex.flex-col
      .toolbar.flex.justify-between.pt-4.px-4.mb-4.items-center
        .name.flex-1.mr-4
          input.input(placeholder="Name", type="text", bind:value='{name}')
        .actions
          button.btn(on:click='{onSave}') Save
      Collapse(title="Configuration", open='{true}')
        .grid.grid-cols-2.gap-4.w-full.p-4
          .source
            .form-item
              label(for='source-connection') Source connection
              select(id='source-connection' name='sourceConnection', bind:value='{sourceConnectionId}')
                +each('connections as c')
                  option(value='{c.id}' selected='{sourceConnectionId === c.id}') {c.name}
            .form-item
              +if('sourceConnectionId')
                  <svelte:component usedIn='source' bind:config={sourceDatasetConfig} this={getComponent(sourceConnectionId)} connection={connections.find((c) => c.id == sourceConnectionId)} />

          .target
            .form-item
              label(for='target-connection') Target connection
              select(id='target-connection' name='targetConnection' bind:value='{targetConnectionId}')
                +each('connections as c')
                  option(value='{c.id}' selected='{targetConnectionId === c.id}') {c.name}
            .form-item
              +if('targetConnectionId')
                  <svelte:component usedIn='target' bind:config={targetDatasetConfig} this={getComponent(targetConnectionId)} connection={connections.find((c) => c.id == targetConnectionId)} />

      +if('sourceDatasetConfig && targetDatasetConfig && sourceConnectionId && targetConnectionId')
        Collapse(title="Mappings")
          .mt-4
            Mappings(
              sourceDataset='{sourceDatasetConfig}'
              targetDataset='{targetDatasetConfig}'
              sourceConnection='{getConnection(sourceConnectionId)}',
              targetConnection='{getConnection(targetConnectionId)}'
            )
</template>
