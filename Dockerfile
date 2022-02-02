FROM golang:1.14.3-alpine AS build 
WORKDIR /src
COPY ./go.mod go.mod
COPY ./go.sum go.sum
RUN go mod download
COPY . . 
RUN go build -o /out/kicker

FROM alpine:3.14 AS bin 
COPY ./assets /assets
COPY --from=build /out/kicker /
CMD "/kicker"