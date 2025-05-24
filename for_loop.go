package basics

import "fmt"

func main() {

	//simple iteration over a range
	for i := 5; i < 10; i++ {
		fmt.Println(i)
	}

	//iterate over collection
	number := []int{1, 2, 3, 4, 5, 6}
	for index, value := range number {
		fmt.Printf("index: %d, value:%d\n", index, value)
	}

	//New Problem
	for i := 1; i < 100; i++ {
		if i%2 == 0 {
			continue

		}
		fmt.Println("Odd Number:-", i)

		if i == 17 {
			break
		}
	}

}
