<script>
	import CloseIcon from '@icons/close.svg'
	import { globalStatus } from '@stores'

	let open = false
	let text = ''
	let type = ''
	globalStatus.subscribe((status) => {
		open = status.open
		text = status.text
		type = status.type
	})

	function onClose() {
		globalStatus.set({ open: false })
	}
</script>

<template lang="pug">
+if('open')
  .global-status(class="{type}")
    .global-status__body.flex.justify-center
      .global-status__message
        | {text}
      .close.cursor-pointer(on:click="{onClose}")
        CloseIcon(class="close-icon tex-white fill-current")
</template>

<style>
	.global-status {
		top: 50px;
		position: fixed;
		width: 100%;
		height: 48px;
		color: white;
	}
	.global-status.error {
		background-color: #e1817c;
	}
	.global-status__body {
		width: 100%;
		padding: 11px 36px;
		line-height: 24px;
		font-size: 1.1rem;
		text-align: center;
		margin: auto;
	}
	.global-status .close {
		position: absolute;
		right: 0;
		margin-right: 10px;
	}
	.global-status svg.close-icon {
		color: #fff !important;
		fill: #fff !important;
	}
</style>
