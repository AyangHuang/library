package concurrent

import (
	"fmt"
	"testing"
)

func TestGetBatchSliceByCount(t *testing.T) {
	batch := GetBatchSliceByCount(1, 12, 5)
	fmt.Println(batch)
}

func TestGetBatchSliceByPer(t *testing.T) {
	batch := GetBatchSliceByPer(1, 12, 5)
	fmt.Println(batch)
}
