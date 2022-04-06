# dl-talebook

Downloading books from [talebook server](https://github.com/talebook/talebook), inspired by
syhily's [gist](https://gist.github.com/syhily/9feb936bcaebf2beec567733810f4666).

## Feature

1. Concurrent download.
2. Download from previous progress.
3. Register account on website.
4. Bypass the ratelimit from cloudflare.

## Install

```
go install github.com/hellojukay/dl-talebook@latest
```

## Usage

Execute `dl-talebook -h` to see how to use this download tools.

### Register account

```text
→ dl-talebook register -h
Some talebook website need a user account for downloading books.
You can use this register command for creating account.

Usage:
  dl-talebook register [flags]

Flags:
  -e, --email string      The account email.
  -h, --help              help for register
  -p, --password string   The account password.
  -u, --username string   The account login name.
  -w, --website string    The talebook website.
```

### Download book

```text
→ dl-talebook download -h
Download the book from talebook.

Usage:
  dl-talebook download [flags]

Flags:
  -c, --cookie string       The cookie file name you want to use, it would be saved under the download directory. (default "cookies")
  -d, --download string     The book directory you want to use, default would be current working directory. (default "")
  -f, --format strings      The file formats you want to download. (default [EPUB,MOBI,PDF])
  -h, --help                help for download
  -i, --initial int         The book id you want to start download. It should exceed 0. (default 1)
  -p, --password string     The account password.
  -g, --progress string     The download progress file name you want to use, it would be saved under the download directory. (default "progress")
  -n, --rename              Rename the book file by book ID. (default true)
  -r, --retry int           The max retry times for timeout download request. (default 5)
  -t, --thread int          The number of concurrent download request. (default 1)
  -o, --timeout duration    The max pending time for download request. (default 10s)
  -a, --user-agent string   Set User-Agent for download request. (default "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
  -u, --username string     The account login name.
  -w, --website string      The talebook website.
```
