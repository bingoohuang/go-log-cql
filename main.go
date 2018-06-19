package main

import (
	"flag"
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/gorilla/mux"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	contextPath      *string
	port             string
	devMode          *bool
	cassandraCluster *string

	authParam go_utils.MustAuthParam
)

func init() {
	contextPath = flag.String("contextPath", "", "context path")
	cassandraCluster = flag.String("cassandraCluster", "127.0.0.1:9042", "cassandra cluster, like: 192.168.0.1:9042 192.168.0.2:9042,")
	httpPortArg := flag.Int("port", 7679, "Port to serve.")
	devMode = flag.Bool("devMode", false, "devMode(disable js/css minify)")

	go_utils.PrepareMustAuthFlag(&authParam)

	flag.Parse()

	port = strconv.Itoa(*httpPortArg)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc(*contextPath+"/favicon.png", go_utils.ServeFavicon("res/favicon.png", MustAsset, AssetInfo))
	handleFunc(r, "/", serveHome, false)
	handleFunc(r, "/{logid}", serveLog, false)
	http.Handle(*contextPath+"/", r)

	fmt.Println("start to listen at ", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleFunc(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request), requiredGzip bool) {
	wrap := go_utils.DumpRequest(f)
	wrap = go_utils.MustAuth(wrap, authParam)

	if requiredGzip {
		wrap = go_utils.GzipHandlerFunc(wrap)
	}

	r.HandleFunc(*contextPath+path, wrap)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	textHtml(w)

	index := string(MustAsset("res/index.html"))
	index = strings.Replace(index, "${contextPath}", *contextPath, -1)
	index = strings.Replace(index, `<ContextLogs/>`, `Please input LogId to query!`, -1)
	index = replaceIndex(index, &EventLogException{})
	index = strings.Replace(index, `<Error/>`, ``, -1)

	index = go_utils.MinifyHtml(index, true)
	w.Write([]byte(index))
}

func serveLog(w http.ResponseWriter, r *http.Request) {
	textHtml(w)

	vars := mux.Vars(r)
	logid := vars["logid"]

	index := string(MustAsset("res/index.html"))
	index = strings.Replace(index, "${contextPath}", *contextPath, -1)

	log, err := findLog(logid)
	if log != nil {
		index = replaceIndex(index, log)
		index = strings.Replace(index, `<Error/>`, ``, -1)
	} else if err != nil {
		index = strings.Replace(index, `<Error/>`, html.EscapeString(err.Error()), -1)
		index = replaceIndex(index, &EventLogException{})
	} else {
		index = strings.Replace(index, `<Error/>`, `LogId=`+logid+` Not Found!`, -1)
		index = replaceIndex(index, &EventLogException{})
	}

	index = go_utils.MinifyHtml(index, true)
	w.Write([]byte(index))
}

func textHtml(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func replaceIndex(index string, log *EventLogException) string {
	index = strings.Replace(index, `<LogId/>`, log.LogId, -1)
	index = strings.Replace(index, `<Hostname/>`, log.Hostname, -1)
	index = strings.Replace(index, `<Logger/>`, log.Logger, -1)
	index = strings.Replace(index, `<Tcode/>`, log.Tcode, -1)
	index = strings.Replace(index, `<Tid/>`, log.Tid, -1)
	index = strings.Replace(index, `<ExceptionNames/>`, html.EscapeString(log.ExceptionNames), -1)
	index = strings.Replace(index, `<Timestamp/>`, log.Timestamp, -1)
	index = strings.Replace(index, `<ContextLogs/>`, html.EscapeString(log.ContextLogs), -1)
	return index
}
