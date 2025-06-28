# boot-blog-agg

`boot-blog-agg`, AKA `gator`, is a command line tool that can aggregate posts from multiple RSS feeds and display the most recent ones to the user.

## prerequisites

- `go` >= 1.22
- Postgres >= 16.9

## installation

Explain to the user how to install the gator CLI using go install.

1. Clone this repo
2. Run `go install` in the cloned directory
3. Run `boot-blog-agg` to get a list of available commands

## Configuration

### First setup

1. Create the config file in your home directory: `touch ~/.gatorconfig.json`
2. Add your postgres access url
   ```json
   {"db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"}
   ```
3. Regsiter a user to be able to use other commands
