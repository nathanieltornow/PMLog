package main

import (
	"flag"
	"github.com/nathanieltornow/PMLog/pedigree"
	pb "github.com/nathanieltornow/PMLog/pedigree/pedigreepb"
	"log"
	"strings"
)

var (
	parentIP = flag.String("parIP", "", "")
)

func main() {
	flag.Parse()
	waitC := make(chan bool)
	if *parentIP != "" {
		parentNodeIP, parentNodePort := getIPandPort(*parentIP)
		parentNodeInfo := pb.NodeInfo{IP: parentNodeIP, Port: parentNodePort}
		_, err := pedigree.NewClient([]*pb.NodeInfo{&parentNodeInfo})
		if err != nil {
			log.Fatalln(err)
		}
	}
	<-waitC

}

func getIPandPort(ipAddr string) (string, string) {
	splits := strings.Split(ipAddr, ":")
	return splits[0], splits[1]
}
