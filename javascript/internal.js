export default {
  message_types: {
    welcome: "welcome",
    disconnect: "disconnect",
    ping: "ping",
    ack: "ack",
    single_chat: "single_chat",
    room_chat: "room_chat"
  },
  default_mount_path: "/anychat",
  protocols: ["anychat-v1-json", "anychat-unsupported"]
}
