package server

import (
	"fmt"
	"io"
)

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
	cb, found := root.PathMarshal(ID)
	if !found {
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
	return root.DepthMarshal(out, depth)
}
