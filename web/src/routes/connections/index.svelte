<script>
	import { onMount } from 'svelte';
	import { endpoint } from '$lib/utils';

	let connections = [];
	onMount(async () => {
		const response = await fetch(endpoint('connections'));
		connections = await response.json();
	});
</script>

<template lang="pug">
  .w-full.h-full.flex.flex-col.pt-4
    .flex.justify-between.items-center
      .title
        | Connections
      .search
        input(type="text" placeholder="Search")
    .bg-white.mt-4.p-2.rounded-md
      table.min-w-full
        thead.border-b
          tr
            th.text-sm.font-bold.text-gray-900.px-6.py-4.text-left(scope="col")
              | Name
            th.text-sm.font-bold.text-gray-900.px-6.py-4.text-left(scope="col")
              | Type
        tbody
          +each('connections as conn')
            tr(class="last:border-b-0  hover:bg-gray-50")
              td.text-sm.text-gray-900.font-light.px-6.py-4.whitespace-nowrap
                | {conn.name}
              td.text-sm.text-gray-900.font-light.px-6.py-4.whitespace-nowrap
                | {conn.connectionType.name}


</template>
