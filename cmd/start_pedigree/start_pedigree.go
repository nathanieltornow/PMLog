package main

import (
	"flag"
	"github.com/nathanieltornow/PMLog/pedigree"
	pb "github.com/nathanieltornow/PMLog/pedigree/pedigreepb"
	"log"
	"strings"
)

var (
	ip       = flag.String("ip", "", "")
	parentIP = flag.String("parIP", "", "")
	isLeader = flag.Bool("leader", false, "")
)

func main() {
	flag.Parse()

	ped, err := pedigree.NewNode(*isLeader)
	if err != nil {
		log.Fatalln(err)
	}
	nodeIP, nodePort := getIPandPort(*ip)
	nodeInfo := pb.NodeInfo{IP: nodeIP, Port: nodePort}
	if *parentIP != "" {
		parentNodeIP, parentNodePort := getIPandPort(*parentIP)
		parentNodeInfo := pb.NodeInfo{IP: parentNodeIP, Port: parentNodePort}
		err = ped.Start(&nodeInfo, &parentNodeInfo)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		err = ped.Start(&nodeInfo, nil)
		if err != nil {
			log.Fatalln(err)
		}
	}

}

func getIPandPort(ipAddr string) (string, string) {
	splits := strings.Split(ipAddr, ":")
	return splits[0], splits[1]
}
