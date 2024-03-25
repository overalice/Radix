<div align="center">
  <a href=#>
    <img src="http://122.51.217.174:8080/radix.png" alt="Logo" width="80" height="80">
  </a>
</div>

<h3 align="center">Radix</h3>

<p align="center">
让你快速开始构建项目的后端接口！
<br />
<a href="https://github.com/overalice/Radix"><strong>浏览文档 »</strong></a>
<br />
<br />
<a href="https://github.com/overalice/Radix">查看 Demo</a>
·
<a href="https://github.com/overalice/Radix/issues">反馈 Bug</a>
·
<a href="https://github.com/overalice/Radix/issues">请求新功能</a>
</p>



<details>
  <summary>目录</summary>
  <ol>
    <li>
      <a href="#关于本项目">关于本项目</a>
    </li>
    <li>
      <a href="#开始">开始</a>
      <ul>
        <li><a href="#依赖">依赖</a></li>
        <li><a href="#安装">安装</a></li>
      </ul>
    </li>
    <li><a href="#使用方法">使用方法</a></li>
    <li><a href="#贡献">贡献</a></li>
  </ol>
</details>


## 关于本项目

本项目提供了一个高效且强大的工具集，用于创建和管理后端接口。

- **简洁的语法和结构：** 使用 Go 语言的简洁性和强大的类型系统，提供了一种清晰易懂的方式来定义和实现后端接口。

- **快速的开发周期：** 通过提供一系列已经封装好的功能和工具，可以极大地减少开发人员编写重复代码的时间，从而加速整个开发周期。

- **灵活的路由和中间件支持：** 通过支持灵活的路由和中间件，能够让开发人员轻松地管理请求和响应，实现更加复杂的业务逻辑。

- **丰富的文档和示例：** 为了帮助开发人员更好地理解和使用本项目，提供了详尽的文档和丰富的示例代码，让他们能够快速上手并开始构建自己的应用程序。



## 开始

这是一份在本地构建项目的指导的例子。
要获取本地副本并且配置运行，你可以按照下面的示例步骤操作。

### 依赖

* Go
  参考：https://golang.google.cn/doc/install

### 安装

1. 初始化项目
   ```sh
   mkdir your-project-name
   cd your-project-name
   go mod init your-project-name
   ```
2. 获取本框架
   ```sh
   go get github.com/overalice/radix
   ```



## 使用方法

创建你的项目文件
   ```sh
touch main.go
   ```

修改其内容

   ```go
package main

import "github.com/overalice/radix"

func main() {
	r := radix.New()
	r.GET("/index", func(ctx *radix.Context) {
		ctx.String("Welcome Radix!")
	})
	r.REST("/person")
	r.Start()
}
   ```

运行项目

```sh
go run main.go
```

转到 [文档](https://github.com/overalice/Radix) 查看更多示例



## 贡献

贡献让开源社区成为了一个非常适合学习、启发和创新的地方。你所做出的任何贡献都是**受人尊敬**的。

如果你有好的建议，请复刻（fork）本仓库并且创建一个拉取请求（pull request）。你也可以简单地创建一个议题（issue），并且添加标签「enhancement」。不要忘记给项目点一个 star！再次感谢！

1. 复刻（Fork）本项目
2. 创建你的 Feature 分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的变更 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到该分支 (`git push origin feature/AmazingFeature`)
5. 创建一个拉取请求（Pull Request）

