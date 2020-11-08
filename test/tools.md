# 网络性能测试工具汇总
网络是指设备与设备之间的连接，网络性能测试往往是针对网络链路中的某中间设备进行的测试，例如：防火墙、网关等，本文整理常见的网络测试工具，以供参考。因为被测试对象只是中间设备，不能对网络数据包直接做出回应，所以一般在网络的一端需要一个Server，另一端需要一个或者多个Client，辅助进行测试

我们知道，网络性能测试主要考虑以下几个部分：
1. 被测试对象。被测试的网络设备是什么，常见的比如：网关、防火墙等中间设备
2. 测试指标。测试需要观测的度量指标，例如：最大连接数、每秒请求数、丢包率、时延等
3. 测试方法。如何搭建环境，如何测试，如何观测指标等。一般针对具体的测试指标，梳理出影响该指标的变量，每次变化一个变量，多次测试来观测指标表现情况
4. 测试结果。网络性能测试结果一般会说：在xCxG等资源消耗下可以达到xxx性能

下面就整理Linux环境中常见的网络性能测试工具和各工具的重点测试指标：

| 工具名称 | 充当角色 | 测试指标 | 说明 | 配合工具 | 特殊限制 |
| --- | --- | --- | --- | --- | --- | --- |
| [iperf3](https://iperf.fr/) | Server、Client | 带宽 | 既可以作为TCP/SCTP/UDP Server，又可以作为TCP/SCTP/UDP Client | 无论作为Server还是作为Client，一般都只跟iperf3配合测试 | 最大并发连接数128 |
| [nginx](http://nginx.org/en/docs/) | Server | - | 作为HTTP Server，处理Client的请求，以此来测试Server跟Client中间被测设备的性能 | ab、wrk | - |
| [ab](http://httpd.apache.org/docs/2.4/programs/ab.html) | Client | RPS等 | Apache Benchmark的缩写，用于测试HTTP相关性能 | httpd、nginx等HTTP Server | - |
| [mtr](http://www.bitwizard.nl/mtr/) | Client | 测试ICMP/TCP/UDP协议丢包率和时延 | Client | ICMP：不需要Server<br>TCP：iperf3等<br>UDP：iperf3等 | - |
| [wrk](https://github.com/wg/wrk) | Client | HTTP性能测试工具 | httpd、nginx等HTTP Server | - | - |
| [netserver](http://www.netperf.org/) | Server | 重点是TCP/UDP批量数据传输请求和响应性能和Berkeley Sockets接口，还可以将DLPI、UDS、ipv6的特性编译进来 | netperf | - | - |
| [netperf](http://www.netperf.org/) | Client | 参考netserver | netserver | - | - |