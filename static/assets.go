package static

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed *
var StaticFS embed.FS
var FileSystem http.FileSystem

func init() {
	subFs, err := fs.Sub(StaticFS, "static")
	if err != nil {
		log.Fatal("静态资源加载失败", err)
	}
	FileSystem = http.FS(subFs)
}

//type StaticRoute struct {
//}
//
//func (s StaticRoute) Setup(router *mux.Router) {
//	httpserver.RouterUtil.AddNoAuthPrefix("/")
//	httpserver.RouterUtil.AddNoAuthPrefix("static")
//
//	router.Handle("/favicon.ico", http.FileServer(FileSystem)).Methods(http.MethodGet, http.MethodOptions)
//	router.PathPrefix("/").Handler(util.MakeHTTPGzipHandler(http.StripPrefix("/", http.FileServer(FileSystem)))).Methods(http.MethodGet, http.MethodOptions)
//}
//
//func NewRoute() ihttpserver.IRoute {
//	opt := &StaticRoute{}
//	return opt
//}
