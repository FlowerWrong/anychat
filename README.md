# new chat

* [chat references](https://github.com/FlowerWrong/erlim/blob/master/api/chat.md#chat-api)

## 技术栈

* [postgresql 11.2+](https://www.postgresql.org/)
* [rails 6](https://rubyonrails.org/)
* [gorilla/websocket](https://github.com/gorilla/websocket)
* [go-micro](https://micro.mu/)
* [opencensus](https://github.com/census-instrumentation/opencensus-go): A stats collection and distributed tracing framework
* [uber zap](https://github.com/uber-go/zap) + [lumberjack](https://github.com/natefinch/lumberjack)
* [yuin/gopher-lua](https://github.com/yuin/gopher-lua)
* [blevesearch/bleve](https://github.com/blevesearch/bleve)
* [xorm](https://github.com/xormplus/xorm)
* [viper](https://github.com/spf13/viper)
* [ant.design pro](http://pro.ant.design/index-cn/)
* [gin](https://github.com/gin-gonic/gin)
* [nats](https://nats.io/documentation/)
* [jaeger](https://github.com/jaegertracing/jaeger)

## 竞品

* [udesk](http://www.udesk.cn/)
* [美洽](https://meiqia.com/)
* [live.chat](https://www.livechatinc.com/)
* [tawk.to](https://www.tawk.to/)

## Features

* 自定义界面
* js SDK
* iOS/android SDK
* react SDK
* react native SDK

## Design

用户列表统一存储在redis，key为user.UUID，value为node host:port，即所在服务器实例，走rpc
