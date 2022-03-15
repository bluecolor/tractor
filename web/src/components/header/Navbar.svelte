<script>
	import Connection from '../../assets/icons/connection.svg';
	import Download from '../../assets/icons/download.svg';
	import Tractor from '../../assets/icons/tractor.svg';
	import Pulse from '../../assets/icons/pulse.svg';
	import './navbar.css';
</script>

<template lang="pug">

nav.nav-top.bg-white.border-b
  .max-w-full.mx-auto
    .flex.items-center.justify-center.h-12
      .justify-between.flex.items-center(class='w-3/5')
        a.cursor-pointer(href='/')
          Tractor(width='2em' class="nav-link")
        .hidden(class='md:block')
          .ml-10.flex.items-baseline.space-x-6
            a.navlink(href='/connections')
              Connection(width='1.4em')
            a.navlink(href='/extractions')
              Download(width='1.4em')
            a.navlink(href='/sessions')
              Pulse(width='1.4em')

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
