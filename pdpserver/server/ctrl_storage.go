package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
)

const (
	missingStorageMsg = `"Server missing policy storage"`
	evalRuleMsg       = `"Attempting to evaluable rule"`
	queryAPI          = `Query debug API:

GET /storage/<path/to/node>?depth=<depth>
POST /storage/<path/to/node> -d '{attributes:[
	{
		"id":"<attr_id>",
		"type":"<attr_type>",
		"value":"<attr_val>"
	}
]}'
`
)

type storageHandler struct {
	s *Server
}

func handleServerErr(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, strconv.Quote(err.Error()), 500)
	}
}

func (handler *storageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.FieldsFunc(r.URL.Path, func(c rune) bool { return c == '/' })
	if len(path) == 0 {
		_, err := w.Write([]byte(queryAPI))
		handleServerErr(w, err)
		return
	}

	switch path[0] {
	case "storage":
		switch r.Method {
		case "GET":
			// dump subtree at path
			handler.storageQuery(w, r, path[1:])
			return
		case "POST":
			// calculate response of policy(set) at path
			handler.storageEval(w, r, path[1:])
			return
		}
	}
	http.Error(w, "API not found", 404)
}

func (handler *storageHandler) storageQuery(w http.ResponseWriter, r *http.Request, path []string) {
	var (
		depth uint64
		err   error
	)

	// parse depth
	queryOpt := r.URL.Query()
	if depthOpt, ok := queryOpt["depth"]; ok {
		depthStr := depthOpt[0]
		depth, err = strconv.ParseUint(depthStr, 10, 31)
		if err != nil {
			http.Error(w, strconv.Quote(err.Error()), 400)
			return
		}
	}

	storage := handler.s.p
	if storage == nil {
		http.Error(w, missingStorageMsg, 404)
		return
	}

	target, err := storage.GetAtPath(path)
	if err != nil {
		var errCode int
		if _, ok := err.(*pdp.PathNotFoundError); ok {
			errCode = 404
		} else {
			errCode = 500
		}
		http.Error(w, strconv.Quote(err.Error()), errCode)
		return
	}

	// dump
	handleServerErr(w, target.MarshalWithDepth(w, int(depth)))
}

func (handler *storageHandler) storageEval(w http.ResponseWriter, r *http.Request, path []string) {
	server := handler.s

	// parse context
	decoder := json.NewDecoder(r.Body)
	var req pb.Request
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Context format error: %+v", err), 400)
		return
	}
	ctx, err := server.newContext(server.c, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Context error: %+v", err), 400)
		return
	}

	storage := server.p
	if storage == nil {
		http.Error(w, missingStorageMsg, 404)
		return
	}

	target, err := storage.GetAtPath(path)
	if err != nil {
		var errCode int
		if _, ok := err.(*pdp.PathNotFoundError); ok {
			errCode = 404
		} else {
			errCode = 500
		}
		http.Error(w, strconv.Quote(err.Error()), errCode)
		return
	}
	if e, ok := target.(pdp.Evaluable); ok {
		// calculate and dump response
		res := e.Calculate(ctx)
		b, err := json.Marshal(res)
		if err == nil {
			_, err = w.Write(b)
		}

		handleServerErr(w, err)
	} else {
		http.Error(w, evalRuleMsg, 404)
	}
}
