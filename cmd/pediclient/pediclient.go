package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

var (
	parentIP = flag.String("parIP", "", "")
)

func main() {
	//flag.Parse()
	//waitC := make(chan bool)
	//if *parentIP != "" {
	//	parentNodeIP, parentNodePort := getIPandPort(*parentIP)
	//	parentNodeInfo := pb.NodeInfo{IP: parentNodeIP, Port: parentNodePort}
	//	_, err := pedigree.NewClient([]*pb.NodeInfo{&parentNodeInfo})
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//}
	//<-waitC
	lis, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		fmt.Println(lis, err)
	}
	lis, err = net.Listen("tcp", "localhost:9000")
	if err != nil {
		fmt.Println(lis, err)
	}
}

func getIPandPort(ipAddr string) (string, string) {
	splits := strings.Split(ipAddr, ":")
	return splits[0], splits[1]
}
