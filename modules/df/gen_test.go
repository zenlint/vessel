package df

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	_SIMPLE_TAML = `app:
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
`
)

func Test_Generate(t *testing.T) {
	Convey("Generate one Dockerfile", t, func() {
		dfs, err := Generate([]byte(_SIMPLE_TAML))
		So(err, ShouldBeNil)

		So(dfs, ShouldNotBeNil)
		So(len(dfs), ShouldEqual, 1)

		app := dfs["app"]
		So(app, ShouldNotBeNil)
		So(string(app), ShouldEqual, _SIMPLE_DOCKERFILE)
	})
}
