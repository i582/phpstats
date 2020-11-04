package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/i582/phpstats/internal/representator"
	"github.com/i582/phpstats/internal/stats/walkers"
)

// RunServer starts the server on the given port.
func RunServer(port int64) {
	if port == 0 {
		port = 8080
	}

	http.HandleFunc("/info/class", infoClassHandler)
	http.HandleFunc("/info/func", infoFunctionHandler)
	http.HandleFunc("/info/namespace", infoNamespaceHandler)
	http.HandleFunc("/exit", exitHandler)
	http.HandleFunc("/analyzeStats", analyzeStatsHandler)

	go http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
}

func infoClassHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	classNames, err := walkers.GlobalCtx.Classes.GetFullClassName(name)
	if err != nil {
		data, _ := representator.GetJsonClassReprWithFlag(nil)
		fmt.Fprintln(w, data)
		return
	}

	var className string

	if len(classNames) > 1 {
		className = classNames[0]
	} else {
		className = classNames[0]
	}

	class, _ := walkers.GlobalCtx.Classes.Get(className)

	data, _ := representator.GetJsonClassReprWithFlag(class)
	fmt.Fprintln(w, data)
}

func infoFunctionHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	funcNameKeys, err := walkers.GlobalCtx.Functions.GetFullFuncName(name)
	if err != nil {
		data, _ := representator.GetJsonFunctionReprWithFlag(nil)
		fmt.Fprintln(w, data)
		return
	}

	var funcKeyIndex int

	if len(funcNameKeys) > 1 {
		funcKeyIndex = 0
	} else {
		funcKeyIndex = 0
	}

	fn, _ := walkers.GlobalCtx.Functions.Get(funcNameKeys[funcKeyIndex])

	data, _ := representator.GetJsonFunctionReprWithFlag(fn)
	fmt.Fprintln(w, data)
}

func infoNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	ns, ok := walkers.GlobalCtx.Namespaces.GetNamespace(name)
	if !ok {
		data, _ := representator.GetJsonNamespaceReprWithFlag(nil)
		fmt.Fprintln(w, data)
		return
	}

	data, _ := representator.GetJsonNamespaceReprWithFlag(ns)
	fmt.Fprintln(w, data)
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

func analyzeStatsHandler(w http.ResponseWriter, r *http.Request) {
	if walkers.GlobalCtx.BarLinting == nil {
		fmt.Fprintf(w, "{\"state\": \"indexing\", \"current\": 0.0}")
		return
	}

	count := float64(walkers.GlobalCtx.BarLinting.Total())
	cur := float64(walkers.GlobalCtx.BarLinting.Current())

	fmt.Fprintf(w, "{\"state\": \"linting\", \"current\": %f}", cur/count)
}
