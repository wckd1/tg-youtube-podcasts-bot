# YouTube Podcasts Feed

This is self-hosted service for extract audio from YouTube videos and build rss feed that can be added to Podcast app.

It allows to subscribe to YouTube channels/playlists using Telegram bot and automatically get updates to feed.
Also, single item can be added without subscription.

The service uses [yt-dlp](https://github.com/yt-dlp/yt-dlp) to pull videos.

## Bot commands

### `start`
Create new user with default playlist

### `pl`
List all playlists
```
/pl
```

Create new playlist with given name
```
/pl -new {playlist name}
```

### `add`

Add video to default playlist
```
/add https://youtube.com/watch?v={id}
```

Add video to specified playlist
```
/add https://youtube.com/watch?v={id} -p {playlist id or name}
```

### `sub`
List subscriptions for default playlist
```
/sub
```

List subscriptions for  specified playlist
```
/sub -p {playlist id or name}
```

Add subscription to channel to default playlist
```
/sub -new https://youtube.com/c/{id}
/sub -new https://youtube.com/channel/{id}
/sub -new https://youtube.com/{@id}
```

Add subscription to youtube playlis to default playlist
```
/sub -new https://youtube.com/watch?v={video_id}&list={id}
/sub -new https://youtube.com/playlist?list={id}
```

Add subscription to specified playlist
```
/sub -new https://youtube.com/c/{id} -p {playlist id or name}
```

Filter string can be added to subscription to get only specified updates
```
/sub -new https://youtube.com/c/{id} -f {some title entry}
```

## API
Servise expose only one endpoint for adding feed to Podcasts app.

- `GET /rss/{playlist_id}` - returns generated rss xml

## Configuration
Add config.yml file following the example of [example-config.yml](https://github.com/wckd1/tg-youtube-podcasts-bot/blob/main/example-config.yml)

- `feed`
    - `update_interval` - interaval for updates check. Should be set in Golang time.Duration syntax (ex. "1h")
- `server`
    - `port` - port for http server to listen to (integer)
- `telegram`
    - `bot_token` - token for Telegram bot to communicate with (string)
    - `debug_mode` - enable extended logging for debug mode (True/False)

## TODO
- Add fetch old episodes on subscribe
- Optimize yt-dlp commands
- Add godoc comments
- Add tests
- Replace text commands with custom keyboard
