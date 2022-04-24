<script>
	import Navbar from '@components/header/Navbar.svelte'
	import GlobalStatus from '@components/GlobalStatus.svelte'
	import '../assets/css/app.css'
	import '../assets/scss/index.scss'
	import 'sweetalert2/src/sweetalert2.scss'
</script>

<Navbar />
<GlobalStatus />
<div class="main-layout flex justify-center">
	<div class="main-layout-content flex w-3/5">
		<slot />
	</div>
</div>

<style>
	.main-layout {
		padding-top: 50px;
		/* height: calc(100vh - 50px); */
		background-color: var(--bg-color);
	}
</style>
