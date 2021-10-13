package mem_log_go

import (
	"strings"
	"testing"
)

func TestLog_Append(t *testing.T) {
	log, _ := NewMemLogGo()
	for i := 0; i < 1000; i++ {
		_ = log.Append(strings.Repeat("h", 4000), 14)
		_ = log.Commit(14, 33)
		_, _ = log.Read(33)
	}
}
