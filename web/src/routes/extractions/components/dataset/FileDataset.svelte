<script>
	export let connection = null;
	export let dataset = {
		config: {}
	};
	console.log(connection);
</script>

<template lang="pug">
.form-item
  label(for="table") Path
  input.input(type="text", name="path", bind:value='{dataset.config.path}')
.form-item
  label(for="table") Path
  input.input(type="text", name="path", bind:value='{dataset.config.path}')


</template>
