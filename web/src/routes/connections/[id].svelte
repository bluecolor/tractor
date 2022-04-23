<script>
	import { page } from '$app/stores'
	import { goto } from '$app/navigation'
	import { onMount } from 'svelte'
	import { endpoint, api } from '$lib/utils'
	import FileConnection from './components/FileConnection.svelte'
	import MySQLConnection from './components/MySQLConnection.svelte'
	let loading = false
	let connection = {}
	let connectionTypes = []
	let connectionTypeId = undefined
	let connectionTypeCode = undefined
	const id = $page.params.id
	onMount(async () => {
		const [con, ctypes] = await Promise.all([
			fetch(endpoint('connections/' + id)),
			fetch(endpoint('connections/types'))
		])
		connection = await con.json()
		connectionTypes = await ctypes.json()
		connectionTypeId = connection.connectionTypeId
		connectionTypeCode = connection.connectionType.code
	})
	const components = {
		file: FileConnection,
		mysql: MySQLConnection
	}
	function onConnectionTypeChange(e) {
		connectionTypeId = e.target.value
		connectionTypeCode = connectionTypes.find((c) => c.id == connectionTypeId).code
	}
	function onSubmit(e) {
		let method = 'PUT'
		let resource = 'connections/' + connection.id
		loading = true
		api(method, resource, connection)
			.then((response) => {
				if (response.ok) {
					goto('/connections')
					alert('Connection saved')
				} else {
					response.text().then((text) => {
						alert('Error saving connection\n' + text)
					})
				}
			})
			.finally(() => {
				loading = false
			})
	}
	function onTest() {
		loading = true
		api('POST', 'connections/test', connection)
			.then((response) => {
				if (response.ok) {
					alert('Connection test successful')
				} else {
					response.text().then((text) => {
						alert('Connection test failed!\n' + text)
					})
				}
			})
			.finally(() => {
				loading = false
			})
	}
</script>

<template lang="pug">
  .w-full.h-full.flex.flex-col.pt-4
    .flex.justify-between.items-center
      .title
        | Connection: {connection?.name}
    .bg-white.mt-4.p-2.rounded-md.flex.items-center.justify-center
      form(action='#' method='POST' class="w-2/3" on:submit|preventDefault='{onSubmit}')
        .px-4.py-5.bg-white(class='sm:p-6')
          .flex.flex-col
            .form-item
              label(for='connection-name') Name
              input.input.mt-1(type='text' id="connection-name" name='name' autocomplete='conneciton-name' bind:value='{connection.name}')

            .form-item
              label(for='connection-type') Connection Type
              select#connection-type(aria-label='connection type' name='type', on:change='{onConnectionTypeChange}' bind:value='{connection.connectionTypeId}')
                +each('connectionTypes as ct')
                  option(value='{ct.id}'  selected='{ct.id === connection.connectionTypeId}' ) {ct.name}
            +if('connectionTypeCode !== undefined')
              <svelte:component this={components[connectionTypeCode]} bind:state={connection.config} />

        .py-3.text-right.space-x-4(class='sm:px-6')
          button.btn.warning(on:click|preventDefault='{onTest}' disabled='{loading}') {loading ? 'Testing...' : 'Test Connection'}
          button.btn() Save

</template>
