package server

import (
	"io/ioutil"
	"runtime"
	"runtime/pprof"

	"github.com/ehazlett/atlas"
	ptypes "github.com/gogo/protobuf/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
	"github.com/ehazlett/atlas/ds"
	"github.com/ehazlett/ttlcache"
)

var (
	empty = &ptypes.Empty{}
)

type Server struct {
	cfg   *atlas.Config
	ds    ds.Datastore
	cache *ttlcache.TTLCache
}

func NewServer(cfg *atlas.Config) (*Server, error) {
	ds, err := atlas.GetDatastore(cfg.Datastore)
	if err != nil {
		return nil, err
	}
	srv := &Server{
		cfg: cfg,
		ds:  ds,
	}
	if cfg.CacheTTL != 0 {
		c, err := ttlcache.NewTTLCache(cfg.CacheTTL)
		if err != nil {
			return nil, err
		}
		srv.cache = c
	}

	return srv, nil
}

func (s *Server) Register(server *grpc.Server) error {
	logrus.Debug("registering nameserver with grpc")
	api.RegisterNameserverServer(server, s)
	return nil
}

func (s *Server) Start() error {
	return s.startDNSServer()
}

func (s *Server) Stop() error {
	return nil
}

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
