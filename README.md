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
Add app.env file following the example of [app.env.example](https://github.com/wckd1/tg-youtube-podcasts-bot/blob/main/app.env.example)

- `BOT_API_TOKEN` - token for Telegram bot to communicate with (string)
- `CHAT_ID` - id of chat/group where updates will be posted (integer)
- `DEBUG_MODE` - enable extended logging for debug mode (True/False)
- `UPDATE_INTERVAL` - interaval for check new updates. Should be set in golang time.Duration syntax (ex. "1h")
- `RSS_KEY` - secret key that will be added to /rss/ endpoint (string)
- `PORT` - port for http server to listen to (integer)

## TODO
- For now, Telegram is used as storage with a limit of [50Mb](https://core.telegram.org/bots/api#sending-files)
- Add sponsorblock
- Retry for downloads/uploads
- Replace text commands with custom keyboard
