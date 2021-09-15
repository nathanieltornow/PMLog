package redis

import "strings"

const (
	seperator = ","
)

func (s *Server) prepareTransaction(txString string) {
	//txs := parseTransactionString(txString)
	//for tx := range txString {
	//
	//}

}

func parseTransactionString(txStr string) []string {
	return strings.Split(txStr, seperator)
}
