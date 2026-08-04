package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/Norzuiso/client/internal/handler"
	"github.com/Norzuiso/client/internal/states"

	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:8080", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

/*
TODO
==================================================
1. Check how to define localaddress and port


*/

func main() {
	clientName := os.Args[1]
	initState := states.GetStateByArgs(os.Args[2])
	idStr := os.Args[3]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Fail to dial: %v", err)
	}
	defer conn.Close()
	clientBroadcast := pb.NewBroadcastClient(conn)

	ctx := context.Background()
	msgHandler := handler.NewMsgHandler(ctx, clientBroadcast)

	msgHandler.Init(clientName, initState, id)

	// ==========================================
	/*
				1. Connect client -> Call ConnectClient from client
				1.1. Get the response from the orchestrator and store it
		 		2. Send first msg to open stream -> MessageType 5



	*/
	// ==========================================
}
