# Gator

Gator is a CLI tool for aggregating RSS feeds, built in Go.

## Requirements

- Go (version 1.27.1 used)
- PostgreSQL

## Installation

Install the gator CLI:

    go install github.com/Contary925/gator@latest

## Configuration

Gator reads its configuration from `~/.gatorconfig.json`. Create this file
with the following structure:

    {
      "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
      "current_user_name": ""
    }

(Fill in your own Postgres credentials)

## Database Setup

Gator uses [goose](https://github.com/pressly/goose) for migrations. Run:

    goose postgres "<your connection string>" up

from the `sql/schema` directory to set up the required tables.

## Usage

    gator <command> [args]

### Commands

- `register <name>` — create a new user
- `login <name>` — log in as an existing user
- `reset` — reset the database (explain what this actually does)
- `users` — list all registered users
- `agg <duration>` — continuously fetch feeds on an interval (e.g. `1m`, `30s`; intervals like `1m30s` work too.)
- `addfeed <name> <url>` — add a new feed (requires login)
- `feeds` — list all feeds
- `follow <url>` — follow an existing feed (requires login)
- `following` — list feeds you're following (requires login)
- `unfollow <url>` — unfollow a feed (requires login)
- `browse [limit]` — view recent posts from followed feeds (requires login)