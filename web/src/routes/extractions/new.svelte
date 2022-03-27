<script>
	import Mappings from './components/Mappings.svelte';
	import Dataset from './components/dataset/Dataset.svelte';
	import { onMount } from 'svelte';
	import { api } from '$lib/utils';
	import Collapse from '@components/Collapse.svelte';
	let name = '';
	let connections = [];
	let sourceDataset = { config: {} };
	let targetDataset = { config: {} };
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
	function onSave() {
		let payload = {};
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
            Dataset(type='source' connections='{connections}' dataset='{sourceDataset}')
          .target
            Dataset(type='target' connections='{connections}' dataset='{targetDataset}')

      Collapse(title="Mappings")
        .mt-4
          Mappings(source='{sourceDataset}' target='{targetDataset}')
</template>
