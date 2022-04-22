<script>
	import { onMount } from 'svelte'
	import { api } from '$lib/utils'
	import Dropdown from '@components/Dropdown.svelte'
	import TrashIcon from '@icons/trash.svg'
	import MoreIcon from '@icons/more.svg'
	import PlusIcon from '@icons/plus.svg'
	import GreaterThanIcon from '@icons/greater-than.svg'

	export let extraction

	let options = [
		{
			label: 'Fetch',
			value: 'fetch'
		},
		{
			label: 'Clear',
			value: 'clear'
		}
	]

	onMount(async () => {
		if (extraction?.sourceDataset?.fields && extraction?.targetDataset?.fields) {
			mapFields()
		}
	})
	function onDropdown(e) {
		const { value } = e.detail
		switch (value) {
			case 'fetch':
				onFetch()
				break
			case 'clear':
				onClear()
				break
		}
	}
	function mapFields() {
		extraction.sourceDataset.fields = (extraction.sourceDataset?.fields ?? []).map((s, i) => {
			return {
				order: i,
				...s
			}
		})
		extraction.targetDataset.fields = (extraction.targetDataset?.fields ?? []).map((t, i) => {
			return {
				order: i,
				...t
			}
		})

		while (extraction.sourceDataset.fields.length < extraction.targetDataset.fields.length) {
			extraction.sourceDataset.fields.push({ order: extraction.sourceDataset.fields.length })
		}
		while (extraction.targetDataset.fields.length < extraction.sourceDataset.fields.length) {
			extraction.targetDataset.fields.push({ order: extraction.targetDataset.fields.length })
		}
		extraction.sourceDataset.fields =
			extraction.sourceDataset.fields.sort((a, b) => a.order - b.order) || []
		extraction.targetDataset.fields =
			extraction.targetDataset.fields.sort((a, b) => a.order - b.order) || []
	}
	function onFetch() {
		Promise.all([
			api('POST', `connections/${extraction.sourceDataset.connectionId}/dataset`, {
				...extraction.sourceDataset.connection.config,
				...extraction.sourceDataset.config
			}),
			api('POST', `connections/${extraction.targetDataset.connectionId}/dataset`, {
				...extraction.targetDataset.connection.config,
				...extraction.targetDataset.config
			})
		]).then(async (responses) => {
			if (responses.every((r) => r.ok)) {
				let s = await responses[0].json()
				let t = await responses[1].json()
				extraction.sourceDataset.fields = s.fields
				extraction.targetDataset.fields = t.fields
				mapFields()
			} else {
				let errms = await Promise.all(responses.map((r) => r.text()))
				alert('Failed to load fields\n' + errms.join('\n'))
			}
		})
	}
	function onDeleteMapping(i) {
		extraction.sourceDataset.fields.splice(i, 1)
		extraction.targetDataset.fields.splice(i, 1)
		extraction.sourceDataset.fields = extraction.sourceDataset.fields.map((s, i) => {
			return {
				...s,
				order: i
			}
		})
		extraction.targetDataset.fields = extraction.targetDataset.fields.map((t, i) => {
			return {
				...t,
				order: i
			}
		})
	}
	function onAddMapping() {
		extraction.sourceDataset.fields.push({ order: extraction.sourceDataset.fields.length })
		extraction.targetDataset.fields.push({ order: extraction.targetDataset.fields.length })
	}
	function onClear() {
		mappings = []
	}
</script>

<template lang="pug">
.mappings
  .mappings
    table.table.min-w-full
      thead
        tr
          th(scope="col" align="left")
            | Source column
          th(scope="col" align="left")
            | Source type
          th(scope="col" align="left")
            | Target column
          th(scope="col" align="left")
            | Target type
          th.actions.flex.justify-end.items-center(align="right")
            .action.icon-btn.mr-3(on:click='{onAddMapping}')
              PlusIcon()
            Dropdown(label="Options" bind:options='{options}' on:select='{onDropdown}')
              div(slot="button")
                MoreIcon.icon-btn()
      tbody
        +each('extraction.sourceDataset.fields as s, i')
          tr(class="last:border-b-0  hover:bg-blue-50")
            td
              input.input(placeholder="Source column", bind:value='{s.name}')
            td
              span.text-gray-600
                | {s.type}
            td
              input.input(placeholder="Target column", bind:value='{extraction.targetDataset.fields[i].name}')
            td
              select.cursor-pointer()
                option(value="string") string
                option(value="integer") integer
                option(value="float") float
                option(value="boolean") boolean
                option(value="date") date
            td
              div.flex.justify-end.items-center
                <span class="cursor-pointer" on:click='{() => onDeleteMapping(i)}'>
                  GreaterThanIcon(class="fill-current text-gray-200 hover:text-blue-500")
                </span>
                <span class="cursor-pointer" on:click='{() => onDeleteMapping(i)}'>
                  TrashIcon(class="trash")
                </span>

</template>

<style lang="postcss">
	table thead tr th {
		@apply font-normal text-base text-gray-500 pl-4 pr-4 pb-2 pt-2;
	}
	table tbody tr td {
		@apply font-normal text-base text-gray-700 pl-4 pr-4 pb-2 pt-2;
	}
</style>
