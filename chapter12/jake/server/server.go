package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"go-network-programming/chapter12/jake/housework/v1"
	"log"
	"net"

	"google.golang.org/grpc"
)

var addr, certFn, keyFn string

func init() {
	flag.StringVar(&addr, "address", "localhost:34443",
		"listen address")
	flag.StringVar(&certFn, "cert", "cert.pem", "certificate file")
	flag.StringVar(&keyFn, "key", "key.pem", "private key file")
}

func main() {
	flag.Parse()

	server := grpc.NewServer()
	rosie := new(Rosie)
	housework.RegisterRobotMaidService(server, rosie.Service())

	cert, err := tls.LoadX509KeyPair(certFn, keyFn)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening for TLS connections on %s ...", addr)
	log.Fatal(server.Serve(tls.NewListener(listener,
		&tls.Config{
			Certificates:             []tls.Certificate{cert},
			CurvePreferences:         []tls.CurveID{tls.CurveP256},
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		},
	)))
}
