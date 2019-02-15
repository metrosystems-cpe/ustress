import React from 'react';
import ReactDOM from 'react-dom';
import './index.scss';
import App from './App';
import WebSocketService from './components/utils/wsservice';

export const WS = new WebSocketService(`ws://${window.location.host}/ustress/api/v1/ws`) 


ReactDOM.render(<App />, document.getElementById('root'));
