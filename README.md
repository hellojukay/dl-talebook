# dl-talebook
Downloading books from [talebook server](https://github.com/talebook/talebook) , inspired from gist https://gist.github.com/syhily/9feb936bcaebf2beec567733810f4666 .

1. 并发任务
2. 断点下载

# Demo
![demo](demo.png)
# Install
```
go install github.com/hellojukay/dl-talebook@latest
```
# Help
```
hellojukay@local dl-talebook (main) $ ./dl-talebook -h
Usage of ./dl-talebook:
  -c int
    	maximum number of concurrent download tasks allowed per second (default 5)
  -dir string
    	data dir (default "./")
  -password string
    	password
  -site string
    	tabebook web site (default "https://book.codefine.site:6870/")
  -timeout duration
    	http timeout (default 10s)
  -user-agent string
    	http userAgent (default "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
  -username string
    	username
```