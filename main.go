package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
)

type dataStruct struct {
	Min  int   `json:"min"`
	Max  int   `json:"max"`
	Data []int `json:"data"`
}

func main() {
	var bins, choice int
	fmt.Println("BINNING DATA")
	fmt.Println("Make sure your data.json file is within the same directory")
	fmt.Println("Please enter your choice:")
	fmt.Println("1. Equal Width Binning")
	fmt.Println("2. Equal Frequency Binning")
	fmt.Println("3. Equal Frequency Binning (Smoothing by Means)")
	fmt.Println("4. Equal Frequency Binning (Smoothing by Bounderies)")
	fmt.Println("5. Equal Frequency Binning (Smoothing by Median)")
	fmt.Println("Enter choice: ")
	fmt.Scan(&choice)
	fmt.Println("Enter bins: ")
	fmt.Scan(&bins)

	var finalData []dataStruct
	if choice == 1 {
		finalData = equalWidth(bins)
	}
	if choice == 2 {
		finalData = equalFrequency(bins)
	}
	if choice == 3 {
		finalData = equalFrequency(bins)
		smoothingMeans(finalData)
	}
	if choice == 4 {
		finalData = equalFrequency(bins)
		smoothingBounderies(finalData)
	}
	if choice == 5 {
		finalData = equalFrequency(bins)
		smoothingMedian(finalData)
	}

	// Writing the file
	file, err := os.Create("final_data.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(finalData); err != nil {
		fmt.Println("Error encoding data to JSON:", err)
		return
	}
}

func equalWidth(bins int) []dataStruct {
	data := readAndSort()
	max := data[len(data)-1]
	min := data[0]

	w := ((max - min) / bins)

	finalData := make([]dataStruct, bins)

	wCount := 0
	for i := 0; i < bins; i++ {
		finalData[i].Min = wCount + 1
		if i == 0 {
			finalData[i].Min = wCount
		}
		finalData[i].Max = wCount + w
		wCount = wCount + w
	}
	for i, itemFinal := range finalData {
		for _, item := range data {
			if itemFinal.Min > item {
				continue
			}
			if itemFinal.Max >= item {
				finalData[i].Data = append(finalData[i].Data, item)
			}
			if itemFinal.Max < item {
				break
			}
		}
	}
	return finalData
}
func equalFrequency(bins int) []dataStruct {
	data := readAndSort()
	lengthFlt := float32(len(data)) / float32(bins)
	length := int(math.Floor(float64(lengthFlt)))

	i := 0
	finalData := make([]dataStruct, bins)

	for _, item := range data {
		if len(finalData[i].Data) == length && i+1 < bins {
			finalData[i].Max = finalData[i].Data[len(finalData[i].Data)-1]
			i += 1
			finalData[i].Min = finalData[i-1].Data[len(finalData[i-1].Data)-1] + 1
		}
		if len(finalData[i].Data) >= length && i+2 > bins {
			finalData[i].Data = append(finalData[i].Data, item)
			finalData[i].Max = finalData[i].Data[len(finalData[i].Data)-1]
			continue
		}
		if len(finalData[i].Data) <= length {
			finalData[i].Data = append(finalData[i].Data, item)
			finalData[i].Max = finalData[i].Data[len(finalData[i].Data)-1]
			continue
		}
	}

	return finalData
}
func smoothingMeans(data []dataStruct) {
	for i, slice := range data {
		total := 0
		means := 0
		for _, item := range slice.Data {
			total += item
		}
		means = total / len(slice.Data)
		for k := range slice.Data {
			data[i].Data[k] = means
		}
	}
}
func smoothingMedian(data []dataStruct) {
	for i, slice := range data {
		for k := range slice.Data {
			index := int(math.Floor(float64(len(slice.Data) / 2)))
			data[i].Data[k] = slice.Data[int(index)]
		}
	}
}
func smoothingBounderies(data []dataStruct) {
	for i, slice := range data {
		for k := range slice.Data {
			index := int(math.Floor(float64(len(slice.Data) / 2)))
			if k <= index {
				data[i].Data[k] = slice.Data[0]
			}
			if k > index {
				data[i].Data[k] = slice.Data[len(slice.Data)-1]
			}
		}
	}
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
