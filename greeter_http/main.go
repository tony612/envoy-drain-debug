/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
)

const (
	defaultName = "world"
	ridKey      = "x-request-id"
)

var grpcAddress = flag.String("grpc-address", "localhost:60051", "")

func main() {
	flag.Parse()

	log.Printf("connecting to grpc %s", *grpcAddress)

	// Set up a connection to the server.
	conn, err := grpc.Dial(*grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName

	http.HandleFunc("/grpc", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*600)
		defer cancel()
		if rid, ok := r.Header["x-request-id"]; ok {
			if len(rid) > 0 {
				ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", rid[0])
			}
		}
		if rid, ok := r.Header["X-Request-Id"]; ok {
			if len(rid) > 0 {
				ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", rid[0])
			}
		}

		resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "Got error %v", err.Error())
		} else {
			fmt.Fprintf(w, "Hello %v, from http server", resp.Message)
		}
	})

	log.Fatal(http.ListenAndServe(":18080", nil))
}
