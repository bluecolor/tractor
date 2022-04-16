<script>
	export let type = 'source';
	export let connection = null;
	export let dataset = {
		config: {}
	};
	$: if (dataset.config.path) {
		dataset.name = dataset.config.path.split('/').pop();
	}
</script>

<template lang="pug">
.form-item
  label(for="table") Path
  input.input(type="text", name="path", bind:value='{dataset.config.path}')
.form-item
  label(for="table") Delimiter
  input.input(type="text", name="path", bind:value='{dataset.config.delimiter}')

</template>
