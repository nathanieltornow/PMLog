package main

import (
	"flag"
)

var (
	IP      = flag.String("IP", "", "")
	peerIPs = flag.String("peers", "", "")
	id      = flag.Int("id", 0, "")
)

func main() {
	flag.Parse()
	//primLog, err := mem_log.NewMemLog()
	//if err != nil {
	//	logrus.Fatalln(err)
	//}
	//secLog, err := mem_log.NewMemLog()
	//if err != nil {
	//	logrus.Fatalln(err)
	//}
	//node, err := app_node.NewNode(uint32(*id), 0, primLog, secLog)
	//if err != nil {
	//	logrus.Fatalln(err)
	//}
	//
	//var peerList []string
	//if *peerIPs != "" {
	//	peerList = strings.Split(*peerIPs, ",")
	//}
	//err = node.Start(*IP, peerList)
	//if err != nil {
	//	logrus.Fatalln(err)
	//}
}
