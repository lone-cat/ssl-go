package storage

func bytesEqual(data1 []byte, data2 []byte) bool {
	return string(data1) == string(data2)
}

func bytesArrEqual(data1 [][]byte, data2 [][]byte) bool {
	if len(data1) != len(data2) {
		return false
	}
	for i, dat1 := range data1 {
		if !bytesEqual(dat1, data2[i]) {
			return false
		}
	}
	return true
}
