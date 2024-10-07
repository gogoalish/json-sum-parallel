package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

type Data struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("usage: go run . <filename> <number of goroutiness>")
	}

	filename := os.Args[1]
	numGoroutines, err := strconv.Atoi(os.Args[2])
	if err != nil || numGoroutines <= 0 {
		log.Fatal("usage: go run . <filename> <number of goroutiness>")
	}

	dataChan := make(chan Data, 100)
	resultChan := make(chan int, numGoroutines)
	var wg sync.WaitGroup

	// чтение
	wg.Add(1)
	go reader(filename, dataChan, &wg)

	// вычисление
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go worker(dataChan, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// результат
	totalSum := 0
	for sum := range resultChan {
		totalSum += sum
	}

	fmt.Printf("Total sum: %d\n", totalSum)
}

func reader(filename string, dataChan chan<- Data, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(1, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if _, err := decoder.Token(); err != nil {
		log.Fatal(2, err)
	}

	for decoder.More() {
		var item Data
		if err := decoder.Decode(&item); err != nil {
			log.Fatal(3, err)
		}
		dataChan <- item
	}

	if _, err := decoder.Token(); err != nil {
		log.Fatal(4, err)
	}

	close(dataChan)
}

func worker(dataChan <-chan Data, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	localSum := 0
	for item := range dataChan {
		localSum += item.A + item.B
	}
	resultChan <- localSum
}
