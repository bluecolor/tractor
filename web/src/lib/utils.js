import Swal from 'sweetalert2/dist/sweetalert2.js'

const endpoint = (path) => {
	return `http://localhost:3000/api/v1/${path}`
}
const wsendpoint = (path) => {
	return `ws://localhost:3000/api/v1/ws/${path}`
}

function api(method, resource, data) {
	return fetch(endpoint(resource), {
		method,
		headers: {
			'content-type': 'application/json'
		},
		body: data && JSON.stringify(data)
	})
}
function clickOutside(element, callbackFunction) {
	function onClick(event) {
		if (!element.contains(event.target)) {
			callbackFunction()
		}
	}

	document.body.addEventListener('click', onClick)

	return {
		update(newCallbackFunction) {
			callbackFunction = newCallbackFunction
		},
		destroy() {
			document.body.removeEventListener('click', onClick)
		}
	}
}

const Alert = Swal.mixin({
	customClass: {
		confirmButton: 'btn',
		cancelButton: 'btn btn-danger'
	},
	buttonsStyling: false,
	showClass: {
		popup: 'animate__animated animate__fadeInDown'
	},
	hideClass: {
		popup: 'animate__animated animate__fadeOutUp'
	}
})

const Toast = Swal.mixin({
	toast: true,
	position: 'top',
	showConfirmButton: false,
	timer: 3000,
	timerProgressBar: false,
	showClass: {
		popup: 'animate__animated animate__fadeInDown'
	},
	hideClass: {
		popup: 'animate__animated animate__fadeOutUp'
	},
	didOpen: (toast) => {
		toast.addEventListener('mouseenter', Swal.stopTimer)
		toast.addEventListener('mouseleave', Swal.resumeTimer)
	}
})

export { endpoint, wsendpoint, api, clickOutside, Alert, Toast }
