package importer

import (
	"io"
	"log"
	"os"
	"encoding/csv"
)


func ImportCSV(filePath *os.File, lineChan chan <- []string) {
	defer close(lineChan)
    reader := csv.NewReader(filePath)
	for {
        record, err := reader.Read()
        if err != nil {
            if err == io.EOF {
                break
            }
            log.Fatalf("Error reading CSV data: %v", err)
        }
		lineChan <- record
    }
}
