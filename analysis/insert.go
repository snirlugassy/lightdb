package main

import (
	"encoding/csv"
	"fmt"
	"lightdb"
	"log"
	"os"
	"path"
	"reflect"
	"time"
)

type Dummy struct {
	XX int
	XY int
	YX string
	YY string
	ZZ bool
	ZX byte
	XZ rune
}

func analyzeInsert(output_path string) {
	runtimes := make(map[int]float64)

	for i := 1; i <= 30; i++ {
		collection := lightdb.Collection{
			Name:     "analysis",
			FilePath: path.Join(os.TempDir(), "analysis.collection"),
			DType:    reflect.TypeOf(Dummy{}),
		}

		sample_size := i * 10000
		dummy_data := make([]interface{}, 0)
		for j := 0; j < sample_size; j++ {
			dummy_data = append(dummy_data, Dummy{
				XX: 1,
				XY: 2,
				YX: "test",
				YY: "test",
				ZZ: true,
				ZX: 3,
				XZ: 4,
			})
		}
		start := time.Now()
		collection.InsertArray(dummy_data)
		runtimes[sample_size] = time.Since(start).Seconds()
	}

	result_file, err := os.Create(output_path)
	if err != nil {
		log.Fatal("error creating insert.csv")
	}

	writer := csv.NewWriter(result_file)
	writer.Write([]string{"sample_size", "duration"})

	for k, v := range runtimes {
		writer.Write([]string{fmt.Sprintf("%d", k), fmt.Sprintf("%f", v)})
	}
	writer.Flush()
	result_file.Close()
}

func main() {
	analyzeInsert("analysis/insert.csv")
}
