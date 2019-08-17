package atlas

import "time"

// Config is the configuration used for the stellar server
type Config struct {
	// BindAddress is the address on which the DNS server will bind
	BindAddress string
	// Datastore is the uri to the preferred datastore backend
	Datastore string
	// GRPCAddress is the address for the grpc server
	GRPCAddress string
	// TLSCertificate is the certificate used for grpc communication
	TLSServerCertificate string
	// TLSKey is the key used for grpc communication
	TLSServerKey string
	// TLSClientCertificate is the client certificate used for communication
	TLSClientCertificate string
	// TLSClientKey is the client key used for communication
	TLSClientKey string
	// TLSInsecureSkipVerify disables certificate verification
	TLSInsecureSkipVerify bool
	// UpstreamDNSAddr is the address to use for external queries
	UpstreamDNSAddr string
	// CacheTTL is the duration for caching dns lookups
	CacheTTL time.Duration
}
