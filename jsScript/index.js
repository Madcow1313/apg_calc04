const buttons = document.getElementsByClassName('simple_buttons')
const mainEntry = document.getElementById('inputField')
const xInput= document.getElementById('input_x_value')
const xyMax = document.getElementById('input_max')
const xyMin = document.getElementById('input_min')

var sendRequest = function(message) {
	const req = new XMLHttpRequest()
	req.open("POST", 'http://localhost:8080/index.html' + "?body=" + '\''+ message +'\'')
	req.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
	console.log(message)
	req.send()
}

xyMax.addEventListener('change', () => {
	sendRequest('xy_max= ' + xyMax.value)
})

xyMin.addEventListener('change', () => {
	sendRequest('xy_min= ' + xyMin.value)
})

var clickFunction = function() {
	if (this.getAttribute('id') === 'button_equals') {
		sendRequest('x= ' + xInput.value)
	} else if (this.getAttribute('id') === 'button_graph') {
		console.log(xyMax.value, xyMin.value)
	} else if (this.getAttribute('value') === 'clear') {
		mainEntry.setAttribute('value', '')
	}
	else {
		mainEntry.setAttribute('value',mainEntry.getAttribute('value') + this.getAttribute('value'))
	}
	sendRequest(this.getAttribute('value'))
}

for (var i = 0; i < buttons.length; i++) {
	buttons[i].addEventListener('click', clickFunction)
}



