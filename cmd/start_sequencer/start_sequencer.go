package main

import (
	"flag"
	"github.com/nathanieltornow/PMLog/order_repl_framework/sequencer"
	"github.com/sirupsen/logrus"
)

var (
	IP     = flag.String("IP", "", "The IP on which the sequencer listens")
	parIP  = flag.String("parIP", "", "The IP of the parent-sequencer")
	root   = flag.Bool("root", false, "If the sequencer is the root-sequencer")
	leader = flag.Bool("leader", false, "If the sequencer is a leader and can therefore reply to requests")
	color  = flag.Int("color", 0, "The color which the sequencer represents")
)

func main() {
	flag.Parse()
	seq := sequencer.NewSequencer(*root, *leader, uint32(*color))
	err := seq.Start(*IP, *parIP)
	if err != nil {
		logrus.Fatalln(err)
	}
}
