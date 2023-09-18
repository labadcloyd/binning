package main

import (
	"encoding/json"
	"log"
	"os"
)

type dataStruct struct {
	min  int
	max  int
	data []int
}

func main() {
	equalWidth(3)
}

func equalWidth(bins int) {
	data := readAndSort()
	max := data[len(data)-1]
	min := data[0]

	w := ((max - min) / bins)

	finalData := make([]dataStruct, bins)

	wCount := 0
	for i := 0; i < bins; i++ {
		finalData[i].min = wCount + 1
		if i == 0 {
			finalData[i].min = wCount
		}
		finalData[i].max = wCount + w
		wCount = wCount + w
	}
	for i, itemFinal := range finalData {
		for _, item := range data {
			if itemFinal.min > item {
				continue
			}
			if itemFinal.max >= item {
				finalData[i].data = append(finalData[i].data, item)
			}
			if itemFinal.max < item {
				break
			}
		}
	}
	log.Println(finalData)
}

// READ AND SORT
func readAndSort() []int {
	file, err := os.Open("data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	// Read open bracket
	_, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}

	data := []int{}
	for dec.More() {
		var m int
		err := dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, m)
	}

	quicksort(data, 0, len(data)-1)

	_, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// SORTING ALGORITHM
func quicksort(arr []int, low int, high int) {
	if low < high {
		pivotIndex := partition(arr, low, high)
		quicksort(arr, low, pivotIndex-1)
		quicksort(arr, pivotIndex+1, high)
	}
}
func partition(arr []int, low int, high int) int {
	pivot := arr[high]
	i := low - 1
	for j := low; j < high; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}
