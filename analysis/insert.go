package main

import (
	"encoding/csv"
	"fmt"
	"github.com/snirlugassy/lightdb"
	"log"
	"os"
	"path"
	"reflect"
	"time"
)

func analyzeInsert(output_path string) {
	result_file, err := os.Create(output_path)
	if err != nil {
		log.Fatal("error creating at " + output_path)
	}

	writer := csv.NewWriter(result_file)
	writer.Write([]string{"Sample Size", "Duration (Sec)"})

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
		// START TIMER
		start := time.Now()

		// WORK
		collection.InsertArray(dummy_data)

		// STOP TIMER
		duration := time.Since(start).Seconds()
		writer.Write([]string{fmt.Sprintf("%d", sample_size), fmt.Sprintf("%f", duration)})
	}

	writer.Flush()
	result_file.Close()
}
