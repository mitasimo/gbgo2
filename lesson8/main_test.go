package main

import "testing"

func TestIterate(t *testing.T) {

	s := struct{}{}

	ref := map[string]struct{}{
		"./tests/f1":     s,
		"./tests/f2":     s,
		"./tests/f3":     s,
		"./tests/sum/s1": s,
		"./tests/sum/s2": s,
		"./tests/sum/s3": s,
	}

	real := make(map[string]struct{})

	filePathChan := make(chan string)
	go func() {
		defer close(filePathChan)
		IterateEntitiesInDirectory("./tests", filePathChan)
	}()
	for path := range filePathChan {
		real[path] = struct{}{}
	}

	// сравнить мапы
	if len(ref) != len(real) {
		t.Fatalf("Количество путей к файлам должны быть %d, но равно %d", len(ref), len(real))
	}

}
