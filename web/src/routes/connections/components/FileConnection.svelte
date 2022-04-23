<script>
	import { onMount } from 'svelte'
	import { api } from '$lib/utils'
	import _ from 'lodash'
	export let state = null

	let providers = []
	onMount(async () => {
		const response = await api('GET', `connections/providers`)
		providers = await response.json()
	})

	const fileFormats = [
		{ name: 'CSV', code: 'csv' },
		{ name: 'JSON', code: 'json' }
	]
</script>

<template lang="pug">
  .form-item
    label(for='file-provider') Provider
    select#file-provider(aria-label='provider' bind:value='{state.provider.code}')
      +each('providers as p')
        option(value='{p.code}' selected='{p.code === state?.provider?.code}') {p.name}
  .form-item
    label(for='file-format') Format
    select#file-format(aria-label='file format' bind:value='{state.format.code}')
      +each('fileFormats as ff')
        option(value='{ff.code}' selected='{ff.code === state?.format}') {ff.name}
</template>
