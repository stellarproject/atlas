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

	ptypes "github.com/gogo/protobuf/types"
	"github.com/gomodule/redigo/redis"
	api "github.com/stellarproject/atlas/api/v1"
)

// Delete removes records from Atlas
func (s *Server) Delete(ctx context.Context, req *api.DeleteRequest) (*ptypes.Empty, error) {
	keys, err := redis.Strings(s.do(ctx, "KEYS", fmt.Sprintf(recordKey, req.Name, "*")))
	if err != nil {
		return empty, err
	}
	for _, key := range keys {
		if _, err := s.do(ctx, "DEL", key); err != nil {
			return empty, err
		}
	}
	return empty, nil
}
