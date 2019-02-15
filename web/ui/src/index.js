import React from 'react';
import ReactDOM from 'react-dom';
import './index.scss';
import App from './App';
import WebSocketService from './components/utils/wsservice';


export const CurrentDomain = `${window.location.protocol}//${window.location.host}`;

var wsUrl = CurrentDomain.indexOf("localhost") != -1 ? 
    `ws://${window.location.host}/ustress/api/v1/ws` : 
    `wss://${window.location.host}/ustress/api/v1/ws`;

export const WS = new WebSocketService(wsUrl);



ReactDOM.render(<App />, document.getElementById('root'));
