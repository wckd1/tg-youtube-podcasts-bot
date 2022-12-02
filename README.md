# YouTube Podcasts Feed

This is self-hosted service for extract audio from YouTube videos and build rss feed that can be added to Podcast app.

It allows to subscribe to YouTube channels/playlists using Telegram bot and automatically get updates to feed.
Also, single items can be added withou subscription.

The service uses [yt-dlp](https://github.com/yt-dlp/yt-dlp) to pull videos and [ffmpeg](https://www.ffmpeg.org/) for audio extraction.

## Bot commands

### `add`, `new`, `sub`
Add subscription or single video to feed.

Download and add single item
```
/add https://youtube.com/watch?v={id}
```

Subscribe to channel
```
/add https://youtube.com/c/{id}
/add https://youtube.com/channel/{id}
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

### `remove`, `rm`, `delete`, `unsub`
Removes specified subscription. Same syntax as fo adding subscription.

## API
Servise expose only one endpoint for adding feed to Podcasts app.

- `GET /rss/{key}` - returns generated rss xml with last 20 entries

## Configuration
Add config.yml file following the example of [example-config.yml](https://github.com/wckd1/tg-youtube-podcasts-bot/blob/main/example-config.yml)

- `feed.update_interval` - interaval for check new updates. Should be set in golang time.Duration syntax (ex. "1h")
- `server.port` - port for http server to listen to (integer)
- `server.rss_key` - secret key that will be added to /rss/ endpoint (string)
- `telegram.bot_token` - token for Telegram bot to communicate with (string)
- `telegram.chat_id` - id of chat/group where updates will be posted (integer)
- `telegram.debug_mode` - enable extended logging for debug mode (True/False)

## TODO
- For now, Telegram is used as storage with a limit of [50Mb](https://core.telegram.org/bots/api#sending-files)
- Add sponsorblock
- Retry for downloads/uploads
- Replace text commands with custom keyboard
