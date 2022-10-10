FROM golang:1.18-alpine AS builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /usr/src/app

ENV CGO_ENABLE 0
ENV GOPROXY='https://goproxy.cn,direct'

ADD go.mod .
ADD go.sum .
RUN go mod download && go mod verify

COPY . .

RUN go build -ldflags="-s -w" -v -o /usr/local/bin/app

FROM alpine as releaser

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache ca-certificates

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia

ENV TZ Asia/Shanghai

WORKDIR /

COPY --from=builder /usr/local/bin/app .

EXPOSE 80

CMD ["/app"]