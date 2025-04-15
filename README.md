
#Gator

## Gator is a simple CLI tool, that adds blogs using the RSS - Really Simple Syndication

## Setup Instructions

1. Prerequisites:
   - Go installed (version 1.x or later)
   - PostgreSQL installed and running

2. Build the application:

go build -o gator

This will create a binary exec file, then have to move that to the root by this, 

`sudo mv gator /usr/local/bin/`

Once you done, you can able to run the program by typing `gator {command} {Arguments}`


3. Make it globally available (optional):

sudo mv gator /usr/local/bin/

Once you set up Go and Postgres, run this command to build, 

`go build -o gator`

## What can we do in the gator

# Gator RSS Reader - Command Reference

## User Management
- `gator register {name}`      - Register/signup as a new user
- `gator user {name}`          - Login as a specific user
- `gator users`                - List all registered users

## Feed Management
- `gator addfeed {name} {url}` - Add a new feed to the system
- `gator feeds`                - List all available feeds
- `gator follow {url}`         - Follow a specific feed
- `gator following`            - Show feeds you're currently following
- `gator unfollow {url}`       - Unfollow a specific feed

## Content
- `gator agg {time_interval}`  - Aggregate posts from feeds (runs in a loop)
  `gator agg 5s`

- `gator browse {limit}`       - View recent posts (default limit: 2)
  `gator browse 2`

## System
- `gator reset`                - Erase and reset everything