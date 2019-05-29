# new chat

* [chat references](https://github.com/FlowerWrong/erlim/blob/master/api/chat.md#chat-api)

## 技术栈

* [postgresql 11.2+](https://www.postgresql.org/)
* [rails 6](https://rubyonrails.org/): db and test only
* [gorilla/websocket](https://github.com/gorilla/websocket)
* [go-micro](https://micro.mu/)
* [opencensus](https://github.com/census-instrumentation/opencensus-go): A stats collection and distributed tracing framework
* [uber zap](https://github.com/uber-go/zap) + [lumberjack](https://github.com/natefinch/lumberjack)
* [otto](https://github.com/robertkrimen/otto): a javascript interpreter
* [blevesearch/bleve](https://github.com/blevesearch/bleve)
* [xormplus](https://github.com/xormplus/xorm)
* [viper](https://github.com/spf13/viper)
* [gin](https://github.com/gin-gonic/gin)
* [nats](https://nats.io/documentation/)
* [jaeger](https://github.com/jaegertracing/jaeger)
* [sonyflake](https://github.com/sony/sonyflake): distributed unique ID generator

## Design

用户列表统一存储在redis，key为user.UUID，value为node host:port，即所在服务器实例，走rpc

## TODO

* [ ] cluster
* [ ] pub/sub based on redis
* [ ] grpc
* [ ] session reconnectable
* [ ] registration center: consul
* [ ] configuration center
* [ ] hot restart
