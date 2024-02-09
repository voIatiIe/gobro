import { MouseEvent, WheelEvent, FormEvent, KeyboardEvent } from 'react';
import { clamp } from '../lib/utils';
import { debounce } from "lodash";

enum MessageType {
    MOUSE_MOVE = 'MOUSE_MOVE',
    MOUSE_CLICK = 'MOUSE_CLICK',
    RESIZE = 'RESIZE',
}

type Coordinates = { x: number, y: number }

interface Command {
    type: MessageType,
    payload: {
        coordinates?: Coordinates,
        viewportParameters?: ViewportParameters,
    }
}

interface ViewportParameters {
    width: number,
    height: number,
    top: number,
    left: number,
}

export class GobroClient {
    private socket: WebSocket;
    private viewportParameters: ViewportParameters;
    private _mouse: Coordinates = { x: 0, y: 0 };
    private lastMouseSent: Coordinates = { x: 0, y: 0 };

    set mouse(value: Coordinates) { this._mouse = value }

    get mouse(): Coordinates {
        const params = this.getViewportParameters();

        return { 
            x: clamp((this._mouse.x - params.left) / params.width),
            y: clamp((this._mouse.y - params.top) / params.height),
        }
    }

    constructor(
        private getViewportParameters: () => ViewportParameters,
        private displayImage: (src: string) => void,
    ) {
        this.viewportParameters = getViewportParameters();

        this.socket = new WebSocket('ws://localhost:8010/ws');

        window.addEventListener('resize', debounce(this.onResize, 500));

        const tracker = this.startMouseTracking();

        this.socket.onmessage = (e) => displayImage(window.URL.createObjectURL(e.data))
        this.socket.onclose = () => clearInterval(tracker);
    }

    onMouseMove = (e: MouseEvent<HTMLImageElement>) => {
        this.mouse = { x: e.clientX, y: e.clientY };
    }

    onClick = (e: MouseEvent<HTMLImageElement>) => {
        this.mouse = { x: e.clientX, y: e.clientY }

        this.syncMouse()
        this.sendCommand({
            type: MessageType.MOUSE_CLICK,
            payload: { coordinates: this.mouse },
        });
    }

    onResize = () => {
        this.sendCommand({
            type: MessageType.RESIZE,
            payload: { viewportParameters: this.getViewportParameters() } 
        });
    }

    sendCommand = (command: Command) => {
        if (this.socket.readyState !== WebSocket.OPEN) { return }

        this.socket.send(JSON.stringify(command));
    }

    startMouseTracking = () => {
        return setInterval(() => {
            if (
                this.lastMouseSent.x === this.mouse.x
                && this.lastMouseSent.y === this.mouse.y
            ) { return }

            this.syncMouse()
            this.sendCommand({
                type: MessageType.MOUSE_MOVE,
                payload: { coordinates: this.mouse },
            });
        }, 30);
    }

    syncMouse = () => {
        this.lastMouseSent.x = this.mouse.x;
        this.lastMouseSent.y = this.mouse.y;
    }
}
