package server

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"

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
	m.RecursionAvailable = true

	query := m.Question[0].Name
	queryType := m.Question[0].Qtype

	logrus.Debugf("nameserver: query=%q", query)
	name := getName(query, queryType)

	logrus.Debugf("nameserver: looking up %s", name)
	resp, err := s.Lookup(context.Background(), &api.LookupRequest{
		Query: name,
	})
	if err != nil {
		logrus.Error(errors.Wrapf(err, "nameserver: error performing lookup for %s", name))
		w.WriteMsg(m)
		return
	}
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

	m.Answer = []dns.RR{}

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
					Ttl:    0,
				},
				A: ip,
			}
		case api.RecordType_CNAME:
			rr = &dns.CNAME{
				Hdr: dns.RR_Header{
					Name:   fqdn(name),
					Rrtype: dns.TypeCNAME,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				Target: string(record.Value),
			}
		case api.RecordType_TXT:
			rr = &dns.TXT{
				Hdr: dns.RR_Header{
					Name:   fqdn(name),
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				Txt: []string{string(record.Value)},
			}
		case api.RecordType_MX:
			rr = &dns.MX{
				Hdr: dns.RR_Header{
					Name:   fqdn(name),
					Rrtype: dns.TypeMX,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				Mx: string(record.Value),
			}
		case api.RecordType_SRV: // srv is unique do to the return format
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
					Ttl:    0,
				},
				Target:   query,
				Priority: o.Priority,
				Weight:   o.Weight,
				Port:     o.Port,
			}
		default:
			logrus.Errorf("nameserver: unsupported record type %s for %s", record.Type, name)
		}

		// set for answer or extra
		if rr.Header().Rrtype == queryType {
			m.Answer = append(m.Answer, rr)
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
