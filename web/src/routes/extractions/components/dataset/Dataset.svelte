<script>
	import MySQLDataset from './MySQLDataset.svelte'
	import FileDataset from './FileDataset.svelte'

	export let type = 'source'
	export let connections = []
	export let dataset = {}
	const components = {
		mysql: MySQLDataset,
		file: FileDataset
	}
	$: component = () => {
		let connection = connections.find((c) => c.id == dataset.connectionId)
		const k = Object.keys(components).find(
			(key) => key.toLowerCase() === connection.connectionType.code.toLowerCase()
		)
		return components[k]
	}
	$: {
		dataset.connection = connections.find((c) => c.id == dataset.connectionId)
	}
</script>

<template lang="pug">
.form-item
  label(for="connection")  {type=='source' ? 'Source' : 'Target'} Connection
  select(name='connection', bind:value='{dataset.connectionId}')
    +each('connections as c')
      option(value='{c.id}' selected='{dataset.connectionId === c.id}') {c.name}
+if('dataset.connectionId')
  .form-item
    label(for="name") Dataset name
    input.input(type="text", name="name", bind:value='{dataset.name}')
  .form-item
    <svelte:component type='{type}' bind:dataset='{dataset}' this={component()} connection={connections.find((c) => c.id == dataset.connectionId)} />

</template>
