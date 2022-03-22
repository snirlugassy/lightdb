package main

import "log"

type Dummy struct {
	XX int
	XY int
	YX string
	YY string
	ZZ bool
	ZX byte
	XZ rune
}

func main() {
	log.Println("Running insert analysis")
	analyzeInsert("analysis/insert.csv")

	log.Println("Running commit & pull analysis")
	analyzeCommitAndPull("analysis/commit_pull.csv")
}
