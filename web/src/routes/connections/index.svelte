<script>
	import Trash from '@icons/trash.svg'
	import { onMount } from 'svelte'
	import { api } from '$lib/utils'
	import Pagination from '@components/Pagination.svelte'
	import _ from 'lodash'

	let connections = []
	let page = {}
	export let filters = {
		q: ''
	}

	onMount(async () => {
		onLoad()
	})
	async function onLoad(params) {
		let url = 'connections?' + new URLSearchParams(params)
		let result = await api('GET', url)
		if (!result.ok) {
			let error = await result.json()
			alert(error.error)
			return
		}
		page = await result.json()
		connections = page.items
	}
	function onSearch() {
		_.debounce(async () => {
			onLoad(filters)
		}, 500)()
	}
	function onDeleteConnection(id) {
		let connection = connections.find((connection) => connection.id === id)
		let ok = confirm('Are you sure you want to delete this connection? ' + connection.name)
		if (ok) {
			api('DELETE', 'connections/' + id).then((response) => {
				if (response.ok) {
					connections = connections.filter((connection) => connection.id !== id)
				} else {
					response.text().then((text) => {
						alert('Failed to delete connection\n' + text)
					})
				}
			})
		}
	}
</script>

<template lang="pug">
  .w-full.h-full.flex.flex-col.pt-4
    .flex.justify-between.items-center
      .title
        | Connections
      .search.space-x-2.inline-flex.items-center()
        .action
          input.input(type="text" placeholder="Search" bind:value="{filters.q}" on:input="{onSearch}")
        a.action(href="/connections/new")
          button.btn Add

    .bg-white.mt-4.p-2.rounded-md.shadow-md
      table.min-w-full
        thead
          tr
            th(scope="col" align="left")
              | Name
            th(scope="col" align="left")
              | Type
            th.actions
        tbody
          +each('connections as conn')
            tr(class="last:border-b-0  hover:bg-gray-50")
              td
                a(href="/connections/{conn.id}")
                  | {conn.name}
              td
                | {conn.connectionType.name}
              td.actions(align="right")
                .flex.justify-end.items-center
                  span(on:click='{onDeleteConnection(conn.id)}')
                    Trash(class="trash")
      .mt-4
        <Pagination page='{page.page}' total='{page.total}' first='{page.first}' last='{page.last}' maxPage='{page.max_page}' visible='{page.visible}' on:paginate='{(p) => onLoad({page: p.detail, ...filters})}'/>

</template>

<style lang="postcss">
	table thead tr th {
		@apply font-normal text-base text-gray-500 pl-4 pr-4 pb-2 pt-2;
	}
	table tbody tr td {
		@apply font-normal text-base text-gray-700 pl-4 pr-4 pb-2 pt-2;
	}
</style>
