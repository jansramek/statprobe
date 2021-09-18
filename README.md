# statprobe
Takes a list of URLs and probes for given HTTP status.

## Install

```
▶ go get -u github.com/jansramek/statprobe
```

## Usage

statprobe accepts line-delimited urls on `stdin` and outputs matching urls to given `-s` status
on `stdout`:

```
▶ cat urls.txt | statprobe -s 403
http://example1.com
http://example2.com
https://example3.com
```

## Concurrency

You can set the concurrency level by the `-c` flag:

```
▶ cat urls.txt | statprobe -c 32
```

## Debug mode

You can toggle debug mode for printing out filtered responses with `-d` flag:

```
▶ cat urls.txt | statprobe -d
http://example1.com
[403] http://example2.com
http://example4.com
[301] https://example6.com
[error] "error message"
```
