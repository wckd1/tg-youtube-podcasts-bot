# Build stage
FROM umputun/baseimage:buildgo-latest as build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /build
COPY . .

RUN go build -o yt-podcast-bot -ldflags "-s -w" ./cmd/yt-podcast-bot

#  Run stage
FROM umputun/baseimage:app-latest

COPY --from=build /build/yt-podcast-bot /srv/yt-podcast-bot
COPY ./configs/config.yml /srv/configs/config.yml

RUN \
    chown -R app:app /srv && \
    chmod +x /srv/yt-podcast-bot
RUN apk --no-cache add ca-certificates python3 py3-pip
RUN pip3 install --no-cache-dir --no-deps -U yt-dlp
WORKDIR /srv

RUN mkdir storage && mkdir storage/downloads

EXPOSE 6745

CMD ["/srv/yt-podcast-bot"]
