// Consult https://www.snowpack.dev to learn about these options
module.exports = {
	extends: '@sveltejs/snowpack-config',
	"proxy": {
		"/oto": "http://localhost:8080/oto",
	}
};