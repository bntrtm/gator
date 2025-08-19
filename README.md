# Gator

## User-Friendly CLI tool for aggregating RSS Feeds

Gator is a blog aggreGATOR with a CLI that interacts with a Postgres database.
Users can separate their feeds into more than one "user," tracked within the relational database.

## Requirements

To install the `gator` CLI, you can use: go install github.com/bntrtm/gator
Gator expects a file be found at the following path: ~/$HOME/.config/gator/.gatorconfig.json.

## User Instructions

1. Set up your config file (_.gatorconfig.json_) with the following contents:
```
{
  "db_url": "protocol://username:password@host:port/database"
}

```
In this case, the protocol will be `postgres`, and the database `gator`.
Postgres uses a default port of `5432`.

Fill in remaining fields (username, password) as accordingly.

2. Run the program with any of the following commands:
    - `register $1`: Register a new user `$1`.
    - `login $1`: Login as user `$1`.
    - `users`: list all users in the database.
    - `addfeed $1 $2`: make available to follow the RSS feed titled `$1` at URL `$2`.
    - `feeds`: list all feeds available to follow.
    - `follow $1`: follow the RSS feed at URL `$1`.
    - `unfollow $1`: unfollow the RSS feed at URL `$1`.
    - `following`: list RSS feeds that the current logged-in user is following.
    - `agg $1`: send GET requests to RSS feeds every `$1` units of time (1s, 1m, 1h...)
    - `browse $1` see the latest `$1` posts across all of the current user's followed RSS feeds (default: 2).
    - `reset`: Deletes ALL users and user data from your database (This command is commented out in production by default. Uncomment the registration of this command in `main.go` and build from source if you want it available)