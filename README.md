# KOReader Sync

Yet another KOReader Sync server, written in Go for my own learning and experimentation.

See `docker-compose.yaml` for an example compose setup.

Build with: `docker build --tag koreader-sync .`

Run with: `docker run -e LOG_LEVEL=DEBUG -p 8080:8080 -v ./data:/app/data koreader-sync`
Or: `docker-compose up -d`