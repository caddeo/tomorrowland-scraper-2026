# Tomorrowland Lineup Scraper

Never really finished it, but the idea is there.

Scrapes the [Tomorrowland Lineup 2026](https://belgium.tomorrowland.com/en/line-up) for the first week and sends it to a configured google sheet.

Built using Golang:

- Colly for web scraping
- Google Sheets API for sending the data to a google sheet.
- Oauth
- Bubbletea for the terminal UI.

## Configuration

Google api token
Google sheet id
...

## Usage

### Scraping

Run `go run ...` ..

### Using archive

Every scrape is stored using a flat json file and can be reloaded using ...
