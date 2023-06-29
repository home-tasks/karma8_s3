package utils

const FileParts = 6
const machineWord = 8

type Sizes = [FileParts]int64

// Pure function
func SplitSize(size int64) (parts Sizes) {
	chunks := size / (FileParts * machineWord)
	remainder := size % (FileParts * machineWord)

	for i := 0; i < FileParts; i++ {
		parts[i] = chunks * machineWord
		if remainder >= machineWord {
			parts[i] += machineWord
			remainder -= machineWord
		}
	}

	parts[len(parts)-1] += remainder

	return
}
