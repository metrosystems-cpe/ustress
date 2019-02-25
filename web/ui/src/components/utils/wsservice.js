import { Observable, Subject } from 'rxjs';
import { map } from 'rxjs/operators';
import { withSnackbar } from 'notistack';


class WebSocketService {
  constructor(url) {
    this.status = false;
    this.url = url;
    if (url !== undefined) {
      this.init(url)
    }
    this.enqueueSnackbar = (m, o) => {}
  }

  reconnect(delay) {
    setTimeout(() => {
      console.warn("Reconnecting ws")
      this.create(this.url)
    }, delay)
  }

  create(url) {

    this.ws = new WebSocket(url)

    let observable = Observable.create((obs => {
      this.ws.onmessage = obs.next.bind(obs);
      this.ws.onerror = obs.error.bind(obs);
      this.ws.onopen = () => {
        this.enqueueSnackbar("Websocket connected", {variant:"success"})
        this.status = true
      }
      this.ws.onclose = () => {
        this.status = false
        this.enqueueSnackbar("Websocket disconnected", {variant: "error"})
        this.reconnect(2500)
      }
    }))

    let observer = {
      next: data => {
        if (this.ws.readyState === WebSocket.OPEN) {
          this.timestamp = new Date();
          this.ws.send(JSON.stringify(data))
          this.enqueueSnackbar("Request successfuly sent", {variant: "success"})
        } else {
          this.status = false
          this.reconnect(2500)
          this.enqueueSnackbar("Websocket Disconnected", {variant: "error"})
          console.warn("Websocket Disconnected");
        }
      }
    }

    return Subject.create(observer, observable)
  }

  init(url) {
    this.feed = this.create(url).pipe(map(res => {
      // For some reason the output after parsing JSON is still string
      // TODO dig deeper
      let data = JSON.parse(res.data);
      return typeof data === 'string' ? JSON.parse(data) : data;
    }))
  }
}

export default WebSocketService