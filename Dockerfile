FROM golang:1.16 as base

WORKDIR /opt/app

FROM base as prod

COPY . .
RUN go build
CMD ["./hsa-homework3"]