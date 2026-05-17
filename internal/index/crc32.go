package index

import (
	"hash/crc32"
	"os"
)

// fileCRC32 calculates file CRC32
func fileCRC32(filePath string) uint32 {
	fbytes, err := os.ReadFile(filePath)
	if err != nil {
		return 0
	}
	return crc32.ChecksumIEEE(fbytes)
}
