## log-alert

[![Test Build Status](https://travis-ci.org/pahoughton/log-alert.png)](https://travis-ci.org/pahoughton/log-alert)

send log scan alerts to prometheus alertmanager

## Features

Scan logfiles specified in log-alert-config.yml and generate alerts
for lines matching supplied regex

### config file

```yaml
global:
  scan_freq: 15m
  metrics_addr: ":9321"

log-files:
  - name: log_httpd_error # alertmanger alert name
    file: /var/http/error.log
    regex:
      - error/i
    annotations:
      sop: http://wiki/sop-log_httpd_error
      match_lines: "{{ $match.lines }}"
      match_first: "{{ $match.lnum_first }}"
      match_last: "{{ $match.lnum_last }}"
      match_range: "{{ $match.lnum_range }}"

  - name: log_httpd_access
    file: "/var/http/access.{{ $date }}.log"
    regex:
      - status=[3-5]\d\d

  - name: log_messages
    file: /var/log/messages
    regex:
      - error/i


alertmanagers:
  - scheme: https
    static_configs:
      - targets:
        - "1.2.3.4:9093"
        - "1.2.3.5:9093"
        - "1.2.3.6:9093"

```

### alert output
```json
[
  {
    "labels": {
      "alertname": "log_httpd_error",
      "regex": "error/i",
      "instance": "hostname:0"
     },
     "annotations": {
       "sop": "http://wiki/sop-log_httpd_error",
       "file": "/var/http/error.log"
       "first_match_line": "2039"
       "last_match_line": "3252"
      }
  },
  {
    "labels": {
      "alertname": "log_messages",
      "regex": "error/i",
      "instance": "hostname:0"
     },
     "annotations": {
       "sop": "http://wiki/sop-messages",
       "file": "/var/http/error.log"
       "match_range": "123-456"
      }
  }
]
```

## Install

Can't

## Usage

install service

## Contribute

https://github.com/pahoughton/log-alert

## Licenses

2019-01-16 (cc) <paul4hough@gmail.com>

GNU General Public License v3.0

See [COPYING](../master/COPYING) for full text.
