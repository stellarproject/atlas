/*
   Copyright 2019 Stellar Project

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
	"fmt"
	"net/http"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const (
	namespace = "atlas"
	subsystem = "dns"
)

var (
	lookupACounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "lookup_a_total",
			Help:      "Total number of A record lookups",
		},
	)
	lookupCNAMECounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "lookup_cname_total",
			Help:      "Total number of CNAME record lookups",
		},
	)
	lookupForwardCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "lookup_forward_total",
			Help:      "Total number of upstream lookups",
		},
	)
	queryDurationHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "query_milliseconds",
			Help:      "Duration of query in milliseconds",
			Buckets: []float64{
				1,
				5,
				10,
				25,
				50,
				100,
				250,
				500,
				1000,
			},
		},
	)
	createRecordCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "create_total",
			Help:      "Total number of record creates",
		},
	)
	deleteRecordCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "delete_total",
			Help:      "Total number of record deletes",
		},
	)
)

func init() {
	prometheus.MustRegister(lookupACounter)
	prometheus.MustRegister(lookupCNAMECounter)
	prometheus.MustRegister(lookupForwardCounter)
	prometheus.MustRegister(queryDurationHistogram)
	prometheus.MustRegister(createRecordCounter)
	prometheus.MustRegister(deleteRecordCounter)
}

func (s *Server) startMetricsServer() error {
	u, err := url.Parse(s.cfg.MetricsAddr)
	if err != nil {
		return err
	}
	if u.Scheme != "http" {
		return fmt.Errorf("metrics: only http endpoints are supported")
	}
	logrus.WithFields(logrus.Fields{
		"addr": s.cfg.MetricsAddr,
	}).Debug("starting metrics server")

	// start emitter listeners
	s.startMetricListener()

	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(u.Host, nil)
}

func (s *Server) startMetricListener() {
	go s.metricListenerLookupATotal()
	go s.metricListenerLookupCNAMETotal()
	go s.metricListenerLookupForwardTotal()
	go s.metricListenerQueryDuration()
	go s.metricListenerCreateRecord()
	go s.metricListenerDeleteRecord()
}

func (s *Server) metricListenerLookupATotal() {
	for range s.emitter.On(emitLookupA) {
		lookupACounter.Inc()
	}
}

func (s *Server) metricListenerLookupCNAMETotal() {
	for range s.emitter.On(emitLookupCNAME) {
		lookupCNAMECounter.Inc()
	}
}

func (s *Server) metricListenerLookupForwardTotal() {
	for range s.emitter.On(emitLookupForward) {
		lookupForwardCounter.Inc()
	}
}

func (s *Server) metricListenerQueryDuration() {
	for evt := range s.emitter.On(emitQueryDuration) {
		queryDurationHistogram.Observe(evt.Float(0))
	}
}

func (s *Server) metricListenerCreateRecord() {
	for evt := range s.emitter.On(emitCreateRecord) {
		createRecordCounter.Add(evt.Float(0))
	}
}

func (s *Server) metricListenerDeleteRecord() {
	for range s.emitter.On(emitDeleteRecord) {
		deleteRecordCounter.Inc()
	}
}
