package sort

// SelectSort -
func SelectSort(data []int) {
	for i := 0; i < len(data)-1; i++ {
		idxMin := i
		for j := i + 1; j < len(data); j++ {
			if data[j] < data[idxMin] {
				idxMin = j
			}
		}
		data[i], data[idxMin] = data[idxMin], data[i]
	}
}
