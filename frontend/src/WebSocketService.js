import { InviteContext } from './InviteContext';
import React, { useContext } from 'react';

class WebSocketService {
  static instance = null;
  callbacks = {};

  static getInstance() {
    if (!WebSocketService.instance) {
      WebSocketService.instance = new WebSocketService();
    }
    return WebSocketService.instance;
  }

  constructor() {
    this.socketRef = null;
  }

  connect(websocketUrl) {
    if (this.socketRef && this.socketRef.readyState !== WebSocket.CLOSED) {
      console.log('WebSocket connection already exists');
      return;
    }

    console.log('Connecting to WebSocket:', websocketUrl);
    this.socketRef = new WebSocket(websocketUrl);

    this.socketRef.onopen = () => {
      console.log('WebSocket connection established');
    };

    this.socketRef.onmessage = (event) => {
      console.log('WebSocket message received:', event.data);
      this.socketNewMessage(event.data);
    };

    this.socketRef.onclose = (event) => {
      console.log('WebSocket connection closed:', event);
    };

    this.socketRef.onerror = (error) => {
      console.log('WebSocket error:', error);
    };
  }

  socketNewMessage(data) {
    const parsedData = JSON.parse(data);
    const { action, args } = parsedData;
    if (this.callbacks[action]) {
      this.callbacks[action](args);
    }
  }

  addCallbacks(actions) {
    this.callbacks = { ...this.callbacks, ...actions };
  }

  sendMessage(data) {
    try {
      this.socketRef.send(JSON.stringify(data));
    } catch (err) {
      console.log(err.message);
    }
  }

  state() {
    return this.socketRef.readyState;
  }

  waitForSocketConnection(callback) {
    const socket = this.socketRef;
    const recursion = this.waitForSocketConnection;
    setTimeout(() => {
      if (socket.readyState === 1) {
        console.log('Connection is made');
        if (callback != null) {
          callback();
        }
        return;
      } else {
        console.log('wait for connection...');
        recursion(callback);
      }
    }, 1); // wait 1 millisecond for the connection...
  }
}

const WebSocketInstance = WebSocketService.getInstance();

export default WebSocketInstance;