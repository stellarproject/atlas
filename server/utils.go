package server

import (
	"crypto/sha256"
	"fmt"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
)

func getRecordID(r *api.Record) string {
	h := sha256.New()
	h.Write([]byte(r.Type.String()))
	h.Write([]byte(r.Name))
	h.Write([]byte(r.Value))
	return fmt.Sprintf("%x", h.Sum(nil))
}
