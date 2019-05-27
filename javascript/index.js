import Connection from "./connection"
import ConnectionMonitor from "./connection_monitor"
import Consumer, { createWebSocketURL } from "./consumer"
import INTERNAL from "./internal"
import adapters from "./adapters"
import logger from "./logger"

export {
  Connection,
  ConnectionMonitor,
  Consumer,
  INTERNAL,
  adapters,
  createWebSocketURL,
  logger,
}

export function createConsumer(url = INTERNAL.default_mount_path) {
  return new Consumer(url)
}