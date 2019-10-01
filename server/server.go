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
	"context"
	"io/ioutil"
	"runtime"
	"runtime/pprof"

	ptypes "github.com/gogo/protobuf/types"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stellarproject/atlas"
	"google.golang.org/grpc"

	api "github.com/stellarproject/atlas/api/v1"
)

const (
	recordKey  = "terra.dns.records.%s.%s"
	publishKey = "terra.dns.update"
)

var (
	empty = &ptypes.Empty{}
)

// Server is an Atlas server
type Server struct {
	cfg  *atlas.Config
	pool *redis.Pool
}

// NewServer returns a new server
func NewServer(cfg *atlas.Config) (*Server, error) {
	pool := redis.NewPool(func() (redis.Conn, error) {
		conn, err := redis.DialURL(cfg.RedisURL)
		if err != nil {
			return nil, errors.Wrap(err, "unable to connect to redis")
		}
		// TODO
		//if auth != "" {
		//	if _, err := conn.Do("AUTH", auth); err != nil {
		//		conn.Close()
		//		return nil, errors.Wrap(err, "unable to authenticate to redis")
		//	}
		//}
		return conn, nil
	}, 10)

	srv := &Server{
		cfg:  cfg,
		pool: pool,
	}

	return srv, nil
}

// Register enables callers to register this service with an existing GRPC server
func (s *Server) Register(server *grpc.Server) error {
	logrus.Debug("registering nameserver with grpc")
	api.RegisterNameserverServer(server, s)
	return nil
}

// Start starts the embedded DNS server
func (s *Server) Start() error {
	if err := s.update(context.Background()); err != nil {
		return err
	}
	return s.startListener()
}

func (s *Server) startListener() error {
	// start listener for pub/sub
	errCh := make(chan error, 1)
	go func() {
		c := s.pool.Get()
		defer c.Close()

		psc := redis.PubSubConn{Conn: c}
		psc.Subscribe(publishKey)
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				// TODO: update
				if err := s.update(context.Background()); err != nil {
					logrus.Error(err)
					continue
				}
			case redis.Subscription:
			default:
				logrus.Debugf("unknown message type %T", v)
			}
		}
	}()

	err := <-errCh
	return err
}

// Stop is used to stop and release resources
func (s *Server) Stop() error {
	if s.pool != nil {
		s.pool.Close()
	}
	return nil
}

// GenerateProfile generates a new Go profile
func (s *Server) GenerateProfile() (string, error) {
	tmpfile, err := ioutil.TempFile("", "atlas-profile-")
	if err != nil {
		return "", err
	}
	runtime.GC()
	if err := pprof.WriteHeapProfile(tmpfile); err != nil {
		return "", err
	}
	tmpfile.Close()
	return tmpfile.Name(), nil
}

func (s *Server) do(ctx context.Context, cmd string, args ...interface{}) (interface{}, error) {
	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	r, err := conn.Do(cmd, args...)
	return r, err
}
