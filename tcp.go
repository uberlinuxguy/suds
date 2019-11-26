package main

import (
	"fmt"
	"net/http"
	"strings"
)

// handleConnection spits out the contents of a table in JSON
func handleConnection(w http.ResponseWriter, req *http.Request) {
	var tableName string
	urlParts := strings.Split(req.URL.Path, "/")
	tableName = urlParts[2]
	tableDump, err := DumpTable(tableName)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "Internal Error: "+err.Error())
		return
	}
	fmt.Fprintln(w, tableDump)
}
