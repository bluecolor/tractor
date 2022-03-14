<script>
	import { onMount } from 'svelte';
	import { endpoint } from '$lib/utils';

	let ptypes = [];
	onMount(async () => {
		const response = await fetch(endpoint(`connections/providers/types`));
		ptypes = await response.json();
		console.log(ptypes);
	});

	const fileFormats = [
		{ name: 'CSV', code: 'csv' },
		{ name: 'JSON', code: 'json' }
	];
	export let state;
</script>

<template lang="pug">
  .form-item
    label(for='file-provider') Provider
    select#file-provider(aria-label='provider')
      +each('ptypes as p')
        option(value='{p.code}' selected='{p.code === state?.provider}') {p.name}
  .form-item
    label(for='file-format') Format
    select#file-format(aria-label='file format')
      +each('fileFormats as ff')
        option(value='{ff.code}' selected='{ff.code === state?.format}') {ff.name}
</template>
