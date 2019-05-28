export default {
  message_types: {
    welcome: 0,
    disconnect: 2,
    ping: 4,
    ack: 11,
    single_chat: 101,
    room_chat: 102
  },
  default_mount_path: "/anychat",
  protocols: ["anychat-v1-json", "anychat-unsupported"]
}
