# gRPC client 实现

我们要实现的 client 应该有从远端获取缓存的能力，通过定义 Fetcher 接口，然后自定义 client 类型实现该接口即可。

实现原理：我们的 grpc client 最终是通过 client stub 发出的请求，我们先获取一个 etcd client 对象，然后通过服务注册发现获取与服务的连接（grpc 通道），接着使用这个连接对象就可以新建一个 grpc client stub了，然后再申请一个超时自动取消的上下文对象；最后，使用上下文并填充 protocol buffer 请求参数去调用 client stud 的方法（这个方法在服务端也实现了）获取响应即可。

具体代码见：distributekv/server.go