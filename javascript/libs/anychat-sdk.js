(function(){function r(e,n,t){function o(i,f){if(!n[i]){if(!e[i]){var c="function"==typeof require&&require;if(!f&&c)return c(i,!0);if(u)return u(i,!0);var a=new Error("Cannot find module '"+i+"'");throw a.code="MODULE_NOT_FOUND",a}var p=n[i]={exports:{}};e[i][0].call(p.exports,function(r){var n=e[i][1][r];return o(n||r)},p,p.exports,r,e,n,t)}return n[i].exports}for(var u="function"==typeof require&&require,i=0;i<t.length;i++)o(t[i]);return o}return r})()({1:[function(require,module,exports){
"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports["default"] = void 0;
var _default = {
  logger: self.console,
  WebSocket: self.WebSocket
};
exports["default"] = _default;

},{}],2:[function(require,module,exports){
"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports["default"] = void 0;

var _adapters = _interopRequireDefault(require("./adapters"));

var _connection_monitor = _interopRequireDefault(require("./connection_monitor"));

var _internal = _interopRequireDefault(require("./internal"));

var _logger = _interopRequireDefault(require("./logger"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { "default": obj }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } }

function _createClass(Constructor, protoProps, staticProps) { if (protoProps) _defineProperties(Constructor.prototype, protoProps); if (staticProps) _defineProperties(Constructor, staticProps); return Constructor; }

// Encapsulate the cable connection held by the consumer. This is an internal class not intended for direct user manipulation.
var message_types = _internal["default"].message_types,
    protocols = _internal["default"].protocols;
var supportedProtocols = protocols.slice(0, protocols.length - 1);
var indexOf = [].indexOf;

var Connection =
/*#__PURE__*/
function () {
  function Connection(consumer) {
    _classCallCheck(this, Connection);

    this.open = this.open.bind(this);
    this.consumer = consumer;
    this.monitor = new _connection_monitor["default"](this);
    this.disconnected = true;
  }

  _createClass(Connection, [{
    key: "send",
    value: function send(data) {
      if (this.isOpen()) {
        this.webSocket.send(JSON.stringify(data));
        return true;
      } else {
        return false;
      }
    }
  }, {
    key: "open",
    value: function open() {
      if (this.isActive()) {
        _logger["default"].log("Attempted to open WebSocket, but existing socket is ".concat(this.getState()));

        return false;
      } else {
        _logger["default"].log("Opening WebSocket, current state is ".concat(this.getState(), ", subprotocols: ").concat(protocols));

        if (this.webSocket) {
          this.uninstallEventHandlers();
        }

        this.webSocket = new _adapters["default"].WebSocket(this.consumer.url, protocols);
        this.installEventHandlers();
        this.monitor.start();
        return true;
      }
    }
  }, {
    key: "close",
    value: function close() {
      var _ref = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {
        allowReconnect: true
      },
          allowReconnect = _ref.allowReconnect;

      if (!allowReconnect) {
        this.monitor.stop();
      }

      if (this.isActive()) {
        return this.webSocket.close();
      }
    }
  }, {
    key: "reopen",
    value: function reopen() {
      _logger["default"].log("Reopening WebSocket, current state is ".concat(this.getState()));

      if (this.isActive()) {
        try {
          return this.close();
        } catch (error) {
          _logger["default"].log("Failed to reopen WebSocket", error);
        } finally {
          _logger["default"].log("Reopening WebSocket in ".concat(this.constructor.reopenDelay, "ms"));

          setTimeout(this.open, this.constructor.reopenDelay);
        }
      } else {
        return this.open();
      }
    }
  }, {
    key: "getProtocol",
    value: function getProtocol() {
      if (this.webSocket) {
        return this.webSocket.protocol;
      }
    }
  }, {
    key: "isOpen",
    value: function isOpen() {
      return this.isState("open");
    }
  }, {
    key: "isActive",
    value: function isActive() {
      return this.isState("open", "connecting");
    } // Private

  }, {
    key: "isProtocolSupported",
    value: function isProtocolSupported() {
      return indexOf.call(supportedProtocols, this.getProtocol()) >= 0;
    }
  }, {
    key: "isState",
    value: function isState() {
      for (var _len = arguments.length, states = new Array(_len), _key = 0; _key < _len; _key++) {
        states[_key] = arguments[_key];
      }

      return indexOf.call(states, this.getState()) >= 0;
    }
  }, {
    key: "getState",
    value: function getState() {
      if (this.webSocket) {
        for (var state in _adapters["default"].WebSocket) {
          if (_adapters["default"].WebSocket[state] === this.webSocket.readyState) {
            return state.toLowerCase();
          }
        }
      }

      return null;
    }
  }, {
    key: "installEventHandlers",
    value: function installEventHandlers() {
      for (var eventName in this.events) {
        var handler = this.events[eventName].bind(this);
        this.webSocket["on".concat(eventName)] = handler;
      }
    }
  }, {
    key: "uninstallEventHandlers",
    value: function uninstallEventHandlers() {
      for (var eventName in this.events) {
        this.webSocket["on".concat(eventName)] = function () {};
      }
    }
  }]);

  return Connection;
}();

Connection.reopenDelay = 500;
Connection.prototype.events = {
  message: function message(event) {
    if (!this.isProtocolSupported()) {
      return;
    }

    var _JSON$parse = JSON.parse(event.data),
        identifier = _JSON$parse.identifier,
        message = _JSON$parse.message,
        reason = _JSON$parse.reason,
        reconnect = _JSON$parse.reconnect,
        type = _JSON$parse.type;

    switch (type) {
      case message_types.welcome:
        this.monitor.recordConnect();
        return this.subscriptions.reload();

      case message_types.disconnect:
        _logger["default"].log("Disconnecting. Reason: ".concat(reason));

        return this.close({
          allowReconnect: reconnect
        });

      case message_types.ping:
        return this.monitor.recordPing();

      case message_types.confirmation:
        return this.subscriptions.notify(identifier, "connected");

      case message_types.rejection:
        return this.subscriptions.reject(identifier);

      default:
        return this.subscriptions.notify(identifier, "received", message);
    }
  },
  open: function open() {
    _logger["default"].log("WebSocket onopen event, using '".concat(this.getProtocol(), "' subprotocol"));

    this.disconnected = false;

    if (!this.isProtocolSupported()) {
      _logger["default"].log("Protocol is unsupported. Stopping monitor and disconnecting.");

      return this.close({
        allowReconnect: false
      });
    }
  },
  close: function close(event) {
    _logger["default"].log("WebSocket onclose event");

    if (this.disconnected) {
      return;
    }

    this.disconnected = true;
    this.monitor.recordDisconnect();
    return this.subscriptions.notifyAll("disconnected", {
      willAttemptReconnect: this.monitor.isRunning()
    });
  },
  error: function error() {
    _logger["default"].log("WebSocket onerror event");
  }
};
var _default = Connection;
exports["default"] = _default;

},{"./adapters":1,"./connection_monitor":3,"./internal":6,"./logger":7}],3:[function(require,module,exports){
"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports["default"] = void 0;

var _logger = _interopRequireDefault(require("./logger"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { "default": obj }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } }

function _createClass(Constructor, protoProps, staticProps) { if (protoProps) _defineProperties(Constructor.prototype, protoProps); if (staticProps) _defineProperties(Constructor, staticProps); return Constructor; }

// Responsible for ensuring the cable connection is in good health by validating the heartbeat pings sent from the server, and attempting
// revival reconnections if things go astray. Internal class, not intended for direct user manipulation.
var now = function now() {
  return new Date().getTime();
};

var secondsSince = function secondsSince(time) {
  return (now() - time) / 1000;
};

var clamp = function clamp(number, min, max) {
  return Math.max(min, Math.min(max, number));
};

var ConnectionMonitor =
/*#__PURE__*/
function () {
  function ConnectionMonitor(connection) {
    _classCallCheck(this, ConnectionMonitor);

    this.visibilityDidChange = this.visibilityDidChange.bind(this);
    this.connection = connection;
    this.reconnectAttempts = 0;
  }

  _createClass(ConnectionMonitor, [{
    key: "start",
    value: function start() {
      if (!this.isRunning()) {
        this.startedAt = now();
        delete this.stoppedAt;
        this.startPolling();
        addEventListener("visibilitychange", this.visibilityDidChange);

        _logger["default"].log("ConnectionMonitor started. pollInterval = ".concat(this.getPollInterval(), " ms"));
      }
    }
  }, {
    key: "stop",
    value: function stop() {
      if (this.isRunning()) {
        this.stoppedAt = now();
        this.stopPolling();
        removeEventListener("visibilitychange", this.visibilityDidChange);

        _logger["default"].log("ConnectionMonitor stopped");
      }
    }
  }, {
    key: "isRunning",
    value: function isRunning() {
      return this.startedAt && !this.stoppedAt;
    }
  }, {
    key: "recordPing",
    value: function recordPing() {
      this.pingedAt = now();
    }
  }, {
    key: "recordConnect",
    value: function recordConnect() {
      this.reconnectAttempts = 0;
      this.recordPing();
      delete this.disconnectedAt;

      _logger["default"].log("ConnectionMonitor recorded connect");
    }
  }, {
    key: "recordDisconnect",
    value: function recordDisconnect() {
      this.disconnectedAt = now();

      _logger["default"].log("ConnectionMonitor recorded disconnect");
    } // Private

  }, {
    key: "startPolling",
    value: function startPolling() {
      this.stopPolling();
      this.poll();
    }
  }, {
    key: "stopPolling",
    value: function stopPolling() {
      clearTimeout(this.pollTimeout);
    }
  }, {
    key: "poll",
    value: function poll() {
      var _this = this;

      this.pollTimeout = setTimeout(function () {
        _this.reconnectIfStale();

        _this.poll();
      }, this.getPollInterval());
    }
  }, {
    key: "getPollInterval",
    value: function getPollInterval() {
      var _this$constructor$pol = this.constructor.pollInterval,
          min = _this$constructor$pol.min,
          max = _this$constructor$pol.max,
          multiplier = _this$constructor$pol.multiplier;
      var interval = multiplier * Math.log(this.reconnectAttempts + 1);
      return Math.round(clamp(interval, min, max) * 1000);
    }
  }, {
    key: "reconnectIfStale",
    value: function reconnectIfStale() {
      if (this.connectionIsStale()) {
        _logger["default"].log("ConnectionMonitor detected stale connection. reconnectAttempts = ".concat(this.reconnectAttempts, ", pollInterval = ").concat(this.getPollInterval(), " ms, time disconnected = ").concat(secondsSince(this.disconnectedAt), " s, stale threshold = ").concat(this.constructor.staleThreshold, " s"));

        this.reconnectAttempts++;

        if (this.disconnectedRecently()) {
          _logger["default"].log("ConnectionMonitor skipping reopening recent disconnect");
        } else {
          _logger["default"].log("ConnectionMonitor reopening");

          this.connection.reopen();
        }
      }
    }
  }, {
    key: "connectionIsStale",
    value: function connectionIsStale() {
      return secondsSince(this.pingedAt ? this.pingedAt : this.startedAt) > this.constructor.staleThreshold;
    }
  }, {
    key: "disconnectedRecently",
    value: function disconnectedRecently() {
      return this.disconnectedAt && secondsSince(this.disconnectedAt) < this.constructor.staleThreshold;
    }
  }, {
    key: "visibilityDidChange",
    value: function visibilityDidChange() {
      var _this2 = this;

      if (document.visibilityState === "visible") {
        setTimeout(function () {
          if (_this2.connectionIsStale() || !_this2.connection.isOpen()) {
            _logger["default"].log("ConnectionMonitor reopening stale connection on visibilitychange. visbilityState = ".concat(document.visibilityState));

            _this2.connection.reopen();
          }
        }, 200);
      }
    }
  }]);

  return ConnectionMonitor;
}();

ConnectionMonitor.pollInterval = {
  min: 3,
  max: 30,
  multiplier: 5
};
ConnectionMonitor.staleThreshold = 6; // Server::Connections::BEAT_INTERVAL * 2 (missed two pings)

var _default = ConnectionMonitor;
exports["default"] = _default;

},{"./logger":7}],4:[function(require,module,exports){
"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.createWebSocketURL = createWebSocketURL;
exports["default"] = void 0;

var _connection = _interopRequireDefault(require("./connection"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { "default": obj }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } }

function _createClass(Constructor, protoProps, staticProps) { if (protoProps) _defineProperties(Constructor.prototype, protoProps); if (staticProps) _defineProperties(Constructor, staticProps); return Constructor; }

// The ActionCable.Consumer establishes the connection to a server-side Ruby Connection object. Once established,
// the ActionCable.ConnectionMonitor will ensure that its properly maintained through heartbeats and checking for stale updates.
// The Consumer instance is also the gateway to establishing subscriptions to desired channels through the #createSubscription
// method.
//
// The following example shows how this can be setup:
//
//   App = {}
//   App.cable = ActionCable.createConsumer("ws://example.com/accounts/1")
//
// For more details on how you'd configure an actual channel subscription, see ActionCable.Subscription.
//
// When a consumer is created, it automatically connects with the server.
//
// To disconnect from the server, call
//
//   App.cable.disconnect()
//
// and to restart the connection:
//
//   App.cable.connect()
var Consumer =
/*#__PURE__*/
function () {
  function Consumer(url) {
    _classCallCheck(this, Consumer);

    this._url = url;
    this.connection = new _connection["default"](this);
  }

  _createClass(Consumer, [{
    key: "send",
    value: function send(data) {
      return this.connection.send(data);
    }
  }, {
    key: "connect",
    value: function connect() {
      return this.connection.open();
    }
  }, {
    key: "disconnect",
    value: function disconnect() {
      return this.connection.close({
        allowReconnect: false
      });
    }
  }, {
    key: "ensureActiveConnection",
    value: function ensureActiveConnection() {
      if (!this.connection.isActive()) {
        return this.connection.open();
      }
    }
  }, {
    key: "url",
    get: function get() {
      return createWebSocketURL(this._url);
    }
  }]);

  return Consumer;
}();

exports["default"] = Consumer;

function createWebSocketURL(url) {
  if (typeof url === "function") {
    url = url();
  }

  if (url && !/^wss?:/i.test(url)) {
    var a = document.createElement("a");
    a.href = url; // Fix populating Location properties in IE. Otherwise, protocol will be blank.

    a.href = a.href;
    a.protocol = a.protocol.replace("http", "ws");
    return a.href;
  } else {
    return url;
  }
}

},{"./connection":2}],5:[function(require,module,exports){
"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.createConsumer = createConsumer;
Object.defineProperty(exports, "Connection", {
  enumerable: true,
  get: function get() {
    return _connection["default"];
  }
});
Object.defineProperty(exports, "ConnectionMonitor", {
  enumerable: true,
  get: function get() {
    return _connection_monitor["default"];
  }
});
Object.defineProperty(exports, "Consumer", {
  enumerable: true,
  get: function get() {
    return _consumer["default"];
  }
});
Object.defineProperty(exports, "createWebSocketURL", {
  enumerable: true,
  get: function get() {
    return _consumer.createWebSocketURL;
  }
});
Object.defineProperty(exports, "INTERNAL", {
  enumerable: true,
  get: function get() {
    return _internal["default"];
  }
});
Object.defineProperty(exports, "adapters", {
  enumerable: true,
  get: function get() {
    return _adapters["default"];
  }
});
Object.defineProperty(exports, "logger", {
  enumerable: true,
  get: function get() {
    return _logger["default"];
  }
});

var _connection = _interopRequireDefault(require("./connection"));

var _connection_monitor = _interopRequireDefault(require("./connection_monitor"));

var _consumer = _interopRequireWildcard(require("./consumer"));

var _internal = _interopRequireDefault(require("./internal"));

var _adapters = _interopRequireDefault(require("./adapters"));

var _logger = _interopRequireDefault(require("./logger"));

function _interopRequireWildcard(obj) { if (obj && obj.__esModule) { return obj; } else { var newObj = {}; if (obj != null) { for (var key in obj) { if (Object.prototype.hasOwnProperty.call(obj, key)) { var desc = Object.defineProperty && Object.getOwnPropertyDescriptor ? Object.getOwnPropertyDescriptor(obj, key) : {}; if (desc.get || desc.set) { Object.defineProperty(newObj, key, desc); } else { newObj[key] = obj[key]; } } } } newObj["default"] = obj; return newObj; } }

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { "default": obj }; }

function createConsumer() {
  var url = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : _internal["default"].default_mount_path;
  return new _consumer["default"](url);
}

consumer = createConsumer();
console.log(consumer);

},{"./adapters":1,"./connection":2,"./connection_monitor":3,"./consumer":4,"./internal":6,"./logger":7}],6:[function(require,module,exports){
"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports["default"] = void 0;
var _default = {
  message_types: {
    welcome: "welcome",
    ping: "ping",
    confirmation: "confirm_subscription",
    rejection: "reject_subscription"
  },
  default_mount_path: "/ws",
  protocols: ["actioncable-v1-json", "actioncable-unsupported"]
};
exports["default"] = _default;

},{}],7:[function(require,module,exports){
"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports["default"] = void 0;

var _adapters = _interopRequireDefault(require("./adapters"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { "default": obj }; }

var _default = {
  log: function log() {
    if (this.enabled) {
      var _adapters$logger;

      for (var _len = arguments.length, messages = new Array(_len), _key = 0; _key < _len; _key++) {
        messages[_key] = arguments[_key];
      }

      messages.push(Date.now());

      (_adapters$logger = _adapters["default"].logger).log.apply(_adapters$logger, ["[Anychat]"].concat(messages));
    }
  }
};
exports["default"] = _default;

},{"./adapters":1}]},{},[1,2,3,4,5,6,7]);
