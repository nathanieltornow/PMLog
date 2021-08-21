package main

import (
	"flag"
	"github.com/nathanieltornow/PMLog/sequencer"
	"github.com/sirupsen/logrus"
)

var (
	IP     = flag.String("IP", "", "")
	parIP  = flag.String("parIP", "", "")
	root   = flag.Bool("root", false, "")
	leader = flag.Bool("leader", false, "")
	color  = flag.Int("color", 0, "")
)

func main() {
	flag.Parse()
	seq := sequencer.NewSequencer(*root, *leader, uint32(*color))
	err := seq.Start(*IP, *parIP)
	if err != nil {
		logrus.Fatalln(err)
	}
}
