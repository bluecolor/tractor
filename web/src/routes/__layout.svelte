<script>
	import Navbar from '@components/header/Navbar.svelte'
	import GlobalStatus from '@components/GlobalStatus.svelte'
	import '../assets/css/app.css'
	import '../assets/scss/index.scss'
	import 'sweetalert2/src/sweetalert2.scss'
	import { loading } from '@stores'
	import PageLoading from '@components/PageLoading.svelte'

	let isLoading = false
	loading.subscribe((value) => {
		isLoading = value
	})
</script>

<template lang="pug">
Navbar
GlobalStatus
div(class="{isLoading ? '' : 'hidden'}")
  PageLoading
.main-layout.flex.justify-center
  .main-layout-content.flex(class="{'w-3/5 ' + (isLoading ? 'hidden' : '')}")
    slot
</template>

<style>
	.main-layout {
		padding-top: 50px;
		/* height: calc(100vh - 50px); */
		background-color: var(--bg-color);
	}
</style>
