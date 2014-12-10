# YAML 文件格式规范

本文档定义了目前版本的 YAML 文件格式规范，将配合样例文件讲解（除非特别注明，否则均为必填项）：

```yaml
# 每个文件有切仅有一个 app 分区
app:
    image: bradrydzewski/go:1.2 # 用于指定使用哪个 Docker 基础镜像
    git: # 用于指定需要克隆的仓库及相关信息
        auth: username:password         # HTTP 基本授权信息（可选）
        path: github.com/drone/drone    # Git 仓库地址
    environment: # 用于指定环境变量（可选）
      - GOROOT=/usr/local/go
      - PATH=$PATH:$GOROOT/bin:$GOPATH/bin
    services: # 用于指定依赖服务，每个依赖需要另外建立分区说明详细信息（可选）
      - redis # 依赖服务名称是任意的，但必须和相关分区名称相匹配
      - mysql
      - log
      - oauth
    script: # 用于指定需要执行的脚本（可选）
      - go get -u github.com/hoisie/redis
      - go get -u github.com/go-martini/martini
      - go get -u github.com/martini-contrib/render
      - go get -u github.com/go-sql-driver/mysql
      - go get -u github.com/go-xorm/xorm
      - go build
    notify: # 用于指定通知信息（可选）
      email: # 用于指定邮箱通知信息
        recipients:
          - u@docker.cn

# 依赖服务分区，名称需要和 app 分区内的名称相匹配（可选）
oauth:
    image: bradrydzewski/go:1.3
    git:
      path: github.com/drone/oauth
    environment:
      - GOROOT=/usr/local/go
      - PATH=$PATH:$GOROOT/bin:$GOPATH/bin
    services: # 依赖服务同样可以具有自己的依赖服务
      - redis-oauth
      - mysql-oauth
    script:
      - go get -u github.com/hoisie/redis
      - go get -u github.com/go-martini/martini
      - go get -u github.com/martini-contrib/render
      - go get -u github.com/go-sql-driver/mysql
      - go get -u github.com/go-xorm/xorm
      - go build
    ports: # 用于指定端口映射
        - "80:80"
    expose: # 用于指定暴露的端口号
        - "3000"

log:
    image: bradrydzewski/go:1.1
    git:
      path: github.com/drone/log
    environment:
      - GOROOT=/usr/local/go
      - PATH=$PATH:$GOROOT/bin:$GOPATH/bin
    services:
      - redis-log
      - mysql-log
    script:
      - go get -u github.com/hoisie/redis
      - go get -u github.com/go-martini/martini
      - go get -u github.com/martini-contrib/render
      - go get -u github.com/go-sql-driver/mysql
      - go get -u github.com/go-xorm/xorm
      - go build
    ports:
        - "80:80"
    expose:
        - "3000"

redis:
    image: bradrydzewski/redis:2.6
    ports:
        - "80:80"
    expose:
        - "3000"

mysql:
    image: bradrydzewski/mysql:5.5
    ports:
        - "22:22"
    expose:
        - "3306"
    environment:
        - RACK_ENV=development
        - SESSION_SECRE

redis-oauth:
    image: bradrydzewski/redis:2.6
    ports:
        - "80:80"
    expose:
        - "3000"

mysql-oauth:
    image: bradrydzewski/mysql:5.5
    ports:
        - "22:22"
    expose:
        - "3306"
    environment:
        - RACK_ENV=development
        - SESSION_SECRE

redis-log:
    image: bradrydzewski/redis:2.6
    ports:
        - "80:80"
    expose:
        - "3000"

mysql-log:
    image: bradrydzewski/mysql:5.5
    ports:
        - "22:22"
    expose:
        - "3306"
    environment:
        - RACK_ENV=development
        - SESSION_SECRE
```
