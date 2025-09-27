package task1

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

type Numbers []int64

func (n Numbers) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	err = binary.Write(writer, binary.LittleEndian, int64(len(n)))
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	for _, num := range n {
		err = binary.Write(writer, binary.LittleEndian, num)
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return writer.Flush()
}

func ReadNumbersFromFile(filename string) ([]int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var numbersCount int64
	err = binary.Read(file, binary.LittleEndian, &numbersCount)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	numbers := make([]int64, numbersCount)

	var i int64
	for ; i < numbersCount; i++ {
		err = binary.Read(file, binary.LittleEndian, &numbers[i])
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
	}

	return numbers, nil
}

func CountEvenFromFile(filename string) (int, error) {
	numbers, err := ReadNumbersFromFile(filename)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, n := range numbers {
		if int(n)%2 == 0 {
			count++
		}
	}

	return count, nil
}

func CreateSqrtFile(filename string, numbers []int64) error {
	nums := make(Numbers, len(numbers))

	for i, n := range numbers {
		if n < 0 {
			return fmt.Errorf("can't calculate sqrt of a negarive number %d (at index %d)", n, i)
		}
		nums[i] = int64(math.Sqrt(float64(n)))
	}

	return nums.WriteToFile(filename)
}

func SumFromFile(filename string) (int64, error) {
	numbers, err := ReadNumbersFromFile(filename)
	if err != nil {
		return 0, err
	}

	var result int64
	for _, n := range numbers {
		result += n
	}

	return result, nil
}

func CreateSequenceFile(filename string, start, step, count int) error {
	if count <= 0 {
		return fmt.Errorf("sequence count must be positive")
	}

	nums := make(Numbers, count)
	for i := 0; i < count; i++ {
		nums[i] = int64(start + i*step)
	}

	return nums.WriteToFile(filename)
}
