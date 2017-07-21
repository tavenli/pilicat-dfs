
# pilicat-dfs
霹雳猫-分布式文件系统

一种可以将网站图片或上传的文件，进行分布式存放的服务，可自动复制到多台物理机器，可满足高可用和负载均衡

# 已编译好的程序包

[http://git.oschina.net/tavenli/pilicat-dfs/releases](http://git.oschina.net/tavenli/pilicat-dfs/releases)


# 功能介绍

> 支持多个dfs-node

> 支持对上传的文件自动复制到多个dfs-node

> 支持 RestfulAPI 接口，方便各种语言调用

> 可以对dfs-node进行分组


# 为什么要使用pilicat-dfs

- 可使每个应用系统上传的附件，不存放在应用系统目录下，即使上传了图片木马，也执行不了，因为dfs是纯静态

- 可同时满足多个应用系统上传附件的需求，将多个应用上传的附件存放在统一的一组服务器，方便管理

- 可以更好的实现CDN，加速对静态资源的访问，因为这里都是纯静态

- 同一份文件，自动分发到多台物理节点，并支持通过 lvs 或 nginx 实现负载均衡，实现高可用和物理容灾

- 还有很多优点，就不全部阐述了...


# 功能使用
1. dfs-node当做单节点使用，非集群高可用方式
- 启动dfs-node
```
cd dfs-node_linux64_v1.0.0
./start.sh
```
- 上传文件测试
```
curl -X POST -F file=@/app/test.jpg http://127.0.0.1:8800/api/file
```
- 服务返回信息
```
{
	"Code": 0,
	"Msg": "success",
	"Data": {
		"FileUrlPath": "/file/2017/07/21/5e30cf328e44824ece5ddc52b629b73c.jpg",
		"OrgFileName": "test.jpg",
		"PubUrl": "http://dsf.hicode.top/file/2017/07/21/5e30cf328e44824ece5ddc52b629b73c.jpg"
	}
}
```

# dfs-node配置文件说明

```
node.name = "dfs-node-1"    //节点名称，在同一个center中唯一
node.public.addr = "0.0.0.0:8700"    //用于对外访问端口，主要供web访问上传后的文件，如果是80端口直接对外，可直接绑到80端口上
node.api.addr = "0.0.0.0:8800"    //用于内网接口，文件上传、覆盖、删除等操作，都通过该端口，通常绑定在内网IP

dfs.center = "192.168.1.200:8000"    //集群高可用服务所在位置，用于自动注册dfs-node
dfs.public.url = "http://dsf.hicode.top:8700"    //用于对外访问的域名
```

后续的功能将根据进度，陆续发布，敬请期待...









