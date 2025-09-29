package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	var totalMemory uint64 = 2 << 30 // 2GB, you can increase this value for more memory usage
	sizeOfInt := uint64(8)           // size of int in bytes on most systems

	// Calculate the number of integers needed to allocate 'totalMemory' bytes
	numIntegers := totalMemory / sizeOfInt

	// Attempt to create a large slice
	hugeSlice := make([]int, numIntegers)

	if hugeSlice == nil {
		log.Info("Memory allocation failed.")
		os.Exit(1)
	}

	// Fill the slice to ensure memory is actually used
	for i := range hugeSlice {
		hugeSlice[i] = i
		log.Info(i)
		time.Sleep(1 * time.Nanosecond)
	}

	log.Info("Memory allocated and filled successfully.")
	log.Infof("Size of slice: %d bytes\n", totalMemory)

	// Keep the program running so memory is not released immediately
	select {}
}
