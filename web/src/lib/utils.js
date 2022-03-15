const endpoint = (path) => {
	return `http://localhost:3000/api/v1/${path}`;
};

function api(method, resource, data) {
	return fetch(endpoint(resource), {
		method,
		headers: {
			'content-type': 'application/json'
		},
		body: data && JSON.stringify(data)
	});
}
function clickOutside(element, callbackFunction) {
	function onClick(event) {
		if (!element.contains(event.target)) {
			callbackFunction();
		}
	}

	document.body.addEventListener('click', onClick);

	return {
		update(newCallbackFunction) {
			callbackFunction = newCallbackFunction;
		},
		destroy() {
			document.body.removeEventListener('click', onClick);
		}
	};
}

export { endpoint, api, clickOutside };
