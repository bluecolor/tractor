<script>
	import EditIcon from '@icons/edit.svg'
	import { onMount } from 'svelte'
	import { api } from '$lib/utils'
	import ExtractionMode from './ExtractionMode.svelte'
	export let dataset
	export let type = 'source'

	let databases = []
	let tables = []
	let editTable = false
	dataset.config = dataset.config || {}

	$: {
		if (!dataset.name) {
			dataset.name = `${dataset.config.database}.${dataset.config.table}`
		}
	}

	onMount(async () => {
		api('POST', 'connections/info', {
			connection: dataset.connection,
			info: 'databases'
		}).then(async (response) => {
			if (response.ok) {
				databases = await response.json()
			} else {
				let errm = await response.text()
				alert('Failed to load databases\n' + errm)
			}
		})
		if (dataset?.config?.database) {
			onDatabaseChange(false)
		}
	})
	function toggleEditTable() {
		editTable = !editTable
	}
	function onDatabaseChange(clearTable = true) {
		api('POST', 'connections/info', {
			connection: dataset.connection,
			info: 'tables',
			options: { database: dataset.config.database }
		}).then(async (response) => {
			if (clearTable) {
				dataset.config.table = null
			}
			if (response.ok) {
				tables = await response.json()
			} else {
				tables = []
			}
		})
	}
	function onTableChange() {}
</script>

<template lang="pug">
.form-item
  label(for="database") Database
  select(name='database', bind:value='{dataset.config.database}' on:change='{onDatabaseChange}')
    +each('databases as d')
      option(value='{d}' selected='{d === dataset.config.database}') {d}

.form-item
  label(for="table") Table
  .flex.justify-between.items-center
    +if('!editTable')
      select(name='table', bind:value='{dataset.config.table}' on:change='{onTableChange}')
        +each('tables as t')
          option(value='{t}' selected='{t === dataset.config.table}') {t}
    +if('editTable')
      input.input(type='text', name='table', bind:value='{dataset.config.table}' placeholder='Table name')
    +if('type == "target"')
      .icon-btn(on:click='{toggleEditTable}')
        EditIcon

ExtractionMode(bind:value='{dataset.config.mode}' type='{type}')

</template>
