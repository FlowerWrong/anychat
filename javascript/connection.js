import adapters from "./adapters"
import ConnectionMonitor from "./connection_monitor"
import INTERNAL from "./internal"
import logger from "./logger"

// Encapsulate the cable connection held by the consumer. This is an internal class not intended for direct user manipulation.

const {message_types, protocols} = INTERNAL
const supportedProtocols = protocols.slice(0, protocols.length - 1)

const indexOf = [].indexOf

class Connection {
  constructor(consumer) {
    this.open = this.open.bind(this)
    this.consumer = consumer
    this.monitor = new ConnectionMonitor(this)
    this.disconnected = true
  }

  send(data) {
    if (this.isOpen()) {
      this.webSocket.send(JSON.stringify(data))
      return true
    } else {
      return false
    }
  }

  open() {
    if (this.isActive()) {
      logger.log(`Attempted to open WebSocket, but existing socket is ${this.getState()}`)
      return false
    } else {
      logger.log(`Opening WebSocket, current state is ${this.getState()}, subprotocols: ${protocols}`)
      if (this.webSocket) { this.uninstallEventHandlers() }
      this.webSocket = new adapters.WebSocket(this.consumer.url, protocols)
      this.installEventHandlers()
      this.monitor.start()
      return true
    }
  }

  close({allowReconnect} = {allowReconnect: true}) {
    if (!allowReconnect) { this.monitor.stop() }
    if (this.isActive()) {
      return this.webSocket.close()
    }
  }

  reopen() {
    logger.log(`Reopening WebSocket, current state is ${this.getState()}`)
    if (this.isActive()) {
      try {
        return this.close()
      } catch (error) {
        logger.log("Failed to reopen WebSocket", error)
      }
      finally {
        logger.log(`Reopening WebSocket in ${this.constructor.reopenDelay}ms`)
        setTimeout(this.open, this.constructor.reopenDelay)
      }
    } else {
      return this.open()
    }
  }

  getProtocol() {
    if (this.webSocket) {
      return this.webSocket.protocol
    }
  }

  isOpen() {
    return this.isState("open")
  }

  isActive() {
    return this.isState("open", "connecting")
  }

  // Private

  isProtocolSupported() {
    return indexOf.call(supportedProtocols, this.getProtocol()) >= 0
  }

  isState(...states) {
    return indexOf.call(states, this.getState()) >= 0
  }

  getState() {
    if (this.webSocket) {
      for (let state in adapters.WebSocket) {
        if (adapters.WebSocket[state] === this.webSocket.readyState) {
          return state.toLowerCase()
        }
      }
    }
    return null
  }

  installEventHandlers() {
    for (let eventName in this.events) {
      const handler = this.events[eventName].bind(this)
      this.webSocket[`on${eventName}`] = handler
      logger.log(`Install event handler ${handler} for ${eventName}`)
    }
  }

  onHandshake(callback) {
    logger.log("Handshake success")
    callback()
  }

  onMessage(cmd, ack, data) {
    logger.log("on message", cmd, ack, data)
  }

  onAckMessage(cmd, ack, data) {
    logger.log("on ack message", cmd, ack, data)
  }

  onChatMessage(cmd, ack, data) {
    logger.log("on chat message", cmd, ack, data)
  }

  onRoomChatMessage(cmd, ack, data) {
    logger.log("on room chat message", cmd, ack, data)
  }

  uninstallEventHandlers() {
    for (let eventName in this.events) {
      this.webSocket[`on${eventName}`] = function() {}
    }
  }

}

Connection.reopenDelay = 500

Connection.prototype.events = {
  message(event) {
    if (!this.isProtocolSupported()) { return }
    logger.log(event.data)
    const {cmd, ack, data} = JSON.parse(event.data)
    switch (cmd) {
      case message_types.welcome:
        return this.monitor.recordConnect()
      case message_types.disconnect:
        logger.log(`Disconnecting. Reason: ${data.reason}`)
        return this.close({allowReconnect: data.reconnect})
      case message_types.ping:
        return this.monitor.recordPing()
      case message_types.ack:
        return this.onAckMessage(cmd, ack, data)
      case message_types.single_chat:
        return this.onChatMessage(cmd, ack, data)
      case message_types.room_chat:
        return this.onRoomChatMessage(cmd, ack, data)
      default:
        return this.onMessage(cmd, ack, data)
    }
  },

  open() {
    logger.log(`WebSocket onopen event, using '${this.getProtocol()}' subprotocol`)
    this.disconnected = false
    if (!this.isProtocolSupported()) {
      logger.log("Protocol is unsupported. Stopping monitor and disconnecting.")
      return this.close({allowReconnect: false})
    }
  },

  close(event) {
    logger.log("WebSocket onclose event", event)
    if (this.disconnected) { return }
    this.disconnected = true
    this.monitor.recordDisconnect()
  },

  error() {
    logger.log("WebSocket onerror event")
  }
}

export default Connection