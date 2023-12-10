package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

var inputData struct {
	ToSort [][]int `json:"to_sort"`
}

type Response struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       int64   `json:"time_ns"`
}

func main() {
	http.HandleFunc("/process-single", singleSort)
	http.HandleFunc("/process-concurrent", concurrentSort)

	port := ":8080"
	fmt.Printf("Server listening on %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}

func singleSort(w http.ResponseWriter, r *http.Request) {

	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	sortedArrays, timeTaken := processSequential(inputData.ToSort)
	response := Response{
		SortedArrays: sortedArrays,
		TimeNS:       timeTaken.Nanoseconds(),
	}

	writeJSONResponse(w, response)
}

func processSequential(toSort [][]int) ([][]int, time.Duration) {
	startTime := time.Now()

	for i := range toSort {
		sort.Ints(toSort[i])
	}

	return toSort, time.Since(startTime)
}

func processConcurrent(toSort [][]int) ([][]int, time.Duration) {
	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(toSort))
	taskCh := make(chan []int, len(toSort))
	for i := 0; i < len(toSort); i++ {
		go func(i int) {
			defer wg.Done()
			subArray := <-taskCh
			sort.Ints(subArray)
			toSort[i] = subArray
		}(i)
	}
	for _, subArray := range toSort {
		taskCh <- subArray
	}

	close(taskCh)
	wg.Wait()
	return toSort, time.Since(startTime)
}

func concurrentSort(w http.ResponseWriter, r *http.Request) {

	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	sortedArrays, timeTaken := processConcurrent(inputData.ToSort)

	response := Response{
		SortedArrays: sortedArrays,
		TimeNS:       timeTaken.Nanoseconds(),
	}

	writeJSONResponse(w, response)
}

func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
