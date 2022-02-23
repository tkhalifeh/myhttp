package main

import (
	"flag"
	"fmt"
	"github.com/tkhalifeh/urlprocessor"
	"log"
	"net/http"
)

func main() {
	parallelPtr := flag.Int("parallel", 10, "number of parallel request")
	flag.Parse()
	urls := flag.Args()

	mainJob, err := urlprocessor.New(urlprocessor.Args{
		HTTPClient:    http.DefaultClient,
		ParallelLimit: *parallelPtr,
		URLs:          urls,
	})
	if err != nil {
		log.Fatalln(err)
	}

	resultCh := make(chan urlprocessor.HashResult)
	errorCh := make(chan error)
	doneCh := make(chan struct{})
	mainJob.Process(resultCh, errorCh, doneCh)

	for {
		select {
		case opResult := <-resultCh:
			fmt.Printf("%v %v\n", opResult.URL, opResult.MD5Hash)
		case opError := <-errorCh:
			fmt.Println(opError)
		case <-doneCh:
			return
		}
	}
}
