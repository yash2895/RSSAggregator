# RSSAggregator

A lightweight, concurrent RSS feed aggregator and reader service built in Go and PostgreSQL.

It allows users to register, track RSS feeds, follow/unfollow subscriptions, and read the latest aggregated posts. A background worker periodically scrapes subscribed feeds concurrently and stores articles in the database.

---

## Features

- **API Key Auth**: PostgreSQL automatically generates a unique SHA-256 API key on user registration. Authenticated endpoints use `Authorization: APIKey <key>`.
- **Feed Management**: Add RSS feeds (name & XML URL) and subscribe/unsubscribe via feed follows.
- **Concurrent Background Scraper**: Periodically fetches and parses RSS feeds across 10 concurrent goroutines, automatically deduplicating posts.
- **Personalized Timeline**: Query latest posts published across feeds you follow (`GET /post`).
- **Type-Safe Queries**: Built with [sqlc](https://sqlc.dev/) for type-safe database access and [goose](https://github.com/pressly/goose) for migrations.

---

## Tech Stack

- **Language**: Go 1.22+
- **Database**: PostgreSQL
- **Libraries & Tools**: `database/sql`, `lib/pq`, `godotenv`, `google/uuid`, `sqlc`, `goose`

---

## Quick Start

### 1. Prerequisites

- [Go](https://go.dev/dl/) (1.22+)
- [PostgreSQL](https://www.postgresql.org/)
- [Goose](https://github.com/pressly/goose):
  ```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
  ```

### 2. Configuration

Create a `.env` file in the project root:

```env
PORT=8080
POSTGRES_URL=postgres://postgres:password@localhost:5432/rssaggregator?sslmode=disable
```

### 3. Database Migrations

Run database migrations using goose:

```bash
goose -dir sql/schema postgres "postgres://postgres:password@localhost:5432/rssaggregator?sslmode=disable" up
```

### 4. Run the Server

```bash
go run .
```

The server will start on the configured port and immediately initiate the background scraping worker.

---

## API Reference

All protected endpoints require the header:
```http
Authorization: APIKey <your-api-key>
```

| Method | Endpoint | Auth | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/healthz` | No | Server readiness probe |
| `GET` | `/error` | No | Test 500 error response |
| `POST` | `/users` | No | Create user (`{"name": "Alice"}`) & returns API key |
| `GET` | `/users` | Yes | Get authenticated user profile |
| `POST` | `/feeds` | Yes | Add RSS feed (`{"name": "...", "url": "..."}`) |
| `POST` | `/feed_follow` | Yes | Follow a feed (`{"feed_id": "<uuid>"}`) |
| `GET` | `/feed_follow` | Yes | List feeds followed by authenticated user |
| `DELETE` | `/feed_follow/{id}`| Yes | Unfollow a feed by follow record ID |
| `GET` | `/post` | Yes | Fetch latest 10 posts from followed feeds |

---

## Example Usage

### 1. Create a User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice"}'
```
*Save the `api_key` returned in the response.*

### 2. Add an RSS Feed
```bash
curl -X POST http://localhost:8080/feeds \
  -H "Authorization: APIKey <YOUR_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Boot.dev", "url": "https://blog.boot.dev/index.xml"}'
```

### 3. Follow a Feed
```bash
curl -X POST http://localhost:8080/feed_follow \
  -H "Authorization: APIKey <YOUR_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"feed_id": "<FEED_ID>"}'
```

### 4. Read Posts
```bash
curl http://localhost:8080/post \
  -H "Authorization: APIKey <YOUR_API_KEY>"
```

---

## Useful Commands

```bash
# Build binary
go build -o RSSAggregator .

# Re-generate sqlc Go code (if queries/schema changed)
sqlc generate

# Rollback last migration
goose -dir sql/schema postgres "$POSTGRES_URL" down
```

---

## Future Improvements

- **Configurable Scraper**: Expose worker concurrency and poll interval as environment variables instead of hardcoded values.
- **Pagination & Filtering**: Support `limit`, `offset`, and keyword search on `GET /post` (currently capped at 10).
- **Feed Discovery**: Mount the unrouted `GET /feeds` endpoint to allow users to discover all existing feeds.
- **Multi-Format Feed Support**: Add parser support for Atom and JSON Feed formats in addition to standard RSS 2.0.
- **Read & Bookmark Status**: Track read/unread state and saved articles on a per-user basis.

