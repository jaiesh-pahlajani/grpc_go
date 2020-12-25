# About this repo

###### Types -

1. Unary - Classic request / response on same TCP connection

2. Server Streaming - Client sends one request, receives a stream of responses on same TCP connection

3. Client Streaming - Client sends a stream of requests, receives a response on same TCP connection

4. Bidirectional - Client and serves send a stream of request and responses on same TCP connection

**The streaming capabilities are powered HTTP/2.0**

###### Notes -

Async by default and do not block threads on request and can serve millions of request in parallel

Grpc can perform client side load balancing

Each language will provide an API to load gRPC with required SSL certificates 

Authentication using interceptors

###### REST vs gRPC 

![Alt text](grpc_vs_rest.png?raw=true "gRPC vs Rest")

###### Error Handling

https://grpc.io/docs/guides/error/

https://avi.im/grpc-errors/

###### Deadlines

Deadlines allow gRPC clients to specify how long they are willing for an RPC to complete before RPC is terminated with `DEADLINE_EXCEEDED`

https://grpc.io/blog/deadlines/

###### SSL

In production grpc calls should be running with encryption enabled.
This is done by generating SSL certificates.
SSL allows comms to be secure end to end.
When we communicate over internet data is visible to all servers that transfers your packet.
Any router can intercept that packet.
Using ssl all messages are encrypted.
Two ways of using SSL(grpc can use both):
1. 1 way verification eg: browser => webserver (ENCRYPTION)
2. 2 way verification eg: SSL auth (AUTHENTICATION)

https://grpc.io/docs/guides/auth/

###### Evans

For calling server without setting up clients -
https://github.com/ktr0731/evans

###### Other Links

https://docs.microsoft.com/en-us/aspnet/core/grpc/comparison?view=aspnetcore-3.0

https://www.slideshare.net/borisovalex/grpc-vs-rest-let-the-battle-begin-81800634

https://husobee.github.io/golang/rest/grpc/2016/05/28/golang-rest-v-grpc.html

https://imagekit.io/demo/http2-vs-http1

https://grpc.io/

https://http2.github.io/