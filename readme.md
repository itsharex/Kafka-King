<p align="center">
  <img src="docs/snap/icon.ico" alt="图片标题">
</p>
<h1 align="center">Kafka King </h1>

<h4 align="center"><strong>简体中文</strong> | <a href="https://github.com/Bronya0/Kafka-King/blob/wails/readme-en.md">English</a></h4>

<div align="center">

![License](https://img.shields.io/github/license/Bronya0/Kafka-King)
![GitHub release](https://img.shields.io/github/release/Bronya0/Kafka-King)
![GitHub All Releases](https://img.shields.io/github/downloads/Bronya0/Kafka-King/total)
![GitHub stars](https://img.shields.io/github/stars/Bronya0/Kafka-King)
![GitHub forks](https://img.shields.io/github/forks/Bronya0/Kafka-King)

<h3 align="center">一个现代、实用的kafka GUI客户端 </h3>

<strong></strong>
</div>

让kafka更好用，make kafka great again!

本项目是一个kafka GUI客户端，支持各个系统。点个star支持作者辛苦开源吧 谢谢❤❤
加群和作者一起交流： <a target="_blank" href="https://qm.qq.com/cgi-bin/qm/qr?k=pDqlVFyLMYEEw8DPJlRSBN27lF8qHV2v&jump_from=webapi&authKey=Wle/K0ARM1YQWlpn6vvfiZuMedy2tT9BI73mUvXVvCuktvi0fNfmNR19Jhyrf2Nz">研发技术交流群：964440643</a>

> 同款好用elasticsearch客户端 `ES-King`，可以一起收藏下：https://github.com/Bronya0/ES-King





# Kafka-King功能清单
- [x] 查看集群节点列表，支持动态配置broker、topic的配置项！
- [x] 支持消费者客户端，按照指定的group进行指定topic、size、timeout的消费，以表格的形式展示消息的各个维度信息！
- [x] 支持PLAIN、SSL、SASL、kerberos、sasl_plaintext等等
- [x] 创建主题（支持批量）、删除主题，指定副本、分区
- [x] 支持根据消费者组统计每个topic的消息总量、提交总量、积压量
- [x] 支持查看topic的分区的详细信息（offset），并支持添加额外的分区
- [x] 支持模拟生产者，批量发送消息，指定headers、分区
- [x] topic、分区健康检查（完成）
- [x] 支持查看消费者组、消费者
- [x] offset巡检报表

# 下载
右侧下载，或者点[下载地址](https://github.com/Bronya0/Kafka-King/releases)，展开【Assets】，选择自己的平台下载，支持windows、macos、linux。

> **使用前请检查kafka集群配置的`advertised.listeners`，如果未配置，或配置的是域名，那么在King中填写连接地址时，请提前在本机电脑的hosts文件中添加对应域名解析，否则会因为无法解析域名而无法连接**
> 
> **如果你的连接需要SSL，那么勾选TLS并勾选忽略验证（有证书的话就下下来，开启tls验证，填入证书路径）。**
> 
> **对于SASL机制用户需要勾选开启SASL，并选择SASL协议（通常是plain），并填入用户名密码。**

如果遇到webview2运行时错误，从微软官网下载重装运行时即可：https://developer.microsoft.com/zh-cn/microsoft-edge/webview2

# 功能截图
offset巡检，v0.33版本上线，最直观的方式查看消息积压情况
![](docs/snap/img_5.png)
topic列表，各种操作
![](docs/snap/img.png)
查看消息
![](docs/snap/img_3.png)


# 捐赠
有条件可以请作者喝杯咖啡，支持项目发展，感谢💕

![image](https://github.com/user-attachments/assets/da6d46da-4e24-41e3-843d-495c6cd32065)

# 构建
只有要研究源码才需要手动构建

cd app

wails dev

## QQ交流群
<a target="_blank" href="https://qm.qq.com/cgi-bin/qm/qr?k=pDqlVFyLMYEEw8DPJlRSBN27lF8qHV2v&jump_from=webapi&authKey=Wle/K0ARM1YQWlpn6vvfiZuMedy2tT9BI73mUvXVvCuktvi0fNfmNR19Jhyrf2Nz">KingTool研发技术交流群：964440643</a>

![](assets/qq.jpg)


# Star星星
[![Stargazers over time](https://starchart.cc/Bronya0/Kafka-King.svg)](https://starchart.cc/Bronya0/Kafka-King)


# License
Apache-2.0 license

# 感谢
- wails：https://wails.io/docs/gettingstarted/installation
- naive ui：https://www.naiveui.com/
- franz-go：https://github.com/twmb/franz-go/
- xicons：https://xicons.org/#/
