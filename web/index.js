var cursorX = 0.5, cursorY = 0.5;
var cursorLastSentX = 0.5, cursorLastSentY = 0.5;

const MessageType = {
    Move: 0,
    LeftClick: 1,
}

clamp = (number, min, max) => { return Math.max(min, Math.min(number, max)) }

setCursorPosition = (x, y) => { cursorX = x; cursorY = y; }

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

    setInterval(() => {
        if (cursorLastSentX === cursorX && cursorLastSentY === cursorY) { return }

        cursorLastSentX = cursorX; cursorLastSentY = cursorY;

        sendClick(null, cursorX, cursorY);
    }, 30);

    sendClick = (button, x, y) => {
        const image = document.getElementById('image');

        const rect = image.getBoundingClientRect()

        x -= rect.left;
        y -= rect.top;

        body = { x: clamp(x / rect.width, 0.0, 1.0), y: clamp(y / rect.height, 0.0, 1.0) };

        switch (button) {
            case null:
                sendCommand({type: MessageType.Move, body: body});
                break;
            case 0:
                sendCommand({type: MessageType.LeftClick, body: body});
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
    const image = document.getElementById('image');

    image.src = window.URL.createObjectURL(e.data);
}