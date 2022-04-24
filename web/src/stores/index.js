import { writable } from 'svelte/store'

export const globalStatus = writable({
	open: false,
	text: 'Error',
	type: 'error'
})
