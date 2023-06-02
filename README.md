# E-WAF

一个高性能且简单的 WAF

TODO

* 阻断 SlowHTTP CC 攻击
* 支持 IP 白名单和黑名单功能
* 支持 User-Agent 的过滤访问
* 日志转储至 ElasticsSearch
* 良好的 Metrics 监控 WAF 本身状态
* 反向代理 (TCP / HTTP)
* 阻断常见的 SQL 注入攻击
* 阻断常见的 XSS 注入攻击
* 简易临时操作面板
* 支持 Reuseport 机制 (需要内核 >= 3.9)

感谢依赖和参考库
* WireGuard
* Logrus
* Echo
