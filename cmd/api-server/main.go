package main

import (
	"flag"
	"fmt"
	"github.com/yametech/echoer/pkg/api"
	"github.com/yametech/echoer/pkg/controller"
	"github.com/yametech/echoer/pkg/storage/mongo"
	"time"
)

var storageUri string

func main() {
	flag.StringVar(&storageUri, "storage_uri", "mongodb://127.0.0.1:27017/admin", "-storage_uri mongodb://127.0.0.1:27017/admin")
	flag.Parse()

	fmt.Println(fmt.Sprintf("echoer api server start...,%v", time.Now()))
	stage, err, errC := mongo.NewMongo(storageUri)
	if err != nil {
		panic(fmt.Sprintf("can't not open storage error (%s)", err))
	}
	server := api.NewServer(stage)

	errChan := make(chan error)

	go func() {
		if err := server.RpcServer(":8081"); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := server.Run(":8080"); err != nil {
			errChan <- err
		}
	}()

	go controller.GC(stage)
	for {
		select {
		case e := <-errC:
			panic(e)
		case e := <-errChan:
			panic(e)
		}
	}

}
