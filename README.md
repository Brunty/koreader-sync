# KOReader Sync

Yet another KOReader Sync server, written in Go for my own learning and experimentation.

API endpoints based on https://github.com/Open-Audiobook/koreader-sync-protocol

## Notes

This is the server part of koreader-sync, how things like document IDs are computed, and progress strings are calculated
are part of the client (e-reader device typically) communicating with this API. This API doesn't concern itself with
how they're formatted, it's just here to store them

## Hosting and running koreader-sync

An example `docker-compose.yml` setup:

```yaml
services:
  app:
    image: brunty/koreader-sync:v0.1
    ports:
      - "8080:8080"
    environment:
      - LOG_LEVEL=DEBUG # Debug puts a lot of info in the logs, not useful if you're not actively debugging something
    volumes:
      - data:/app/data # this is where the sqlite db file is stored
    restart: unless-stopped

volumes:
  data:
```

Then run `docker-compose up -d` and the server will be running at `http://localhost:8080` and the database will be
stored in the Docker volume

## Connecting Your KOReader Device

1. Open a book (or document) in KOReader on your device.
2. Go to the progress sync settings (Settings > Progress Sync > Custom Sync Server) - enter your server URL
(http://localhost:8080 if it's running locally, or your own URL if behind something like a reverse proxy / tunnel)
3. Select "Register / Login" and enter your username and password (to set up your user the first time, select
"Register", for subsequence devices you can select "Login")
4. Select "Push progress from this device now" and progress will be sent and stored in koreader-sync
5. Set up automatic sync to your preferences ("Automatically keep documents in sync" and then configuring how often sync
is sent)

