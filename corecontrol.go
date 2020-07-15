package main

import (
	"fmt"
	"os"

	"io/ioutil"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cores = kingpin.Flag("cores", "Total number of cores that should be online").Required().Short('c').Int()
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func write(content string, filename string) {
	zero := []byte(content)
	err := ioutil.WriteFile(filename, zero, 0644)
	checkError(err)
}

func isOnline(coreIndex int) bool {

	file, err := os.Open(fmt.Sprintf("/sys/devices/system/cpu/cpu%d/online", coreIndex))
	checkError(err)
	defer file.Close()

	var isOnline int
	_, err = fmt.Fscanf(file, "%d", &isOnline)
	checkError(err)
	return isOnline == 1
}

func main() {
	kingpin.Parse()

	file, err := os.Open("/sys/devices/system/cpu/present")
	checkError(err)
	defer file.Close()

	var lowestCoreIndex int
	var highestCoreIndex int
	_, err = fmt.Fscanf(file, "%d-%d", &lowestCoreIndex, &highestCoreIndex)
	checkError(err)
	totalCores := highestCoreIndex - lowestCoreIndex + 1

	fmt.Println("Available core count = ", totalCores)

	for i := 1; i < totalCores; i++ {
		filename := fmt.Sprintf("/sys/devices/system/cpu/cpu%d/online", i)
		online := isOnline(i)
		if i < *cores {
			if !online {
				fmt.Println("Turning on core of index:", i)
				write("1", filename)
			}
		}
		if i >= *cores {
			if online {
				fmt.Println("Turning off core of index:", i)
				write("0", filename)
			}
		}
	}
}
