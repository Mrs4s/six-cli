# 6盘 CLI工具操作指南

> 建议在 Unix-like 系统下使用本工具的shell模式
>
> 理论上不支持win10 TH2以下的原生命令行, 因为[在微软Windows 10更新TH2之前，Windows操作系统的Win32控制台是不支持ANSI转义序列的](https://zh.wikipedia.org/zh-hans/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97)

## 安装

````go
go get github.com/Mrs4s/six-cli
````

> 也可从 [Releases](https://github.com/Mrs4s/six-cli) 直接下载二进制 

##  使用方法

* 待完善，目前尚未开发完成。

## 命令列表

### 登录

````shell
guest@six-pan:/$ login
请输入用户名: mrs4s
请输入密码:           # 密码不会回显，输完直接回车即可
登录完成, 欢迎: mrs4s
````

### 切换工作目录

````shell
mrs4s@six-pan:/$ cd target  # 进入子目录
mrs4s@six-pan:/target$ cd /target/sub/test  # 以完整路径切换目录
mrs4s@six-pan:/target/sub/test$ cd ..  # 返回上一层
mrs4s@six-pan:/target/sub$ cd ../../  # 向上返回N层
````

### 列出文件

````shell
mrs4s@six-pan:/$ ls  # 列出当前目录所有对象
dir1	dir2	file1	file2
mrs4s@six-pan:/$ ls -d  # 按文件夹过滤
dir1	dir2
mrs4s@six-pan:/$ ls /target  # 列出目标目录所有对象
dir1	file1
mrs4s@six-pan:/$ ls /target -R  # 遍历列出子目录对象 (鉴于负载考虑不递归子目录)
.:
dir1	file1
./dir1:
test1	test2
mrs4s@six-pan:/$ ls -a  # 输出文件详细信息
序号  创建时间             文件大小  文件名
0     2019-01-01 00:00:00 100.00GB  dir1
...
````

.. 未完待续