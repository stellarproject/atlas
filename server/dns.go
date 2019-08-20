/*
   Copyright 2019 Evan Hazlett <ejhazlett@gmail.com>

   Permission is hereby granted, free of charge, to any person obtaining a copy of
   this software and associated documentation files (the "Software"), to deal in the
   Software without restriction, including without limitation the rights to use, copy,
   modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
   and to permit persons to whom the Software is furnished to do so, subject to the
   following conditions:

   The above copyright notice and this permission notice shall be included in all copies
   or substantial portions of the Software.

   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
   INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
   PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
   FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
   TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE
   USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package server

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/containerd/typeurl"
	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
	"github.com/ehazlett/atlas/api/types"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// TODO: support multiple RR per name

func (s *Server) startDNSServer() error {
	dns.HandleFunc(".", s.handler)

	u, err := url.Parse(s.cfg.BindAddress)
	if err != nil {
		return err
	}

	var proto string
	switch u.Scheme {
	case "tcp":
		proto = "tcp4"
	case "udp":
		proto = "udp4"
	default:
		return fmt.Errorf("unsupported address protocol: %s; expected tcp or udp", u.Scheme)
	}

	srv := &dns.Server{
		Addr: u.Host,
		Net:  proto,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Errorf("error starting dns server on %s: %s", s.cfg.BindAddress, err)
		}
	}()

	return nil
}

func (s *Server) handler(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.RecursionAvailable = true

	query := m.Question[0].Name
	queryType := m.Question[0].Qtype

	logrus.Debugf("nameserver: query=%q", query)
	name := getName(query, queryType)

	logrus.Infof("nameserver: looking up %s", name)
	resp, err := s.Lookup(context.Background(), &api.LookupRequest{
		Query: name,
	})
	if err != nil {
		logrus.Error(errors.Wrapf(err, "nameserver: error performing lookup for %s", name))
		w.WriteMsg(m)
		return
	}

	logrus.Debugf("lookup results: %+v", resp.Records)
	// forward if empty
	if len(resp.Records) == 0 {
		logrus.WithFields(logrus.Fields{
			"query":    name,
			"upstream": s.cfg.UpstreamDNSAddr,
		}).Debug("forwarding")
		x, err := dns.Exchange(r, s.cfg.UpstreamDNSAddr)
		if err != nil {
			logrus.Errorf("nameserver: error forwarding lookup: %+v", err)
			w.WriteMsg(m)
			return
		}
		x.SetReply(r)
		w.WriteMsg(x)
		return
	}

	// defer WriteMsg to ensure a response
	defer w.WriteMsg(m)

	// cache
	started := time.Now()
	c := s.cache.Get(name)
	if c == nil {
		if err := s.cacheRecords(name, resp.Records); err != nil {
			logrus.WithField("name", name).Warn("error caching results")
		}
		c = s.cache.Get(name)
	}

	m.Answer = []dns.RR{}
	m.Extra = []dns.RR{}
	ttl := uint32(c.TTL.Seconds() + 1)
	for _, record := range resp.Records {
		var rr dns.RR
		switch record.Type {
		case api.RecordType_A:
			ip := net.ParseIP(string(record.Value))
			rr = &dns.A{
				Hdr: dns.RR_Header{
					Name:   fqdn(name),
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    ttl,
				},
				A: ip,
			}
		case api.RecordType_CNAME:
			rr = &dns.CNAME{
				Hdr: dns.RR_Header{
					Name:   fqdn(name),
					Rrtype: dns.TypeCNAME,
					Class:  dns.ClassINET,
					Ttl:    ttl,
				},
				Target: fqdn(string(record.Value)),
			}
			// recurse to resolve name to A and add
			resp, err := s.Lookup(context.Background(), &api.LookupRequest{
				Query: string(record.Value),
			})
			if err != nil {
				logrus.Error(errors.Wrapf(err, "nameserver: error performing recursive lookup for %s", record.Value))
				w.WriteMsg(m)
				return
			}

			logrus.Debugf("looking up A records for cname %s: %+v", record.Value, resp.Records)
			for _, r := range resp.Records {
				if r.Type != api.RecordType_A {
					continue
				}
				ip := net.ParseIP(string(r.Value))
				m.Answer = append(m.Answer, &dns.A{
					Hdr: dns.RR_Header{
						Name:   fqdn(string(r.Name)),
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    ttl,
					},
					A: ip,
				})
			}
		case api.RecordType_TXT:
			rr = &dns.TXT{
				Hdr: dns.RR_Header{
					Name:   fqdn(name),
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    ttl,
				},
				Txt: []string{string(record.Value)},
			}
		case api.RecordType_MX:
			rr = &dns.MX{
				Hdr: dns.RR_Header{
					Name:   fqdn(name),
					Rrtype: dns.TypeMX,
					Class:  dns.ClassINET,
					Ttl:    ttl,
				},
				Mx: fqdn(string(record.Value)),
			}
		case api.RecordType_SRV: // srv is unique due to the return format
			v, err := typeurl.UnmarshalAny(record.Options)
			if err != nil {
				logrus.Errorf("ns: unmarshalling record options: %s", err)
				return
			}
			o, ok := v.(*types.SRVOptions)
			if !ok {
				logrus.Error("ns: invalid type for record options; expected SRVOptions")
			}
			rr = &dns.SRV{
				Hdr: dns.RR_Header{
					Name:   formatSRV(name, o),
					Rrtype: dns.TypeSRV,
					Class:  dns.ClassINET,
					Ttl:    ttl,
				},
				Target:   query,
				Priority: o.Priority,
				Weight:   o.Weight,
				Port:     o.Port,
			}
		default:
			logrus.Errorf("nameserver: unsupported record type %s for %s", record.Type, name)
		}

		lookupDuration := time.Since(started)
		logrus.Debugf("lookup duration: %s", lookupDuration)

		// set for answer or extra
		if rr.Header().Rrtype == queryType || rr.Header().Rrtype == dns.TypeCNAME {
			m.Answer = append(m.Answer, rr)
		} else {
			m.Extra = append(m.Extra, rr)
		}
	}
}

func getName(query string, queryType uint16) string {
	// adjust lookup for srv
	if queryType == dns.TypeSRV {
		p := strings.Split(query, ".")
		v := strings.Join(p[2:], ".")
		return v[:len(v)-1]
	}
	return query[:len(query)-1]
}

func formatSRV(name string, opts *types.SRVOptions) string {
	return fmt.Sprintf("_%s._%s.%s.", opts.Service, opts.Protocol, name)
}

func fqdn(name string) string {
	return name + "."
}
