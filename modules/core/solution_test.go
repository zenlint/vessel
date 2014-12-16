package core

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	_SIMPLE_YAML = `app:
    image: bradrydzewski/go:1.2 
    git: 
        auth: username:password         
        path: github.com/drone/drone    
    environment: 
        - GOROOT=/usr/local/go
        - PATH=$PATH:$GOROOT/bin:$GOPATH/bin
    script: 
        - go get -u github.com/hoisie/redis
        - go get -u github.com/go-martini/martini
        - go get -u github.com/martini-contrib/render
        - go get -u github.com/go-sql-driver/mysql
        - go get -u github.com/go-xorm/xorm
        - go build
    cmd: ["./gogs", "web"]
    notify: 
        email:
            recipients:
                - u@docker.cn`
	_SIMPLE_DOCKERFILE = `FROM bradrydzewski/go:1.2

ENV GOROOT /usr/local/go
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

RUN go get -u github.com/hoisie/redis
RUN go get -u github.com/go-martini/martini
RUN go get -u github.com/martini-contrib/render
RUN go get -u github.com/go-sql-driver/mysql
RUN go get -u github.com/go-xorm/xorm
RUN go build

CMD ["./gogs", "web"]
`

	_COMPLEX_YAML = `app:
    image: bradrydzewski/go:1.2 
    git:
        auth: username:password         
        path: github.com/drone/drone   
    environment: 
        - GOROOT=/usr/local/go
        - PATH=$PATH:$GOROOT/bin:$GOPATH/bin
    script: 
        - go get -u github.com/hoisie/redis
        - go get -u github.com/go-martini/martini
        - go build
    cmd: ["./gogs", "web"]
    services:
        - redis 
        - mysql
        - log
    notify: 
        email:
            recipients:
                - u@docker.cn

log:
    image: bradrydzewski/go:1.1
    git:
        path: github.com/drone/log
    environment:
        - GOROOT=/usr/local/go
        - PATH=$PATH:$GOROOT/bin:$GOPATH/bin
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
    services:
        - redis-log
        - mysql-log

redis:
    image: bradrydzewski/redis:2.6
    ports:
        - "80:80"
    expose:
        - "3000"

mysql:
    image: bradrydzewski/mysql:5.5
    environment:
        - RACK_ENV=development
        - SESSION_SECRE
    ports:
        - "22:22"
    expose:
        - "3306"

redis-log:
    image: bradrydzewski/redis:2.6
    ports:
        - "80:80"
    expose:
        - "3000"

mysql-log:
    image: bradrydzewski/mysql:5.5
    environment:
        - RACK_ENV=development
        - SESSION_SECRE
    ports:
        - "22:22"
    expose:
        - "3306"`
)

func Test_NewSolutionFromBytes(t *testing.T) {
	Convey("New solution from bytes", t, func() {
		sln, err := NewSolutionFromBytes([]byte(_COMPLEX_YAML))
		So(err, ShouldBeNil)
		So(sln, ShouldNotBeNil)

		So(sln.Name, ShouldEqual, "app")
		So(sln.Image, ShouldEqual, "bradrydzewski/go:1.2")
		So(strings.Join(sln.Environment, "\n"), ShouldEqual, "GOROOT=/usr/local/go\nPATH=$PATH:$GOROOT/bin:$GOPATH/bin")
		So(strings.Join(sln.Script, "\n"), ShouldEqual, "go get -u github.com/hoisie/redis\ngo get -u github.com/go-martini/martini\ngo build")
		So(strings.Join(sln.Cmd, " "), ShouldEqual, "./gogs web")
		So(sln.Notify.Email.Recipients[0], ShouldEqual, "u@docker.cn")

		Convey("Dependencies of app service", func() {
			So(len(sln.Services), ShouldEqual, 3)

			s := sln.Services[0]
			So(s.Ports[0], ShouldEqual, "80:80")
			So(s.Expose[0], ShouldEqual, "3000")
		})
	})
}

func Test_Service_Dockerfile(t *testing.T) {
	Convey("Generate Dockerfile from service", t, func() {
		sln, err := NewSolutionFromBytes([]byte(_SIMPLE_YAML))
		So(err, ShouldBeNil)
		So(sln, ShouldNotBeNil)

		data, err := sln.Service.Dockerfile()
		So(err, ShouldBeNil)
		So(string(data), ShouldEqual, _SIMPLE_DOCKERFILE)
	})
}
