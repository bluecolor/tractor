<script>
	import PlayIcon from '@icons/play.svg'
	import MoreIcon from '@icons/more.svg'
	import FilterIcon from '@icons/filter.svg'
	import Dropdown from '@components/Dropdown.svelte'
	import { onMount } from 'svelte'
	import { api, wsendpoint } from '$lib/utils'
	import { session } from '$app/stores'
	import Pagination from '@components/Pagination.svelte'
	import _ from 'lodash'
	import RichSelect from '@components/RichSelect.svelte'
	import { createPopper } from '@popperjs/core'

	let extractions = []
	let connections = []
	let page = {}
	let sessionLog = ''
	export let q = ''
	let options = [
		{
			label: 'Delete',
			value: 'delete'
		},
		{
			label: 'Sessions',
			value: 'sessions'
		}
	]
	export let filters = {
		sc: '',
		tc: ''
	}
	let filtersOpen = false

	function onSearch() {
		_.debounce(async () => {
			onLoad({ q })
		}, 500)()
	}

	function onClearFilters() {
		filters = {
			sc: '',
			tc: ''
		}
		onLoad()
	}
	async function onLoad(params) {
		let url = 'extractions?' + new URLSearchParams(params)
		let result = await api('GET', url)
		if (!result.ok) {
			let error = await result.json()
			alert(error.error)
			return
		}
		page = await result.json()
		extractions = page.items.map((r) => {
			if (r.sessions.length > 0) {
				r.status = r.sessions[0].status
				r.session = r.sessions[0]
			} else {
				r.status = null
			}
			return r
		})
		extractions = [...(page.items || [])]
	}
	function updateExtraction(feed) {
		if (!feed.sessionId) {
			return
		}
		extractions = extractions.map((e) => {
			if (e.session?.id == feed.sessionId) {
				const status = feed.type.toLowerCase()
				if (['success', 'error', 'running'].indexOf(status) > -1) {
					e.status = status
				}
			}
			return e
		})
		extractions = [...extractions]
	}
	function subscribe() {
		const url = wsendpoint('session/feeds')
		const client = new WebSocket(url)
		client.addEventListener('open', () => {
			console.log('Connected to session feed')
		})
		client.addEventListener('message', (event) => {
			const feed = JSON.parse(event.data)
			if (feed.sender == 'Driver') {
				updateExtraction(feed)
			}
		})
		client.addEventListener('close', () => {
			console.log('Disconnected from session feed')
		})
	}
	function onRunExtraction(id) {
		api('POST', `extractions/${id}/run`).then(async (response) => {
			if (response.ok) {
				console.log('Extraction run')
				const session = await response.json()
				extractions = extractions.map((e) => {
					if (e.id === session.extraction.id) {
						e.status = session.status
						e.session = session
					}
					return e
				})
			} else {
				response.text().then((text) => {
					alert('Failed to run extraction\n' + text)
				})
			}
		})
	}
	function onDeleteExtraction(id) {
		let extraction = extractions.find((e) => e.id === id)
		let ok = confirm('Are you sure you want to delete this extraction? ' + extraction.name)
		if (ok) {
			api('DELETE', 'extractions/' + id).then((response) => {
				if (response.ok) {
					extractions = extractions.filter((e) => e.id !== id)
				} else {
					response.text().then((text) => {
						alert('Failed to delete extraction\n' + text)
					})
				}
			})
		}
	}
	function onDropdown(e, id) {
		switch (e.detail.value) {
			case 'delete':
				onDeleteExtraction(id)
				break
			case 'sessions':
				window.location.href = `/sessions?extraction=${id}`
		}
	}
	function onStatusClick(extraction) {
		const tooltip = document.querySelector('#statuslog')

		if (!tooltip.classList.contains('hidden')) {
			tooltip.classList.add('hidden')
			return
		}

		const status = document.querySelector('#status-' + extraction.id)
		tooltip.classList.remove('hidden')
		sessionLog = extraction.session.logs
		createPopper(status, tooltip, {
			placement: 'right',
			modifiers: [
				{
					name: 'offset',
					options: {
						offset: [0, 20]
					}
				}
			]
		})
	}
	onMount(async () => {
		let c = await api('GET', 'connections')
		connections = await c.json()
		onLoad({ page: 0, ...filters })
		subscribe()
	})
</script>

