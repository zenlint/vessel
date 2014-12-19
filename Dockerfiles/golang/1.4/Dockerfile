FROM debian:wheezy

RUN apt-get update && apt-get install -y \
		ca-certificates curl gcc libc6-dev make \
		bzr git mercurial \
		--no-install-recommends \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.4
RUN curl -sSL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz \
		| tar -v -C /usr/src -xz
ENV GOROOT /usr/src/go
ENV GOPATH /go
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH