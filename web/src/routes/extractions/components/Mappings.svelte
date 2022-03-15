<script>
	import { api } from '$lib/utils';

	export let sourceConnection, targetConnection;
</script>

<template lang="pug">
.mappings
  .toolbar.flex.justify-between.items-center
    .sub-title Mappings
    button.warning() Reset
</template>
