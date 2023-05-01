package main

import (
	// "encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/chaosknight/bitnet/entity"
	"github.com/gorilla/mux"
)

// var Insts_ = []*entity.FinesseConfig{
// 	&entity.FinesseConfig{
// 		InstID:      "DOGE-USDT-SWAP",
// 		Firstcount:  10,
// 		Maxcount:    2000,
// 		Addposition: -10,
// 		Maxloss:     2,
// 	},
// 	&entity.FinesseConfig{
// 		InstID:      "OP-USDT-SWAP",
// 		Firstcount:  1,
// 		Maxcount:    100,
// 		Addposition: -10,
// 		Maxloss:     2,
// 	},
// }
// var Lginuser_ = []string{
// 	"541bba76-8d22-45a4-b370-91d763d58055",
// 	"Chaos_256398",
// 	"92A8CABC7A56D45E9FE4B2296FC601A5",
// }

func main() {
	// cnf := entity.Sysconfig{
	// 	User: [3]string{"541bba76-8d22-45a4-b370-91d763d58055",
	// 		"Chaos_256398",
	// 		"92A8CABC7A56D45E9FE4B2296FC601A5"},
	// 	Insts: []*entity.FinesseConfig{
	// 		&entity.FinesseConfig{
	// 			InstID:      "DOGE-USDT-SWAP",
	// 			Firstcount:  10,
	// 			Maxcount:    2000,
	// 			Addposition: -10,
	// 			Maxloss:     2,
	// 		},
	// 		&entity.FinesseConfig{
	// 			InstID:      "OP-USDT-SWAP",
	// 			Firstcount:  1,
	// 			Maxcount:    100,
	// 			Addposition: -10,
	// 			Maxloss:     2,
	// 		},
	// 	},
	// }

	// jsonstr, err := json.Marshal(cnf)
	// if err == nil {
	// 	log.Println(string(jsonstr))
	// }

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
