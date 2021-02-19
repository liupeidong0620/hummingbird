package tunnel

import (
	"github.com/liupeidong0620/hummingbird/adapter"
)

func generateNATKey(m *adapter.Metadata) string {
	return m.SourceAddress()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
