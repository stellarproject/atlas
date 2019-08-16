package atlas

// go modules are not currently detecting dependencies that are defined in
// .proto files, this is a workaround to allow them to be defined and vendored

import (
	_ "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/gogoproto"
	_ "github.com/gogo/protobuf/proto"
)
