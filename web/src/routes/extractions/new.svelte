<script>
	import Mappings from './components/Mappings.svelte'
	import Dataset from './components/dataset/Dataset.svelte'
	import { onMount } from 'svelte'
	import { api } from '$lib/utils'
	import Collapse from '@components/Collapse.svelte'
	let name = ''
	let connections = []
	let extraction = {}
	onMount(async () => {
		api('GET', 'connections?size=-1').then(async (response) => {
			if (response.ok) {
				let result = await response.json()
				connections = result.items
			} else {
				let errm = await response.text()
				alert('Failed to load connections\n' + errm)
			}
		})
	})
	function onSave() {
		api('POST', 'extractions', extraction).then(async (response) => {
			if (response.ok) {
				let extraction = await response.json()
				console.log(extraction)
			} else {
				let errm = await response.text()
				alert('Failed to save extraction\n' + errm)
			}
		})
	}
</script>

<template lang="pug">
  .w-full.flex.flex-col.pt-4.mb-4
    .flex.justify-between.items-center
      .title
        | New Extraction {extraction.name}
    .bg-white.mt-4.rounded-md.flex.flex-col
      .toolbar.flex.justify-between.pt-4.px-4.mb-4.items-center
        .name.flex-1.mr-4
          input.input(placeholder="Name", type="text", bind:value='{name}')
        .actions
          button.btn(on:click='{onSave}') Save
      Collapse(title="Configuration", open='{true}')
        .grid.grid-cols-2.gap-4.w-full.p-4
          .source
            Dataset(type='source' connections='{connections}' bind:dataset='{extraction.sourceDataset}')
          .target
            Dataset(type='target' connections='{connections}' bind:dataset='{extraction.targetDataset}')

      Collapse(title="Mappings")
        .mt-4
          Mappings(bind:extraction='{extraction}')
</template>
