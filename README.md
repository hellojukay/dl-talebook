# dl-talebook
[![check](https://github.com/hellojukay/dl-talebook/actions/workflows/go.yml/badge.svg)](https://github.com/hellojukay/dl-talebook/actions/workflows/go.yml)

Downloading books from [talebook server](https://github.com/talebook/talebook) , inspired by syhily's [gist](https://gist.github.com/syhily/9feb936bcaebf2beec567733810f4666).


1. 整站爬取
2. 断点下载
3. 无第三方依赖，跨平台，可以便捷的运行在群晖上

# Demo
![demo](demo.gif)
# Install
```
go install github.com/hellojukay/dl-talebook@latest
```
# Help
```
hellojukay@local dl-talebook (main) $ ./dl-talebook -h
  -continue
        continue an incomplete download (default true)
  -cookie string
        http cookie
  -dir string
        data dir (default "./")
  -password string
        password
  -site string
        tabebook web site (default "https://book.codefine.site:6870/")
  -start-index int
        start book id
  -timeout duration
        http timeout (default 10s)
  -user-agent string
        http userAgent (default "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
  -username string
        username
  -verbose
        show debug log
  -version
        show progream version

```
