#### eatojoy-order-backend
基于gin-gonic/gin 框架搭建的订单模块微服务MVC架构的空架子。

#### 此项目集成开发常用的功能：
* 基于go-xorm的数据库操作,主从分离配置,连接维持，防止长时间空闲连接报错；
* 配置文件读取
* 计划任务
* swagger文档生成
* Hprose RPC请求
* Log日志生成
* 【尚未实现】基于redis存储的session;
* 【尚未实现】基于redis存储的cache操作；
...
#### 其它注意事项
* 依赖包用dep管理，安装dep
~~~
go get -u github.com/golang/dep/cmd/dep
sudo apt install go-dep
~~~
* 定时器使用[gocron](https://github.com/jasonlvhit/gocron)，

#### 命令：
* cd项目根目录
* 项目初始化

    ~~~
    dep init
    ~~~
* 更新依赖包

    ~~~
    dep ensure -update
    ~~~
* 增加第三方包

    ~~~
    dep ensuer -add 包路径
    ~~~
### swagger文档生成：
* 安装swag

    ~~~
    go get -u github.com/swaggo/swag/cmd/swag
    ~~~
    * 编译swag命令
    
        ~~~
        cd vendor/github.com/swaggo/swag/cmd/swag
        go install -v
        ~~~
        
* 使用[gin-swagger](https://github.com/swaggo/gin-swagger)进行处理
    *安装gin-swagger

    ~~~
    dep ensure -add -v github.com/swaggo/gin-swagger
    dep ensure -add -v github.com/swaggo/gin-swagger/swaggerFiles
    ~~~
* 生成swagger文档

~~~
    swag init
~~~
* swagger api编写例子

   ~~~
   https://segmentfault.com/a/1190000013808421#articleHeader5
   ~~~

#### 运行
~~~
go run *.go
~~~
### 地址:
~~~
http://localhost:5070/home/app/sdafe/123123<br/>
~~~

### swagger文档
~~~
http://localhost:5070/swagger/index.html
~~~