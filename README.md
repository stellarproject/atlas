# Atlas

The Atlas Nameserver service provides a simple programmable [DNS](https://www.cloudflare.com/learning/dns/what-is-dns/) service.
Atlas uses the same core library that CoreDNS uses ([miekg/dns](https://github.com/miekg/dns)).

[![Go Report Card](https://goreportcard.com/badge/github.com/ehazlett/atlas)](https://goreportcard.com/report/github.com/ehazlett/atlas) [![Docs](https://godoc.org/github.com/ehazlett/atlas?status.svg)](http://godoc.org/github.com/ehazlett/atlas) [![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fehazlett%2Fatlas%2Fbadge&style=flat)](https://actions-badge.atrox.dev/ehazlett/atlas/goto)

# Usage

```
NAME:
   atlas - simple dns service

USAGE:
   atlas [global options] command [command options] [arguments...]

VERSION:
   0.1.0 (c568ab98) linux/amd64

AUTHOR:
   @ehazlett

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -D                  enable debug logging
   --bind value, -b value       bind address for the DNS service (default: "udp://0.0.0.0:53")
   --datastore value, -d value  uri for datastore backend (default: "localdb:///etc/atlas/atlas.db")
   --address value, -a value    grpc address (default: "tcp://127.0.0.1:9000")
   --upstream-dns value         upstream dns server (default: "9.9.9.9:53")
   --cache-ttl value            builtin cache ttl (default: disabled) (default: 0s)
   --help, -h                   show help
   --version, -v                print the version
```

# Example

Atlas has a GRPC api that enables remote management of the internal DNS store.

## Add a Record

```
NAME:
   actl create - create nameserver record

USAGE:
   actl create [command options] <NAME> <VALUE>

OPTIONS:
   --type value, -t value  resource record type (A, CNAME, TXT, SRV, MX) (default: "A")
```

To add a new `A` record:

```bash
$> actl create -t A foo.local 127.0.0.1
added foo.local=127.0.0.1 (A)
```

This will create a new `A` record for `foo.local` that resolves to `127.0.0.1`.  You
can use `dig` to query Atlas:

```
$> dig @127.0.0.1 foo.local

; <<>> DiG 9.11.5-P1-1ubuntu2.5-Ubuntu <<>> @localhost foo.local
; (1 server found)
;; global options: +cmd
;; Got answer:
;; WARNING: .local is reserved for Multicast DNS
;; You are currently testing what happens when an mDNS query is leaked to DNS
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 2425
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;foo.local.                     IN      A

;; ANSWER SECTION:
foo.local.              10      IN      A       127.0.0.1

;; Query time: 5 msec
;; SERVER: 127.0.0.1#53(127.0.0.1)
;; WHEN: Sat Aug 17 01:06:36 EDT 2019
;; MSG SIZE  rcvd: 52

```

## List Records

You can list all records in the Atlas store:

```bash
$> actl list
NAME                TYPE                VALUE               OPTIONS
foo.local           A                   127.0.0.1
```

## Delete Records

You can delete records from Atlas as well:

```bash
$> actl remove foo.local
removed foo.local
```

