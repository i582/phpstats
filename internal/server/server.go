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
	http.HandleFunc("/exit", ExitHandler)
	http.HandleFunc("/analyzeStats", AnalyzeStatsHandler)

	go http.ListenAndServe("localhost:8000", nil)
}

func InfoClassHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	classNames, err := stats.GlobalCtx.Classes.GetFullClassName(name)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	var className string

	if len(classNames) > 1 {
		className = classNames[0]
	} else {
		className = classNames[0]
	}

	class, _ := stats.GlobalCtx.Classes.Get(className)

	data, _ := representator.GetJsonClassRepr(class)
	fmt.Fprintln(w, data)
}

func InfoFunctionHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	funcNameKeys, err := stats.GlobalCtx.Funcs.GetFullFuncName(name)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	var funcKeyIndex int

	if len(funcNameKeys) > 1 {
		funcKeyIndex = 0
	} else {
		funcKeyIndex = 0
	}

	fn, _ := stats.GlobalCtx.Funcs.Get(funcNameKeys[funcKeyIndex])

	data, _ := representator.GetJsonFunctionRepr(fn)
	fmt.Fprintln(w, data)
}

func ExitHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

func AnalyzeStatsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, float64(stats.BarLinting.Total())/float64(stats.GlobalCtx.Files.Len()))
}
