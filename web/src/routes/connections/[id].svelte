<script>
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { endpoint } from '$lib/utils';

	let connection = {};
	onMount(async () => {
		const response = await fetch(endpoint(`connections/${$page.params.id}`));
		connection = await response.json();
		console.log(connection);
	});
</script>

<template lang="pug">
  .w-full.h-full.flex.flex-col.pt-4
    .flex.justify-between.items-center
      .title
        | Connection: {connection?.name}
    .bg-white.mt-4.p-2.rounded-md
      form(action='#' method='POST')
        .px-4.py-5.bg-white(class='sm:p-6')
          .flex.flex-col
            .form-item
              label(for='street-address') Name
              input.mt-1(type='text' name='name' autocomplete='conneciton-name' value='{connection?.name}')
        .py-3.text-right(class='sm:px-6')
          button() Save


</template>
