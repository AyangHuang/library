package concurrent

// GetBatchSliceByCount 根据 [startIndex, endIndex] 划分批次，分成 count 个批次
func GetBatchSliceByCount(start, end, count int) [][]int {
	var (
		size           = end - start + 1
		batch          = make([][]int, 0, count)
		batchPer       = size / count
		batchRemainder = size % count
	)

	for i := 0; i < batchRemainder; i++ {
		batch = append(batch, []int{start, start + batchPer})
		start += batchPer + 1
	}
	for start <= end {
		batch = append(batch, []int{start, start + batchPer - 1})
		start += batchPer
	}
	return batch
}

// GetBatchSliceByPer 根据 [start, end] 划分批次，每个批次 per 个，进行分批次
func GetBatchSliceByPer(start, end, per int) [][]int {
	if start >= end || per <= 0 {
		return nil
	}
	batch := make([][]int, 0, (end-start+1)/per+1)
	for start <= end {
		innerEnd := start + per - 1
		if innerEnd > end {
			innerEnd = end
		}
		batch = append(batch, []int{start, innerEnd})
		start = innerEnd + 1
	}
	return batch
}
