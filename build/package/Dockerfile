# Build stage
FROM umputun/baseimage:buildgo-latest as build

ENV GOFLAGS="-mod=mod"
ENV CGO_ENABLED=0

ADD . /build
WORKDIR /build

RUN go build -o /build/yt-podcast-bot -ldflags "-s -w" cmd/yt-podcast-bot/main.go

#  Run stage
FROM umputun/baseimage:app-latest

COPY --from=build /build/yt-podcast-bot /srv/yt-podcast-bot
COPY app.env /srv

RUN \
    chown -R app:app /srv && \
    chmod +x /srv/yt-podcast-bot
RUN apk --no-cache add ca-certificates ffmpeg python3 py3-pip
RUN pip3 install --no-cache-dir --no-deps -U yt-dlp

WORKDIR /srv
RUN mkdir storage && mkdir storage/downloads

EXPOSE 6745

CMD ["/srv/yt-podcast-bot"]
ENTRYPOINT ["/init.sh"]
