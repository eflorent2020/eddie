FROM ubuntu:latest

ENV GOLANG_VERSION 1.9.2
ENV NODE_VERSION 8.9.3
ENV NPM_VERSION 3.10.10

EXPOSE 4567

RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config wget ca-certificates curl git xz-utils libxrender1 \
		libxext6 libfontconfig1 


RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y locales

RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    dpkg-reconfigure --frontend=noninteractive locales && \
    update-locale LANG=en_US.UTF-8

ENV LANG en_US.UTF-8 

RUN set -eux; \
  goRelArch='linux-amd64'; \
  goRelSha256='de874549d9a8d8d8062be05808509c09a88a248e77ec14eb77453530829ac02b' ; \
  url="https://storage.googleapis.com/golang/go${GOLANG_VERSION}.${goRelArch}.tar.gz"; \
	wget -O go.tgz "$url"; \
	echo "${goRelSha256} *go.tgz" | sha256sum -c -; \
	tar -C /usr/local -xzf go.tgz; \
	rm go.tgz; \
	\
	export PATH="/usr/local/go/bin:$PATH"; \
	go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# gpg keys listed at https://github.com/nodejs/node#release-team
RUN gpg --keyserver pool.sks-keyservers.net --recv-keys 94AE36675C464D64BAFA68DD7434390BDBE9B9C5
RUN gpg --keyserver pool.sks-keyservers.net --recv-keys FD3A5288F042B6850C66B31F09FE44734EB7990E
RUN gpg --keyserver pool.sks-keyservers.net --recv-keys 71DCFD284A79C3B38668286BC97EC7A07EDE3FC1
RUN gpg --keyserver pool.sks-keyservers.net --recv-keys DD8F2338BAE7501E3DD5AC78C273792F7D83545D
RUN gpg --keyserver pool.sks-keyservers.net --recv-keys C4F0DFFF4E8C1A8236409D08E73BC641CC11F4C8
RUN gpg --keyserver pool.sks-keyservers.net --recv-keys B9AE9905FFD7803F25714661B63B535A4C206CA9
RUN gpg --keyserver pool.sks-keyservers.net --recv-keys 56730D5401028683275BD23C23EFEFE93C4CFFFE
RUN gpg --keyserver pool.sks-keyservers.net --recv-keys 77984A986EBC2AA786BC0F66B01FBB92821C587A	

RUN curl -SLO "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.gz" \
	&& curl -SLO "https://nodejs.org/dist/v$NODE_VERSION/SHASUMS256.txt.asc" \
	&& gpg --verify SHASUMS256.txt.asc \
	&& grep " node-v$NODE_VERSION-linux-x64.tar.gz\$" SHASUMS256.txt.asc | sha256sum -c - \
	&& tar -xzf "node-v$NODE_VERSION-linux-x64.tar.gz" -C /usr/local --strip-components=1 \
	&& rm "node-v$NODE_VERSION-linux-x64.tar.gz" SHASUMS256.txt.asc \
	&& npm install -g npm@"$NPM_VERSION" \
	&& npm cache clear

RUN useradd -ms /bin/bash eddie
WORKDIR /home/eddie
ADD . .

RUN cd eddie-frontend ; \
  npm install ; \
  npm run build

ENV GOBIN=$GOPATH/bin
RUN mkdir -p /go/src/github.com/unidoc && cd /go/src/github.com/unidoc && rm -fr unidoc && git clone https://github.com/emmanuel-florent/unidoc.git 
RUN go get ./...

#RUN mkdir -p /go/src/github.com/unidoc && cd /go/src/github.com/unidoc && rm -fr unidoc && git clone https://github.com/emmanuel-florent/unidoc.git 
RUN go build

RUN wget https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/0.12.4/wkhtmltox-0.12.4_linux-generic-amd64.tar.xz
RUN tar -xvf wkhtmltox-0.12.4_linux-generic-amd64.tar.xz
RUN mv wkhtmltox/bin/wkhtmltopdf .
#RUN apt-get remove -qy perl git curl wget && rm -rf /var/lib/apt/lists/*
#RUN rm -rf node_modules *.go eddie-frontend/src eddie-frontend/test

ENTRYPOINT [ "eddie" ]
