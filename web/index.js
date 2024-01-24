var cursorX = 0.5, cursorY = 0.5;

const MessageType = {
    Move: 0,
    LeftClick: 1,
}

setCursorPosition = (x, y) => { cursorX = x / window.innerWidth; cursorY = y / window.innerHeight; }

document.onmousemove = (event) => { setCursorPosition(event.clientX, event.clientY) };


newWSConnection = () => {
    const ws = new WebSocket('ws://localhost:8010/ws');

    logOutput = (message) => {
        const output = document.getElementById('output');

        output.value += message + '\n';
        output.scrollTop = output.scrollHeight;
    };
    logOutput('Connecting to WebSocket server...');

    ws.onopen = function () { logOutput('WebSocket connection opened') };
    ws.onclose = function () { logOutput('WebSocket connection closed') };
    ws.onmessage = function (event) { displayImage(event) };
    ws.onerror = function (error) { logOutput('WebSocket Error: ' + error) };


    sendCommand = (command) => {
        if (ws.readyState !== WebSocket.OPEN) { return }

        const message = JSON.stringify(command);

        logOutput(message);

        ws.send(message);
    }

    setInterval(() => { return sendCommand({type: MessageType.Move, body: { x: cursorX, y: cursorY }}) }, 30);

    sendClick = (button, x, y) => {
        if (ws.readyState !== WebSocket.OPEN) { return }

        switch (button) {
            case 0:
                logOutput('Left click')
                sendCommand({type: MessageType.LeftClick, body: { x: x, y: y }});
                break;
            default:
                break;
        }
    }

    document.onclick = (event) => { sendClick(event.button, event.clientX, event.clientY) };
}

clearOutput = (e) => {
    e.preventDefault();
    document.getElementById('output').value = '';
}

displayImage = (e) => {
    let image = document.getElementById('image');

    image.src = window.URL.createObjectURL(e.data);
}