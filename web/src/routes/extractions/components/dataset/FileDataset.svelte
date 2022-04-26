<script>
	export let type = 'source'
	export let connection = null
	export let dataset = null
	dataset.config = dataset.config || {}
	if (!dataset.config.delimiter) {
		dataset.config.delimiter = ','
	}
	$: if (dataset?.config?.path && !dataset.name) {
		dataset.name = dataset.config.path.split('/').pop()
	}
</script>

<template lang="pug">
.form-item
  label(for="table") Path
  input.input(type="text", name="path", bind:value='{dataset.config.path}')
+if('connection.config.format.code == "csv"')
  .form-item
    label(for="table") Delimiter
    input.input(type="text", name="path", bind:value='{dataset.config.delimiter}')

</template>
