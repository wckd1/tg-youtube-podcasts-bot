# YouTube Podcasts Feed

This is self-hosted service for extract audio from YouTube videos and build rss feed that can be added to Podcast app.

It allows to subscribe to YouTube channels/playlists using Telegram bot and automatically get updates to feed.
Also, single item can be added without subscription.

The service uses [yt-dlp](https://github.com/yt-dlp/yt-dlp) to pull videos.

## Bot commands

### `reg`
Create new user with default playlist

### `add`
Add subscription or single video to feed.

Add single item
```
/add https://youtube.com/watch?v={id}
```

Subscribe to channel
```
/add https://youtube.com/c/{id}
/add https://youtube.com/channel/{id}
/add https://youtube.com/{@id}
```

Subscribe to playlist
```
/add https://youtube.com/watch?v={video_id}&list={id}
/add https://youtube.com/playlist?list={id}
```

Filter string can be added to subscription to get only specified updates
```
/add https://youtube.com/c/{id} {some title entry}
```

### `remove`
Removes specified subscription. Same syntax as fo adding subscription.

## API
Servise expose only one endpoint for adding feed to Podcasts app.

- `GET /rss/{key}` - returns generated rss xml with configured limit

## Configuration
Add config.yml file following the example of [example-config.yml](https://github.com/wckd1/tg-youtube-podcasts-bot/blob/main/example-config.yml)

- `feed`
    - `update_interval` - interaval for updates check. Should be set in Golang time.Duration syntax (ex. "1h")
    - `limit` - items count in xml output
- `server`
    - `port` - port for http server to listen to (integer)
    - `rss_key` - secret key that will be added to /rss/ endpoint (string)
- `telegram`
    - `bot_token` - token for Telegram bot to communicate with (string)
    - `debug_mode` - enable extended logging for debug mode (True/False)

## TODO
- Add multi-user support
- Add manually created playlists
- Add fetch old episodes on subscribe
- Optimize yt-dlp commands
- Add godoc comments
- Add tests
- Replace text commands with custom keyboard
