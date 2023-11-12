const buttons = document.getElementsByClassName('simple_buttons')
const mainEntry = document.getElementById('inputField')
const xInput= document.getElementById('input_x_value')
const xyMin = document.getElementById('input_min')
const xyMax = document.getElementById('input_max')

let start = true


var sendRequest = function(message, type = 'POST', location = 'http://localhost:8080/') {
	const req = new XMLHttpRequest()
	if (message === " + ") {
		message = " plus "
	}
	if (message === " / ") {
		message = " divide "
	}
	req.open(type, location + "?body=" + '\''+ message +'\'')
	console.log(type, location + "?body=" + '\''+ message +'\'')
	req.send()
	req.onload = () => {
		if (message == '=' || message == '')
			this.location.reload()
	}
}

var clickFunction = function() {
	if (start) {
		start = false
		mainEntry.setAttribute('value', '')
	}
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
	} else if (this.getAttribute('id') === 'button_next' || this.getAttribute('id') === 'button_last'
	|| this.getAttribute('id') === 'button_prev' || this.getAttribute('id') === 'button_history_clear') {
		sendRequest(this.getAttribute('id'))
		sendRequest('', "GET")
	} else {
		mainEntry.setAttribute('value',mainEntry.getAttribute('value') + this.getAttribute('value'))
		sendRequest(this.getAttribute('value'))
	}
}

for (var i = 0; i < buttons.length; i++) {
	buttons[i].addEventListener('click', clickFunction)
}

xyMin.addEventListener('focusout', () => {
	sendRequest('xy_min= ' + xyMin.value)
})

xyMax.addEventListener('focusout', () => {
	sendRequest('xy_max= ' + xyMax.value)
})


