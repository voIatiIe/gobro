var cursorX = 0.5, cursorY = 0.5;

setCursorPosition = (x, y) => { cursorX = x / window.innerWidth; cursorY = y / window.innerHeight; }

document.onmousemove = (event) => { setCursorPosition(event.clientX, event.clientY) };


newWSConnection = () => {
    const ws = new WebSocket('ws://localhost:8010/ws');

    logOutput = (message) => { document.getElementById('output').value += message + '\n' };

    ws.onopen = function () { logOutput('WebSocket connection opened') };
    ws.onclose = function () { logOutput('WebSocket connection closed') };
    ws.onmessage = function (event) { logOutput('Message received: ' + event.data) };
    ws.onerror = function (error) { logOutput('WebSocket Error: ' + error) };

    sendCursorPosition = () => {
        if (ws.readyState !== WebSocket.OPEN) { return }

        const message = JSON.stringify({ x: cursorX, y: cursorY });
        ws.send(message);
    }

    setInterval(sendCursorPosition, 60);
}

clearOutput = (e) => {
    e.preventDefault();
    document.getElementById('output').value = '';
}
