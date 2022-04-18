<script>
	import { page } from '$app/stores'
	import Mappings from './components/Mappings.svelte'
	import Dataset from './components/dataset/Dataset.svelte'
	import { onMount } from 'svelte'
	import { api } from '$lib/utils'
	import Collapse from '@components/Collapse.svelte'
	let name = ''
	let connections = []

	const id = $page.params.id
	let extraction = {
		name: '',
		sourceDataset: {},
		targetDataset: {}
	}

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
		api('GET', `extractions/${id}`).then(async (response) => {
			if (response.ok) {
				extraction = await response.json()
			} else {
				let errm = await response.text()
				alert('Failed to load extraction\n' + errm)
			}
		})
	})
	function onSave() {
		let payload = {
			name,
			sourceDataset,
			targetDataset
		}
		api('POST', 'extractions', payload).then(async (response) => {
			if (response.ok) {
				let extraction = await response.json()
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
        | Extraction {extraction.name}
    .bg-white.mt-4.rounded-md.flex.flex-col
      .toolbar.flex.justify-between.pt-4.px-4.mb-4.items-center
        .name.flex-1.mr-4
          input.input(placeholder="Name", type="text", bind:value='{extraction.name}')
        .actions
          button.btn(on:click='{onSave}') Save
      Collapse(title="Configuration", open='{true}')
        .grid.grid-cols-2.gap-4.w-full.p-4
          .source
            Dataset(type='source' connections='{connections}' dataset='{extraction.sourceDataset}')
          .target
            Dataset(type='target' connections='{connections}' dataset='{extraction.targetDataset}')

      Collapse(title="Mappings")
        .mt-4
          Mappings(extraction='{extraction}')
</template>
