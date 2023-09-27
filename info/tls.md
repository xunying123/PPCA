# TLS 劫持

对于 HTTPS 等 TLS 连接, 一个简单的代理服务器无法知道加密连接中所传输的内容, 毕竟这是 TLS 的目的之一 (防止连接上的中间人窃听连接内容)。TLS 通过 CA 签发的证书来验证服务器的真实性, 此后与服务器建立加密连接, 密钥对中间方不可见。

```
+--------------+   +-------------------+   +-------------+
| User         |   | Proxy             |   | Server      |
| +----------+ |   | Encrypted traffic |   | +---------+ |
| | Encrypt -+-+->-+--------->---------+->-+-> Decrypt | |
| +----------+ |   |                   |   | +---------+ |
+--------------+   +-------------------+   +-------------+
```

但是, 你可以在代理时, 假装自己是服务器, 建立两个连接:


```
+--------------+   +----------------------------+   +-------------+
| User         |   | Proxy                      |   | Server      |
| +----------+ |   | +----------+  +----------+ |   | +---------+ |
| | Encrypt -+-+->-+-> Decrypt -+--> Encrypt -+-+->-+-> Decrypt | |
| +----------+ |   | +----------+  +----------+ |   | +---------+ |
+--------------+   +----------------------------+   +-------------+
```

当然, 在一般情况下, 代理服务器并不能获得真正服务器所具有的、由权威机构签发的 TLS 证书。但是, 当用户和代理服务器由同一人操作的时候, 用户可以无条件信任代理服务器; 因此, 代理服务器可以自己创建一个根证书来实时签发代理时所使用的 TLS 证书。

推荐调用 OpenSSL 等 TLS 库完成 TLS 协议的相关工作。