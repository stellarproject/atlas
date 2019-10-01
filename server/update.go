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
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	v1 "github.com/stellarproject/atlas/api/v1"
)

const (
	confTemplate = `# managed by atlas
port={{ .Port }}
server={{ .Upstream }}
log-queries
local=/localnet/
{{ range .Records }}
{{ getType .Type }}=/{{ .Name }}/{{ .Value }}{{ end }}
`
)

type config struct {
	Port     int
	Upstream string
	Records  []*v1.Record
}

func getType(t v1.RecordType) string {
	switch strings.ToLower(t.String()) {
	case "a":
		return "address"
	}
	return fmt.Sprintf("# unknown type %s", t.String())
}

func (s *Server) update(ctx context.Context) error {
	logrus.Infof("updating config %s", s.cfg.ConfigPath)
	t, err := template.New("config").Funcs(template.FuncMap{
		"getType": getType,
	}).Parse(confTemplate)
	if err != nil {
		return err
	}

	records, err := s.getRecords(ctx)
	if err != nil {
		return err
	}

	tmpFile, err := ioutil.TempFile("", "atlas-conf-")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	cfg := &config{
		Port:     s.cfg.Port,
		Upstream: s.cfg.UpstreamDNSAddr,
		Records:  records,
	}
	if err := t.Execute(tmpFile, cfg); err != nil {
		return err
	}
	tmpFile.Close()

	if err := os.Rename(tmpFile.Name(), s.cfg.ConfigPath); err != nil {
		return err
	}

	return nil
}
