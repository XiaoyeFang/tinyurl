// Copyright © 2018 joy  <lzy@spf13.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"sync"
	"url-shortener/cmd"
	"url-shortener/config"
	"url-shortener/grpcurl"
	"url-shortener/handler"
	"url-shortener/protos"
	"url-shortener/storage/postgres"
)

func main() {

	cmd.Execute()
	glog.Flush()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		httpStart()

	}()
	go func() {

		grpcStart()
	}()
	wg.Wait()
	//

}

func httpStart() {
	// Set use storage, select [Postgres, Filesystem, Redis ...]
	svc, err := postgres.New(config.UrlConfig.Postgres.Host, config.UrlConfig.Postgres.Port, config.UrlConfig.Postgres.User,
		config.UrlConfig.Postgres.Password, config.UrlConfig.Postgres.DB)
	if err != nil {
		glog.V(0).Infof("postgres.New %v\n", err)
	}
	//defer svc.Close()

	// Create a server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.UrlConfig.Server.Host, config.UrlConfig.Server.Port),
		Handler: handler.New(config.UrlConfig.Options.Prefix, svc),
	}

	// Start server
	glog.V(0).Infof("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		glog.V(0).Infof("%v", err)
	} else {
		glog.V(0).Infoln("Server closed!")
	}
}

func grpcStart() {
	lis, err := net.Listen("tcp", config.UrlConfig.Grpclisten)
	if err != nil {
		panic(err)
	}
	glog.Errorf("am-grpc-port %v", config.UrlConfig.Grpclisten)
	server := grpc.NewServer()
	srv := grpcurl.ShortenerServer{}
	protos.RegisterUrlShortEnerServiceServer(server, &srv)
	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
