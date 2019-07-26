package main

import (
	"net/http"
)

type coroutineFileShareObj struct {
}

// var coroutineFileShareObjSrv *http.Server

func coroutineFileShare() coroutineFileShareObj {
	return coroutineFileShareObj{}
}

func (obj coroutineFileShareObj) run(port string, sharePath string) {
	// srv := &http.Server{Addr: ":" + port}
	// coroutineFileShareObjSrv = srv
	// http.Handle("/", http.FileServer(http.Dir(sharePath)))
	// coroutineFileShareObjSrv.ListenAndServe()

	// http.Handle("/", http.FileServer(http.Dir(sharePath)))
	// log.Fatal(srv.ListenAndServe(":"+port, nil))
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(sharePath)))
	http.ListenAndServe(":"+port, mux)
}

// func (obj coroutineFileShareObj) stop() {
// 	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
// 	if err := coroutineFileShareObjSrv.Shutdown(ctx); err != nil {
// 		// handle err
// 	}
// 	log.Println("coroutineFileShareObj stop")
// }
