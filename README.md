# KOReader Sync

Yet another KOReader Sync server, written in Go for my own learning and experimentation.

See `docker-compose.yaml` for an example compose setup.


docker run -e LOG_LEVEL=DEBUG -p 8080:8080 -v ./data:/app/data go-ko-sync