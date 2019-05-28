var anychat=function(t){var e={};function n(o){if(e[o])return e[o].exports;var i=e[o]={i:o,l:!1,exports:{}};return t[o].call(i.exports,i,i.exports,n),i.l=!0,i.exports}return n.m=t,n.c=e,n.d=function(t,e,o){n.o(t,e)||Object.defineProperty(t,e,{enumerable:!0,get:o})},n.r=function(t){"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(t,"__esModule",{value:!0})},n.t=function(t,e){if(1&e&&(t=n(t)),8&e)return t;if(4&e&&"object"==typeof t&&t&&t.__esModule)return t;var o=Object.create(null);if(n.r(o),Object.defineProperty(o,"default",{enumerable:!0,value:t}),2&e&&"string"!=typeof t)for(var i in t)n.d(o,i,function(e){return t[e]}.bind(null,i));return o},n.n=function(t){var e=t&&t.__esModule?function(){return t.default}:function(){return t};return n.d(e,"a",e),e},n.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},n.p="",n(n.s=0)}([function(t,e,n){"use strict";n.r(e);var o={logger:self.console,WebSocket:self.WebSocket},i={log(...t){this.enabled&&(t.push(Date.now()),o.logger.log("[Anychat]",...t))}};const s=()=>(new Date).getTime(),r=t=>(s()-t)/1e3,c=(t,e,n)=>Math.max(e,Math.min(n,t));class l{constructor(t){this.visibilityDidChange=this.visibilityDidChange.bind(this),this.connection=t,this.reconnectAttempts=0}start(){this.isRunning()||(this.startedAt=s(),delete this.stoppedAt,this.startPolling(),addEventListener("visibilitychange",this.visibilityDidChange),i.log(`ConnectionMonitor started. pollInterval = ${this.getPollInterval()} ms`))}stop(){this.isRunning()&&(this.stoppedAt=s(),this.stopPolling(),removeEventListener("visibilitychange",this.visibilityDidChange),i.log("ConnectionMonitor stopped"))}isRunning(){return this.startedAt&&!this.stoppedAt}recordPing(){this.pingedAt=s()}recordConnect(){this.reconnectAttempts=0,this.recordPing(),delete this.disconnectedAt,i.log("ConnectionMonitor recorded connect")}recordDisconnect(){this.disconnectedAt=s(),i.log("ConnectionMonitor recorded disconnect")}startPolling(){this.stopPolling(),this.poll()}stopPolling(){clearTimeout(this.pollTimeout)}poll(){this.pollTimeout=setTimeout(()=>{this.reconnectIfStale(),this.poll()},this.getPollInterval())}getPollInterval(){const{min:t,max:e,multiplier:n}=this.constructor.pollInterval,o=n*Math.log(this.reconnectAttempts+1);return Math.round(1e3*c(o,t,e))}reconnectIfStale(){this.connectionIsStale()&&(i.log(`ConnectionMonitor detected stale connection. reconnectAttempts = ${this.reconnectAttempts}, pollInterval = ${this.getPollInterval()} ms, time disconnected = ${r(this.disconnectedAt)} s, stale threshold = ${this.constructor.staleThreshold} s`),this.reconnectAttempts++,this.disconnectedRecently()?i.log("ConnectionMonitor skipping reopening recent disconnect"):(i.log("ConnectionMonitor reopening"),this.connection.reopen()))}connectionIsStale(){return r(this.pingedAt?this.pingedAt:this.startedAt)>this.constructor.staleThreshold}disconnectedRecently(){return this.disconnectedAt&&r(this.disconnectedAt)<this.constructor.staleThreshold}visibilityDidChange(){"visible"===document.visibilityState&&setTimeout(()=>{!this.connectionIsStale()&&this.connection.isOpen()||(i.log(`ConnectionMonitor reopening stale connection on visibilitychange. visbilityState = ${document.visibilityState}`),this.connection.reopen())},200)}}l.pollInterval={min:3,max:30,multiplier:5},l.staleThreshold=6;var a=l,h={message_types:{welcome:0,ping:4,ack:11,single_chat:101,room_chat:102},default_mount_path:"/anychat",protocols:["anychat-v1-json","anychat-unsupported"]};const{message_types:u,protocols:d}=h,p=d.slice(0,d.length-1),g=[].indexOf;class f{constructor(t){this.open=this.open.bind(this),this.consumer=t,this.monitor=new a(this),this.disconnected=!0}send(t){return!!this.isOpen()&&(this.webSocket.send(JSON.stringify(t)),!0)}open(){return this.isActive()?(i.log(`Attempted to open WebSocket, but existing socket is ${this.getState()}`),!1):(i.log(`Opening WebSocket, current state is ${this.getState()}, subprotocols: ${d}`),this.webSocket&&this.uninstallEventHandlers(),this.webSocket=new o.WebSocket(this.consumer.url,d),this.installEventHandlers(),this.monitor.start(),!0)}close({allowReconnect:t}={allowReconnect:!0}){if(t||this.monitor.stop(),this.isActive())return this.webSocket.close()}reopen(){if(i.log(`Reopening WebSocket, current state is ${this.getState()}`),!this.isActive())return this.open();try{return this.close()}catch(t){i.log("Failed to reopen WebSocket",t)}finally{i.log(`Reopening WebSocket in ${this.constructor.reopenDelay}ms`),setTimeout(this.open,this.constructor.reopenDelay)}}getProtocol(){if(this.webSocket)return this.webSocket.protocol}isOpen(){return this.isState("open")}isActive(){return this.isState("open","connecting")}isProtocolSupported(){return g.call(p,this.getProtocol())>=0}isState(...t){return g.call(t,this.getState())>=0}getState(){if(this.webSocket)for(let t in o.WebSocket)if(o.WebSocket[t]===this.webSocket.readyState)return t.toLowerCase();return null}installEventHandlers(){for(let t in this.events){const e=this.events[t].bind(this);this.webSocket[`on${t}`]=e}}uninstallEventHandlers(){for(let t in this.events)this.webSocket[`on${t}`]=function(){}}}f.reopenDelay=500,f.prototype.events={message(t){if(!this.isProtocolSupported())return;i.log(t.data);const{cmd:e,ack:n,data:o}=JSON.parse(t.data);switch(e){case u.welcome:return this.monitor.recordConnect();case u.disconnect:return i.log(`Disconnecting. Reason: ${o.reason}`),this.close({allowReconnect:o.reconnect});case u.ping:return this.monitor.recordPing();case u.ack:case u.single_chat:case u.room_chat:default:return}},open(){if(i.log(`WebSocket onopen event, using '${this.getProtocol()}' subprotocol`),this.disconnected=!1,!this.isProtocolSupported())return i.log("Protocol is unsupported. Stopping monitor and disconnecting."),this.close({allowReconnect:!1})},close(t){i.log("WebSocket onclose event"),this.disconnected||(this.disconnected=!0,this.monitor.recordDisconnect())},error(){i.log("WebSocket onerror event")}};var m=f;class b{constructor(t){this._url=t,this.connection=new m(this)}get url(){return v(this._url)}send(t){return this.connection.send(t)}connect(){return this.connection.open()}disconnect(){return this.connection.close({allowReconnect:!1})}ensureActiveConnection(){if(!this.connection.isActive())return this.connection.open()}}function v(t){if("function"==typeof t&&(t=t()),t&&!/^wss?:/i.test(t)){const e=document.createElement("a");return e.href=t,e.href=e.href,e.protocol=e.protocol.replace("http","ws"),e.href}return t}function S(t=h.default_mount_path){return new b(t)}n.d(e,"createConsumer",function(){return S}),n.d(e,"Connection",function(){return m}),n.d(e,"ConnectionMonitor",function(){return a}),n.d(e,"Consumer",function(){return b}),n.d(e,"INTERNAL",function(){return h}),n.d(e,"adapters",function(){return o}),n.d(e,"createWebSocketURL",function(){return v}),n.d(e,"logger",function(){return i})}]);