package stats

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func InfoClassHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	onlyMetrics := r.URL.Query().Get("onlyMetrics") == "true"

	classNames, err := GlobalCtx.Classes.GetFullClassName(name)
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

	class, _ := GlobalCtx.Classes.Get(className)

	if onlyMetrics {
		fmt.Fprintln(w, class.OnlyMetricsString())
		return
	}

	data := class.FullString(0, true)
	data = strings.TrimLeft(strings.TrimRight(data, "\n"), "\n")

	fmt.Fprintln(w, data)
}

func InfoFunctionHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	funcNameKeys, err := GlobalCtx.Funcs.GetFullFuncName(name)
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

	fn, _ := GlobalCtx.Funcs.Get(funcNameKeys[funcKeyIndex])

	data := fn.PluginFunctionString()
	data = strings.TrimLeft(strings.TrimRight(data, "\n"), "\n")

	fmt.Fprintln(w, data)
}

func ExitHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

func AnalyzeStatsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, float64(BarLinting.Total()), float64(GlobalCtx.Files.Len()))
}