<template lang="pug">
  .w-full.h-full.flex.flex-col.pt-4
    .flex.justify-between.items-center
      .title
        | Extractions
      .search.space-x-2.inline-flex.items-center()
        .action
          <span class="action" on:click='{() => filtersOpen = !filtersOpen}'>
            FilterIcon(class="icon-btn")
          </span>
        .action
          input.input(type="text" placeholder="Search" bind:value="{q}" on:input="{onSearch}")
        a.action(href="/extractions/new")
          button.btn Add
    +if('filtersOpen')
      .bg-white.mt-4.p-2.rounded-md
        .flex.items-center.gap-x-2
          .form-item.w-full
            label(for="source") Source connection
            select(name='source', bind:value='{filters.sc}')
              +each('connections as c')
                option(value='{c.id}') {c.name}
          .form-item.w-full
            label(for="target") Target connection
            select(name='source', bind:value='{filters.tc}')
              +each('connections as c')
                option(value='{c.id}') {c.name}

        .actions.flex.justify-start.gap-x-3
          <button class="btn" on:click="{() => onLoad({page:0, ...filters})}"> Apply </button>
          <button class="btn danger" on:click="{() => onClearFilters()}">Clear </button>

    .bg-white.mt-4.p-2.rounded-md
      #statuslog.hidden.relative
        .bubble(role="tooltip")
          | {sessionLog}
        .triLeft
      table.min-w-full
        thead
          tr
            th(scope="col" align="left")
              | Name
            th(scope="col" align="left")
              | Source
            th(scope="col" align="left")
              | Target
            th(scope="col" align="left")
              | Status
            th.actions
        tbody
          +each('extractions as e')
            tr(class="last:border-b-0  hover:bg-gray-50")
              td
                a(href="/extractions/{e.id}")
                  | {e.name}
              td
                | {e.sourceDataset.name}<span class="text-gray-400">@{e.sourceDataset.connection.name} </span>
              td
                | {e.targetDataset.name}<span class="text-gray-400">@{e.targetDataset.connection.name} </span>
              td(align="left")
                <div id="{'status-' + e.id}" class="{e.status=='error' ? 'cursor-pointer': ''}" aria-describedby="tooltip" on:click="{() => onStatusClick(e)}">
                  +if('e.status === "success"')
                    .flex.justify-center.items-center.font-medium.py-1.px-2.bg-white.rounded-full.text-green-700.bg-green-100.border.border-green-300
                      .text-xs.font-normal.leading-none.max-w-full.flex-initial { e.status }
                    +elseif('e.status === "running"')
                      .flex.justify-center.items-center.font-medium.py-1.px-2.bg-white.rounded-full.text-blue-700.bg-blue-100.border.border-blue-300
                        .text-xs.font-normal.leading-none.max-w-full.flex-initial { e.status }
                    +elseif('e.status === "error"')
                      .flex.justify-center.items-center.font-medium.py-1.px-2.bg-white.rounded-full.text-red-700.bg-red-100.border.border-red-300
                        .text-xs.font-normal.leading-none.max-w-full.flex-initial { e.status }
                    +elseif('e.status != null')
                      .flex.justify-center.items-center.font-medium.py-1.px-2.bg-white.rounded-full.text-gray-700.bg-gray-100.border.border-gray-300
                        .text-xs.font-normal.leading-none.max-w-full.flex-initial { e.status }
                </div>

              td.actions
                div.flex.justify-end.items-center
                  span(on:click='{onRunExtraction(e.id)}')
                    PlayIcon(class="icon-btn")
                  <Dropdown label="Options" bind:options='{options}' on:select='{(x) => onDropdown(x, e.id)}'>
                    div(slot="button")
                      MoreIcon.icon-btn()
                  </Dropdown>
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

	#statuslog {
		z-index: 100;
	}
	.bubble {
		min-width: 20em;
		min-height: 20em;
		border-radius: 2px;
		padding: 0.5em;
		top: -1.7em;
		box-sizing: border-box;
		position: absolute;
		background: #888;
		font: 1em Candara, Calibri, Segoe, 'Segoe UI', Optima, Arial, sans-serif;
		text-align: center;
		color: #fff;
	}

	.triLeft {
		width: 0;
		height: 0;
		border-top: 15px solid transparent;
		border-right: 15px solid #888;
		border-bottom: 15px solid transparent;
		border-left: 15px solid transparent;

		position: absolute;
		left: -30px;
		top: -15px;
	}
</style>
