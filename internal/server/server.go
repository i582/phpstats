package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/i582/phpstats/internal/representator"
	"github.com/i582/phpstats/internal/stats"
)

func RunServer() {
	http.HandleFunc("/info/class", InfoClassHandler)
	http.HandleFunc("/info/func", InfoFunctionHandler)
	http.HandleFunc("/info/namespace", InfoNamespaceHandler)
	http.HandleFunc("/exit", ExitHandler)
	http.HandleFunc("/analyzeStats", AnalyzeStatsHandler)

	go http.ListenAndServe("localhost:8080", nil)
}

func InfoClassHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	classNames, err := stats.GlobalCtx.Classes.GetFullClassName(name)
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

	class, _ := stats.GlobalCtx.Classes.Get(className)

	data, _ := representator.GetJsonClassReprWithFlag(class)
	fmt.Fprintln(w, data)
}

func InfoFunctionHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	funcNameKeys, err := stats.GlobalCtx.Funcs.GetFullFuncName(name)
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

	fn, _ := stats.GlobalCtx.Funcs.Get(funcNameKeys[funcKeyIndex])

	data, _ := representator.GetJsonFunctionReprWithFlag(fn)
	fmt.Fprintln(w, data)
}

func InfoNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	ns, ok := stats.GlobalCtx.Namespaces.GetNamespace(name)
	if !ok {
		data, _ := representator.GetJsonNamespaceReprWithFlag(nil)
		fmt.Fprintln(w, data)
		return
	}

	data, _ := representator.GetJsonNamespaceReprWithFlag(ns)
	fmt.Fprintln(w, data)
}

func ExitHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

func AnalyzeStatsHandler(w http.ResponseWriter, r *http.Request) {
	if stats.BarLinting == nil {
		fmt.Fprintf(w, "{\"state\": \"indexing\", \"current\": 0.0}")
		return
	}

	count := float64(stats.BarLinting.Total())
	cur := float64(stats.BarLinting.Current())

	fmt.Fprintf(w, "{\"state\": \"linting\", \"current\": %f}", cur/count)
}
