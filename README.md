# blogagg

A small RSS aggregator that runs in your terminal. You add feeds, follow the ones you care about, leave `blogagg agg` running in the background to collect posts, and browse the latest ones whenever you want.

I built it as a guided project (the course calls it "gator") to learn Go with Postgres, sqlc and goose. It's small, but it works, and it survives the messy reality of RSS feeds (inconsistent date formats, duplicate posts on every re-scrape, HTML stuffed into descriptions).

## What you need

- **Go** 1.22 or newer, only to install it. Once it's built, `blogagg` is a single static binary and doesn't need the Go toolchain.
- **PostgreSQL** 15 or newer, running locally or anywhere you can reach.
- **goose**, to set up the database schema:

  ```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
  ```

## Install

```bash
go install github.com/refugee8888/blogagg@latest
```

This puts the `blogagg` binary in `$(go env GOPATH)/bin`. Make sure that's on your `PATH`.

## Setup

### 1. Create the database

The database is called `gator`:

```bash
psql postgres -c "CREATE DATABASE gator;"
```

### 2. Run the migrations

Clone the repo and apply the schema:

```bash
git clone https://github.com/refugee8888/blogagg.git
cd blogagg/sql/schema
goose postgres "postgres://INSERT_USERNAME:INSERT_PASSWORD@localhost:5432/gator" up
```

Replace `INSERT_USERNAME` and `INSERT_PASSWORD` with your Postgres credentials, and `localhost:5432` with your host and port if they're different.

### 3. Create the config file

blogagg reads its config from `~/.gatorconfig.json`. Create it with your connection string:

```json
{
  "db_url": "postgres://INSERT_USERNAME:INSERT_PASSWORD@localhost:5432/gator?sslmode=disable"
}
```

Use the same credentials as above.

blogagg writes `current_user_name` to this file itself when you register or log in, so you don't need to add it.

## Usage

```bash
blogagg register paul          # create a user and log in as them
blogagg addfeed "Hacker News" https://news.ycombinator.com/rss
blogagg agg 1m                 # scrape feeds every minute (leave it running)
blogagg browse 10              # show the 10 latest posts from feeds you follow
```

### All commands

| Command | What it does |
|---|---|
| `register <name>` | Create a user and log in as them |
| `login <name>` | Switch to an existing user |
| `users` | List all users, marking the current one |
| `addfeed <name> <url>` | Add a feed and follow it automatically |
| `feeds` | List every feed in the database |
| `follow <url>` | Follow a feed someone already added |
| `following` | List the feeds you follow |
| `unfollow <url>` | Stop following a feed |
| `agg <interval>` | Scrape feeds forever on an interval like `30s`, `1m` or `1h` |
| `browse [limit]` | Show the latest posts from your feeds (default 2) |
| `reset` | Delete all users and everything attached to them. Dev only. |

A note on `agg`: don't set the interval too low against real sites. Something like `1m` is plenty, and you don't want to hammer people's servers.

## Development

```bash
go run . <command> [args]
```

Queries live in `sql/queries` and are compiled with [sqlc](https://sqlc.dev). If you change them, run `sqlc generate`.
