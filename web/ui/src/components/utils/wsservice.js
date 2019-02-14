import { Observable, Subject } from 'rxjs';
import { map } from 'rxjs/operators';


class WebSocketService {
  constructor(url) {
    this.status = false;
    this.url = url;
    if (url !== undefined) {
      this.init(url)
    }
  }

  create(url) {
    this.ws = new WebSocket(url)
    let observable = Observable.create((obs => {
      this.ws.onmessage = obs.next.bind(obs);
      this.ws.onerror = obs.error.bind(obs);
      this.ws.onclose = () => {
        this.status = false
        console.log("Websocket disconnected")

      }
    }))

    let observer = {
      next: data => {
        if (this.ws.readyState === WebSocket.OPEN) {
          this.timestamp = new Date();
          this.ws.send(JSON.stringify(data))
        } else {
          this.init(this.url)
          this.status = false
          console.warn("Websocket Disconnected");
        }
      }
    }

    return Subject.create(observer, observable)
  }

  init(url) {
    this.status = true
    this.feed = this.create(url).pipe(map(res => {
      // For some reason the output after parsing JSON is still string
      // TODO dig deeper
      let data = JSON.parse(res.data);
      return typeof data === 'string' ? JSON.parse(data) : data;
    }))
  }
}

export default WebSocketService