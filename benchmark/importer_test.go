package benchmark

import (
	"gahenna/internal/importer"
	"os"
	"log"
	"testing"
)

func BenchmarkImport(b *testing.B) {
    const path = "../assets/customers_dev.csv"

    info, err := os.Stat(path)
    if err != nil {
        b.Fatal(err)
    }

    b.SetBytes(info.Size())
    b.ReportAllocs()
    b.ResetTimer()

    file, err := os.Open(path)
    if err != nil {
        log.Fatalf("Failed to open CSV file: %v", err)
    }
    defer file.Close()


    for i := 0; i < b.N; i++ {
		lineChan := make(chan []string, 1)
		go importer.ImportCSV(file, lineChan)
		for _ = range lineChan {

		}

	}
}
