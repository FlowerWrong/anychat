export default {
  message_types: {
    welcome: "welcome",
    ping: "ping",
    confirmation: "confirm_subscription",
    rejection: "reject_subscription"
  },
  default_mount_path: "/ws",
  protocols: ["actioncable-v1-json", "actioncable-unsupported"]
}
