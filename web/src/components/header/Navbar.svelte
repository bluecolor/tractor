<script>
	import Tractor from '@icons/tractor.svg'
	import './navbar.css'
</script>

<template lang="pug">

nav.nav-top.bg-white.shadow-sm
  .max-w-full.mx-auto
    .flex.items-center.justify-center.h-12
      .justify-between.flex.items-center(class='w-3/5')
        a.cursor-pointer(href='/')
          Tractor(width='2em' class="nav-link")
        .hidden(class='md:block')
          .ml-10.flex.items-baseline.space-x-6
            a.navlink.text-gray-600(href='/connections')
              | Connections
            a.navlink.text-gray-600(href='/extractions')
              | Extractions
            a.navlink.text-gray-600(href='/sessions')
              | Sessions

</template>

<style>
	.nav-top {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 100;
	}
	.nav-top .navlink {
		-webkit-font-smoothing: antialiased;
		cursor: pointer;
	}
</style>
