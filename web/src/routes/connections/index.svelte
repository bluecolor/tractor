<script>
	import Trash from '../../assets/icons/trash.svg';

	import { onMount } from 'svelte';
	import { endpoint, api } from '$lib/utils';

	let connections = [];
	onMount(async () => {
		const response = await fetch(endpoint('connections'));
		connections = await response.json();
	});

	function onDeleteConnection(id) {
		let connection = connections.find((connection) => connection.id === id);
		let ok = confirm('Are you sure you want to delete this connection? ' + connection.name);
		if (ok) {
			api('DELETE', 'connections/' + id).then((response) => {
				if (response.ok) {
					connections = connections.filter((connection) => connection.id !== id);
				} else {
					response.text().then((text) => {
						alert('Failed to delete connection\n' + text);
					});
				}
			});
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
          input(type="text" placeholder="Search")
        .action
          button Add

    .bg-white.mt-4.p-2.rounded-md
      table.min-w-full
        thead.border-b
          tr
            th.text-sm.font-bold.text-gray-900.px-6.py-4.text-left(scope="col")
              | Name
            th.text-sm.font-bold.text-gray-900.px-6.py-4.text-left(scope="col")
              | Type
            th.actions
        tbody
          +each('connections as conn')
            tr(class="last:border-b-0  hover:bg-gray-50")
              td.text-sm.text-gray-900.font-light.px-6.py-4.whitespace-nowrap
                a(href="/connections/{conn.id}")
                  | {conn.name}
              td.text-sm.text-gray-900.font-light.px-6.py-4.whitespace-nowrap
                | {conn.connectionType.name}
              td.actions
                span(on:click='{onDeleteConnection(conn.id)}')
                  Trash(class="w-6 h-6 ml-2 cursor-pointer fill-current hover:text-red-500")

</template>

<style lang="scss">
	tr:hover a {
		text-decoration: underline;
	}
</style>
