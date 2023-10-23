const buttons = document.getElementsByClassName('simple_buttons')
const mainEntry = document.getElementById('inputField')
const xInput= document.getElementById('input_x_value')
const xyMin = document.getElementById('input_min')
const xyMax = document.getElementById('input_max')

var sendRequest = function(message, type = 'POST', location = 'http://localhost:8080/') {
	const req = new XMLHttpRequest()
	req.open(type, location + "?body=" + '\''+ message +'\'')
	console.log(type, location + "?body=" + '\''+ message +'\'')
	req.send()
	req.onload = () => {
		if (message == '=')
			this.location.reload()
	}
}

var clickFunction = function() {
	if (this.getAttribute('id') === 'button_equals') {
		sendRequest('x= ' + xInput.value)
		sendRequest('=')
	} else if (this.getAttribute('value') === 'clear') {
		mainEntry.setAttribute('value', '')
		sendRequest('clear')
	} else if (this.getAttribute('id') === 'button_help') {
		// sendRequest(this.getAttribute('value'), 'GET', 'http://localhost:8080/help.html')
		location.href = 'help.html'
	} else if (this.getAttribute('id') === 'button_graph') {
		// location.href = 'graph_window.html'
		window.open('graph_window.html')
	} else {
		mainEntry.setAttribute('value',mainEntry.getAttribute('value') + this.getAttribute('value'))
		sendRequest(this.getAttribute('value'))
	}
}

for (var i = 0; i < buttons.length; i++) {
	buttons[i].addEventListener('click', clickFunction)
}

xyMin.addEventListener('focus', () => {
	sendRequest('xy_min= ' + xyMin.value)
})

xyMax.addEventListener('focus', () => {
	sendRequest('xy_max= ' + xyMax.value)
})


