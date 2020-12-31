# TCP & HTTP connections test tool
## get started
### start server
help info:
```
Usage of ./server:
  -http-listen-address string
    	address http server listen on (default "0.0.0.0:8080")
  -http-response-body-size int
    	http response body size (default 64)
  -http-response-code int
    	http response status code (default 200)
  -protocol string
    	start server protocol (default "fake")
  -report-interval duration
    	interval of report (default 1s)
  -tcp-address string
    	address tcp server listen (default "0.0.0.0:9090")
  -tcp-pkg-size int
    	tcp max package size (default 4096)
```

*Example 1*: start tcp server:
```
$ go run server.go --protocol tcp
2020/12/31 17:08:33 start server:  tcp
2020/12/31 17:08:33 tcp  listen on  0.0.0.0:9090
2020/12/31 17:08:33 | connections     | packages        | bytes           |
2020/12/31 17:08:34 | 0               | 0               | 0               |
2020/12/31 17:08:35 | 0               | 0               | 0               |
2020/12/31 17:08:36 | 0               | 0               | 0               |
2020/12/31 17:08:37 | 0               | 0               | 0               |
2020/12/31 17:08:38 | 0               | 0               | 0               |
2020/12/31 17:08:39 | 0               | 0               | 0               |
2020/12/31 17:08:40 | 5               | 5               | 500             |
2020/12/31 17:08:41 | 5               | 10              | 1000            |
2020/12/31 17:08:42 | 5               | 15              | 1500            |
2020/12/31 17:08:43 | 5               | 20              | 2000            |
2020/12/31 17:08:44 | 0               | 21              | 2100            |
2020/12/31 17:08:45 | 0               | 21              | 2100            |
2020/12/31 17:08:46 | 0               | 21              | 2100            |
```

*Example 2*: start http server:
```
$ go run server.go --protocol http
2020/12/31 17:09:44 start server:  http
2020/12/31 17:09:44 http  listen on  0.0.0.0:8080
2020/12/31 17:09:44 | requests        | send bytes      | receive bytes   |
2020/12/31 17:09:45 | 0               | 0               | 0               |
2020/12/31 17:09:46 | 0               | 0               | 0               |
2020/12/31 17:09:47 | 0               | 0               | 0               |
2020/12/31 17:09:48 | 0               | 0               | 0               |
2020/12/31 17:09:49 | 0               | 0               | 0               |
2020/12/31 17:09:50 | 2               | 128             | 2000            |
2020/12/31 17:09:51 | 2               | 128             | 2000            |
2020/12/31 17:09:52 | 2               | 128             | 2000            |
2020/12/31 17:09:53 | 4               | 256             | 4000            |
```

### start client
help info:
```
Usage of ./client:
  -http-parallel int
    	http connect parallel (default 1)
  -http-request-body-size int
    	http request body size (default 64)
  -http-request-interval duration
    	http request interval (default 1s)
  -http-request-method string
    	http request method (default "POST")
  -http-request-timeout duration
    	http request timeout (default 1s)
  -http-request-url string
    	http request url (default "http://127.0.0.1:8080/")
  -protocol string
    	start clientset protocol (default "fake")
  -report-interval duration
    	interval of report (default 1s)
  -tcp-address string
    	tcp server address connect to (default ":9090")
  -tcp-package-interval duration
    	interval of packages (default 1s)
  -tcp-package-size int
    	package send to server in bytes (default 64)
  -tcp-parallel int
    	how many tcp connections to server (default 1)
```

*Example 1*: start tcp client
```
$ go run client.go --protocol tcp --tcp-package-size 100 --tcp-parallel 5 --report-interval 1s
2020/12/31 17:12:15 start clientset:  tcp
2020/12/31 17:12:15 | connections     | packages        | bytes           |
2020/12/31 17:12:16 | 5               | 6               | 600             |
2020/12/31 17:12:17 | 5               | 13              | 1300            |
2020/12/31 17:12:18 | 5               | 15              | 1500            |
2020/12/31 17:12:19 | 5               | 23              | 2300            |
2020/12/31 17:12:20 | 5               | 25              | 2500            |
^C2020/12/31 17:12:21 clientset stop:  interrupt
```

*Example 2*: start http client
```
$ go run client.go --protocol http --http-parallel 2 --http-request-body-size 1000 --http-request-method POST --http-request-interval 3s --report-interval 1s
2020/12/31 17:13:01 start clientset:  http
2020/12/31 17:13:01 | connections     | requests        | send bytes      | receive bytes   |
2020/12/31 17:13:02 | 2               | 2               | 2000            | 128             |
2020/12/31 17:13:03 | 2               | 2               | 2000            | 128             |
2020/12/31 17:13:04 | 2               | 2               | 2000            | 128             |
2020/12/31 17:13:05 | 2               | 4               | 4000            | 256             |
2020/12/31 17:13:06 | 2               | 4               | 4000            | 256             |
2020/12/31 17:13:07 | 2               | 4               | 4000            | 256             |
2020/12/31 17:13:08 | 2               | 6               | 6000            | 384             |
2020/12/31 17:13:09 | 2               | 6               | 6000            | 384             |
^C2020/12/31 17:13:10 clientset stop:  interrupt
```

## more protocol
if you want more protocol supported, please add your implementation of Server and ClientSet like TCPServer and TCPClientSet