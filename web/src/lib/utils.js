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

export { endpoint, api };
