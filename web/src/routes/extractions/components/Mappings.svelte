<script>
	import { onMount } from 'svelte'
	import { api } from '$lib/utils'
	import Dropdown from '@components/Dropdown.svelte'
	import TrashIcon from '@icons/trash.svg'
	import MoreIcon from '@icons/more.svg'
	import PlusIcon from '@icons/plus.svg'
	import GreaterThanIcon from '@icons/greater-than.svg'

	export let extraction = {
		sourceDataset: {
			fields: []
		},
		targetDataset: {
			fields: []
		}
	}

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
	export let mappings = []
	$: {
		extraction.sourceDataset.fields = mappings.map((m, i) => ({ ...m.source, order: i }))
		extraction.targetDataset.fields = mappings.map((m, i) => ({ ...m.target, order: i }))
	}

	onMount(async () => {
		if (extraction?.sourceDataset?.fields && extraction?.targetDataset?.fields) {
			mapFields({
				sf: extraction.sourceDataset.fields,
				tf: extraction.targetDataset.fields
			})
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
	function mapFields({ sf, tf }) {
		let sourceFields = (sf ?? source.fields).map((s, i) => {
			return {
				order: i,
				...s
			}
		})
		let targetFields = (tf ?? target.fields).map((t, i) => {
			return {
				order: i,
				...t
			}
		})

		while (sourceFields.length < targetFields.length) {
			sourceFields.push({ order: sourceFields.length })
		}
		while (targetFields.length < sourceFields.length) {
			targetFields.push({ order: targetFields.length })
		}
		sourceFields = sourceFields.sort((a, b) => a.order - b.order) || []
		targetFields = targetFields.sort((a, b) => a.order - b.order) || []
		mappings = sourceFields.map((sf, i) => {
			return {
				source: sf,
				target: targetFields[i]
			}
		})
		mappings = [...mappings]
	}
	function onFetch() {
		Promise.all([
			api('POST', `connections/${source.connectionId}/dataset`, source.config),
			api('POST', `connections/${target.connectionId}/dataset`, target.config)
		]).then(async (responses) => {
			if (responses.every((r) => r.ok)) {
				let s = await responses[0].json()
				let t = await responses[1].json()
				mapFields({ sf: s.fields, tf: t.fields })
			} else {
				let errms = await Promise.all(responses.map((r) => r.text()))
				alert('Failed to load fields\n' + errms.join('\n'))
			}
		})
	}
	function onDeleteMapping(i) {
		mappings.splice(i, 1)
		mappings = [...mappings]
	}
	function onAddMapping() {
		mappings.push({
			source: { order: mappings.length },
			target: { order: mappings.length }
		})
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
        +each('mappings as m, i')
          tr(class="last:border-b-0  hover:bg-blue-50")
            td
              input.input(placeholder="Source column", bind:value='{m.source.name}')
            td
              span.text-gray-600
                | {m.source.type}
            td
              input.input(placeholder="Target column", bind:value='{m.target.name}')
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
