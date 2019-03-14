package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

var runInParllel bool

func Quicksort(nums []int, parallel bool) ([]int, time.Duration) {
	started := time.Now()
	ch := make(chan int)
	runInParllel = parallel

	go quicksort(nums, ch)

	sorted := make([]int, len(nums))
	i := 0
	for next := range ch {
		sorted[i] = next
		i++
	}
	return sorted, time.Since(started)
}

func quicksort(nums []int, ch chan int) {

	// Choose first number as pivot
	pivot := nums[0]

	// Prepare secondary slices
	smallerThanPivot := make([]int, 0)
	largerThanPivot := make([]int, 0)

	// Slice except pivot
	nums = nums[1:]

	// Go over slice and sort
	for _, i := range nums {
		switch {
		case i <= pivot:
			smallerThanPivot = append(smallerThanPivot, i)
		case i > pivot:
			largerThanPivot = append(largerThanPivot, i)
		}
	}

	var ch1 chan int
	var ch2 chan int

	// Now do the same for the two slices
	if len(smallerThanPivot) > 1 {
		ch1 = make(chan int, len(smallerThanPivot))
		if runInParllel {
			go quicksort(smallerThanPivot, ch1)
		} else {
			quicksort(smallerThanPivot, ch1)
		}
	}
	if len(largerThanPivot) > 1 {
		ch2 = make(chan int, len(largerThanPivot))
		if runInParllel {
			go quicksort(largerThanPivot, ch2)
		} else {
			quicksort(largerThanPivot, ch2)
		}
	}

	// Wait until the sorting finishes for the smaller slice
	if len(smallerThanPivot) > 1 {
		for i := range ch1 {
			ch <- i
		}
	} else if len(smallerThanPivot) == 1 {
		ch <- smallerThanPivot[0]
	}
	ch <- pivot

	if len(largerThanPivot) > 1 {
		for i := range ch2 {
			ch <- i
		}
	} else if len(largerThanPivot) == 1 {
		ch <- largerThanPivot[0]
	}

	close(ch)
}

func randomArray(len int) []int {
	a := make([]int, len)
	for i := 0; i <= len-1; i++ {
		a[i] = rand.Intn(len)
	}
	return a
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	totalNumbers := 1000
	var data = [][]string{{"maxRandom", "totalNumbers", "SequentialTime", "ParallelTime"}}

	for totalNumbers < 40000000 {

		a := randomArray(totalNumbers)
		seq_a := make([]int, len(a))
		copy(seq_a, a)
		_, secs := Quicksort(a, false)

		fmt.Println("Sequential Quicksort took : ", secs)

		_, par_secs := Quicksort(seq_a, true)

		fmt.Println("Parallel Quicksort took : ", par_secs)

		x := []string{strconv.Itoa(totalNumbers), strconv.Itoa(totalNumbers), strconv.Itoa(int(secs.Nanoseconds())), strconv.Itoa(int(par_secs.Nanoseconds()))}
		data = append(data, x)
		totalNumbers += 150000

	}

	file, err := os.Create("results_golang.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}
}
