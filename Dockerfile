FROM golang:1.14.3-alpine AS build 
WORKDIR /src
COPY . . 
RUN go build -o /out/kicker

FROM alpine:3.14 AS bin 
COPY ./assets /assets
COPY --from=build /out/kicker /
CMD "/kicker"