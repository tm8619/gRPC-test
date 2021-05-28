package main

import (
	"encoding/json"
	"github.com/tm8619/gRPC-vs-REST/gRPC/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"testing"
)

const (
	address = "localhost:8083"
)

func main() {

}

type getNumbersOutput struct {
	Numbers []int `json:"numbers"`
}

func BenchmarkRest(b *testing.B) {
	url := "http://localhost:8080/numbers?from=1&to=1000000"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, _ := http.Get(url)
		var result getNumbersOutput
		json.NewDecoder(resp.Body).Decode(&result)
		if len(result.Numbers) == 0 {
			log.Fatal("could not get numbers at rest")
		}
	}
}

func BenchmarkRestWithoutDecode(b *testing.B) {
	url := "http://localhost:8080/numbers?from=1&to=1000000"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := http.Get(url)
		if err != nil {
			log.Fatalf("could not get numbers at rest2 %v:", err)
		}
	}
}

func BenchmarkGrpc(b *testing.B) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := service.NewGRPCClient(conn)

	// Contact the server and print out its response.
	var from, to int64
	from, to = 1, 1000000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := c.GetNumbers(context.Background(), &service.GetNumbersInput{
			From: from,
			To:   to,
		})
		if err != nil {
			log.Fatalf("could not get numbers: %v", err)
		}
		if len(resp.Numbers) == 0 {
			log.Fatal("could not get numbers at gRPC")
		}
	}
}
