package main

import (
	"log"
	"net/http"
	"time"

	"github.com/chaosknight/bitnet/entity"
	"github.com/gorilla/mux"
)

func main() {
	syscnf := entity.LoadCnf("./cfg.json")
	if syscnf == nil {
		log.Println("配置文件错误")
		return
	}

	apikey := entity.NewAPIKey(syscnf.User[0], syscnf.User[2], syscnf.User[1])

	syss := NewSys(apikey, syscnf.Insts)
	syss.Init()

	router := mux.NewRouter()
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
