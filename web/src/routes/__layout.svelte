<script>
	import Navbar from '@components/header/Navbar.svelte';
	import '../assets/css/app.css';
	import '../assets/scss/index.scss';
</script>

<Navbar />

<div class="main-layout flex justify-center">
	<div class="main-layout-content h-full flex items-center w-3/5">
		<slot />
	</div>
</div>

<style>
	.main-layout {
		height: calc(100vh - 50px);
		background-color: var(--tertiary-color);
	}
</style>
