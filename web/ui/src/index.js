import React from 'react';
import ReactDOM from 'react-dom';
import './index.scss';
import App from './App';
import WebSocketService from './components/utils/wsservice';



const random_rgba = () => {
    var o = Math.round, r = Math.random, s = 255;
    return 'rgba(' + o(r()*s) + ',' + o(r()*s) + ',' + o(r()*s) + ',' + '0.5' + ')';
}

const cpool = [];


for (let index = 0; index < 100; index++) {
    cpool.push(random_rgba())
}
export const color_pool = cpool;


// Put to false before pushing to master
const dev = false;

export const CurrentDomain = !dev ? `${window.location.protocol}//${window.location.host}`: 'http://localhost:8080';

var wsUrl = CurrentDomain.indexOf("localhost") != -1 ? 
    `ws://localhost:8080/ustress/api/v1/ws` : 
    `wss://${window.location.host}/ustress/api/v1/ws`;

export const WS = new WebSocketService(wsUrl);



ReactDOM.render(<App />, document.getElementById('root'));

