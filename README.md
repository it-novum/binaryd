# binaryd - A micro web server build in Go

binaryd is a web server to execute pre-defined commands through an HTTP API. It could be useful in a container
environment.

Commands can defined via a ini file like so
```ini
[ps]
command=/usr/bin/ps -eaf

[whoami]
command=whoami
```

## Fair warning
At the moment, `binaryd` has no support for authentication or encryption.
Making commands executable through an HTTP API can be dangerous.
**Only make use of this, if you know what you are doing! You have been warned!**

## Usage
```
./binaryd /path/to/config.ini
```


## Query using `curl` (plain text results)

### Get available commands and help
```
$ curl http://127.0.0.1:9099/

binaryd HTTP wrapper to execute pre-defined commands through HTTP.

Execute command and get plain text result:
    - curl -X GET http://xxx.xxx.xxx.xxx:9099
    - curl -X GET http://xxx.xxx.xxx.xxx:9099/<command>
    - curl -X GET http://xxx.xxx.xxx.xxx:9099/ps

Execute command and get result as JSON
    - curl -X GET http://xxx.xxx.xxx.xxx:9099/json/<command>
    - curl -X GET http://xxx.xxx.xxx.xxx:9099/json/ps

Available commands are:
ps
whoami
```

### Run command `ps`
```
$ curl http://127.0.0.1:9099/ps
your 131072x1 screen size is bogus. expect trouble
UID        PID  PPID  C STIME TTY          TIME CMD
root         1     0  0 11:19 ?        00:00:00 /init
root         9     1  0 11:20 ?        00:00:00 /init
root        10     9  0 11:20 ?        00:00:02 /init
dziegler    11    10  0 11:20 pts/0    00:00:00 -bash
dziegler    27    11  0 11:20 pts/0    00:00:03 ssh vmware
root       168     1  0 14:02 ?        00:00:00 /init
root       169   168  0 14:02 ?        00:00:00 /init
dziegler   170   169  0 14:02 pts/1    00:00:00 -bash
dziegler   186   170  0 14:02 pts/1    00:00:00 ssh dziegler-docker.itsm.love
root       187     1  0 14:02 ?        00:00:00 /init
root       188   187  0 14:02 ?        00:00:00 /init
dziegler   189   188  0 14:02 pts/2    00:00:00 -bash
dziegler   222   189  0 14:02 pts/2    00:00:01 ssh dziegler-docker.itsm.love
root       247     1  0 17:52 ?        00:00:00 /init
root       248   247  0 17:52 ?        00:00:00 /init
dziegler   249   248  0 17:52 pts/3    00:00:00 -bash
root       427     1  0 17:53 ?        00:00:00 /init
root       428   427  0 17:53 ?        00:00:00 /init
dziegler 20044   249  2 19:45 pts/3    00:00:00 curl http://127.0.0.1:9099/ps
dziegler 20045 20007  0 19:45 pts/4    00:00:00 /usr/bin/ps -eaf
```

### Run command `whoami`
```
$ curl http://127.0.0.1:9099/whoami
dziegler
```

## Query using `curl` (result as json)

### Run command `ps`
```json
$ curl http://127.0.0.1:9099/json/ps
{"stdout":"your 131072x1 screen size is bogus. expect trouble\nUID        PID  PPID  C STIME TTY          TIME CMD\nroot         1     0  0 11:19 ?        00:00:00 /init\nroot         9     1  0 11:20 ?        00:00:00 /init\nroot        10     9  0 11:20 ?        00:00:02 /init\ndziegler    11    10  0 11:20 pts/0    00:00:00 -bash\ndziegler    27    11  0 11:20 pts/0    00:00:03 ssh vmware\nroot       168     1  0 14:02 ?        00:00:00 /init\nroot       169   168  0 14:02 ?        00:00:00 /init\ndziegler   170   169  0 14:02 pts/1    00:00:00 -bash\ndziegler   186   170  0 14:02 pts/1    00:00:00 ssh dziegler-docker.itsm.love\nroot       187     1  0 14:02 ?        00:00:00 /init\nroot       188   187  0 14:02 ?        00:00:00 /init\ndziegler   189   188  0 14:02 pts/2    00:00:00 -bash\ndziegler   222   189  0 14:02 pts/2    00:00:01 ssh dziegler-docker.itsm.love\nroot       247     1  0 17:52 ?        00:00:00 /init\nroot       248   247  0 17:52 ?        00:00:00 /init\ndziegler   249   248  0 17:52 pts/3    00:00:00 -bash\nroot       427     1  0 17:53 ?        00:00:00 /init\ndziegler 20055   249  0 19:46 pts/3    00:00:00 curl http://127.0.0.1:9099/json/ps\ndziegler 20056 20007  0 19:46 pts/4    00:00:00 /usr/bin/ps -eaf\n","rc":0,"execution_unix_timestamp_sec":1684259184}
```

### Run command `whoami`
```json
$ curl http://127.0.0.1:9099/json/whoami
{"stdout":"dziegler\n","rc":0,"execution_unix_timestamp_sec":1684259273}
```

## License
```
MIT License

Copyright (c) 2023 it-novum GmbH

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
