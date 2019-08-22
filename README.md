# Atlas

The Atlas Nameserver service provides a simple programmable [DNS](https://www.cloudflare.com/learning/dns/what-is-dns/) service.
Atlas uses the same core library that CoreDNS uses ([miekg/dns](https://github.com/miekg/dns)).

[![Go Report Card](https://goreportcard.com/badge/github.com/ehazlett/atlas)](https://goreportcard.com/report/github.com/ehazlett/atlas) [![Docs](https://godoc.org/github.com/ehazlett/atlas?status.svg)](http://godoc.org/github.com/ehazlett/atlas) [![Build Status](https://action-badges.now.sh/JasonEtco/action-badges)](https://github.com/ehazlett/atlas/actions)

# Installation

## Docker
[![Docker Image](https://img.shields.io/badge/ehazlett%2Fatlas%3Alatest-docker%20image-blue)](https://cloud.docker.com/u/ehazlett/repository/docker/ehazlett/atlas)

A [Docker](https://www.docker.com) image is built upon every merge to master.  The latest is [ehazlett/atlas:latest](https://cloud.docker.com/u/ehazlett/repository/docker/ehazlett/atlas)

To run with Docker, you will need to map ports (-p) for all ports you want to expose (at least 53/udp and 9000/tcp).

### Examples

Run Atlas in Docker publishing on default ports:

```bash
$> docker run -ti -d \
	--name atlas \
	-p 53:53/udp \
	-p 9000:9000 \
	ehazlett/atlas:latest -b udp://0.0.0.0:53
```

Run Atlas using Google DNS for upstream with a 30s cache:
```bash
$> docker run -ti -d \
	--name atlas \
	-p 53:53/udp \
	-p 9000:9000 \
	ehazlett/atlas:latest -b udp://0.0.0.0:53 --upstream-dns 8.8.8.8:53 --cache-ttl 30s
```

## Manual

[![Linux Latest Build](https://img.shields.io/badge/linux-latest-green)](https://ehazlett-public.s3.us-east-2.amazonaws.com/atlas/atlas-linux-latest.zip) [![FreeBSD Latest Build](https://img.shields.io/badge/freebsd-latest-green)](https://ehazlett-public.s3.us-east-2.amazonaws.com/atlas/atlas-freebsd-latest.zip) [![Windows Latest Build](https://img.shields.io/badge/windows-latest-green)](https://ehazlett-public.s3.us-east-2.amazonaws.com/atlas/atlas-windows-latest.zip)

- Download a release either from the "latest" builds above or from the [Releases](https://github.com/ehazlett/atlas/releases).
- Extract the release to your `PATH`:
  - Linux: `unzip -d /usr/local/bin atlas-linux-latest.zip`
  - FreeBSD: `unzip -d /usr/local/bin atlas-freebsd-latest.zip`
  - Windows: `unzip -d C:\Windows\system32 atlas-windows-latest.zip`

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
   actl create [command options] [arguments...]

OPTIONS:
   --name value, -n value    record name
   --record value, -r value  record to add (format: <TYPE>:<VALUE>)
```

To add a new `A` record:

```bash
$> actl create -n foo.int r A:127.0.0.1
added 1 record
```

This will create a new `A` record for `foo.int` that resolves to `127.0.0.1`.  You
can use `dig` to query Atlas:

```
$> dig @127.0.0.1 foo.int

; <<>> DiG 9.11.5-P1-1ubuntu2.5-Ubuntu <<>> @localhost foo.int
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 2425
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;foo.int.                     IN      A

;; ANSWER SECTION:
foo.int.              10      IN      A       127.0.0.1

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
foo.int             A                   127.0.0.1
```

## Delete Records

You can delete records from Atlas as well:

```bash
$> actl remove foo.int
removed foo.int
```

# Metrics
Atlas has the ability to expose [Prometheus[] compatible metrics.  By default, the endpoint is disabled.  To enable, add the `--metrics-addr` flag:


You can then scrape the `<metrics-addr>/metrics` endpoint.  You will see something like this:

```
# HELP atlas_dns_create_total Total number of record creates
# TYPE atlas_dns_create_total counter
atlas_dns_create_total 2
# HELP atlas_dns_delete_total Total number of record deletes
# TYPE atlas_dns_delete_total counter
atlas_dns_delete_total 1
# HELP atlas_dns_lookup_a_total Total number of A record lookups
# TYPE atlas_dns_lookup_a_total counter
atlas_dns_lookup_a_total 201
# HELP atlas_dns_lookup_cname_total Total number of CNAME record lookups
# TYPE atlas_dns_lookup_cname_total counter
atlas_dns_lookup_cname_total 201
# HELP atlas_dns_lookup_forward_total Total number of upstream lookups
# TYPE atlas_dns_lookup_forward_total counter
atlas_dns_lookup_forward_total 452
# HELP atlas_dns_query_milliseconds Duration of query in milliseconds
# TYPE atlas_dns_query_milliseconds histogram
atlas_dns_query_milliseconds_bucket{le="1"} 402
atlas_dns_query_milliseconds_bucket{le="5"} 402
atlas_dns_query_milliseconds_bucket{le="10"} 402
atlas_dns_query_milliseconds_bucket{le="25"} 837
atlas_dns_query_milliseconds_bucket{le="50"} 851
atlas_dns_query_milliseconds_bucket{le="100"} 854
atlas_dns_query_milliseconds_bucket{le="250"} 854
atlas_dns_query_milliseconds_bucket{le="500"} 854
atlas_dns_query_milliseconds_bucket{le="1000"} 854
atlas_dns_query_milliseconds_bucket{le="+Inf"} 854
atlas_dns_query_milliseconds_sum 6565.383260999997
atlas_dns_query_milliseconds_count 854
```
