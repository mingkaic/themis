package server

import (
	"fmt"
	"io"
)

const storageAPI = `PDP Storage Control API

GET /search?query=<id> [get id of storage node and return the path to it]
GET /storage/<start path...>?depth=<depth> [dump storage subtree of input depth at path]
`

// GetPath used by endpoint querying PolicyStorage
// get path of object with input id
func (s *Server) GetPath(out io.Writer, ID string) error {
	if s.p == nil {
		return fmt.Errorf("Server missing policy storage")
	}
	root, err := s.p.GetByPath("")
	if err != nil {
		return err
	}
	cb := root.MarshalPath(ID)
	if cb == nil {
		return fmt.Errorf("ID %s not found", ID)
	}
	return cb(out)
}

// DumpPath used by endpoint querying PolicyStorage
// dumps everything within input depth of object at input path
func (s *Server) DumpPath(out io.Writer, path string, depth int) error {
	if s.p == nil {
		return fmt.Errorf("Server missing policy storage")
	}
	root, err := s.p.GetByPath(path)
	if err != nil {
		return err
	}
	return root.MarshalWithDepth(out, depth)
}
