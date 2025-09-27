package main

import (
	"bufio"
	"fmt"
	"go-lb3/task1"
	"math/rand"
	"os"
)

/* Написати пакет функцій для роботи з файлами, що містять числові дані. Врахувати ситуацію помилок та панік */
/* 1. Сформувати файл із модулів цілих чисел. Знайти Кількість парних чисел серед компонентів файлу */
/* 2. Сформувати файл із квадратного коріння цілих чисел. Знайти Суму компонент файлу */
/* 3. Сформувати файл із чисел послідовності. Знайти Суму компонент файлу */
func main() {
	var err error
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("File name: ")
	if !scanner.Scan() {
		fmt.Printf("Failed to read stdin\n")
		return
	}
	filename := scanner.Text()

	var numbersCount int
	fmt.Print("Numbers count for sub-task 1: ")
	_, err = fmt.Scan(&numbersCount)
	if err != nil {
		fmt.Printf("Failed to read number count: %s\n", err)
		return
	}

	numbers := make(task1.Numbers, numbersCount)
	for i := 0; i < numbersCount; i++ {
		numbers[i] = int64(rand.Intn(128))
	}

	fmt.Printf("Numbers: %v\n", numbers)

	err = numbers.WriteToFile(filename)
	if err != nil {
		fmt.Printf("Failed to write numbers to a file: %s\n", err)
		return
	}

	evenCount, err := task1.CountEvenFromFile(filename)
	if err != nil {
		fmt.Printf("Failed to count even numbers in a file: %s\n", err)
		return
	}
	fmt.Printf("Even numbers count: %d\n", evenCount)

	fmt.Println("--------------------------")

	err = task1.CreateSqrtFile(filename, numbers)
	if err != nil {
		fmt.Printf("Failed to write sqrt numbers to a file: %s\n", err)
		return
	}

	numbers, _ = task1.ReadNumbersFromFile(filename)
	fmt.Printf("Numbers: %v\n", numbers)

	numsSum, err := task1.SumFromFile(filename)
	if err != nil {
		fmt.Printf("Failed to find sum of numbers in a file: %s\n", err)
		return
	}
	fmt.Printf("Sum of all numbers in a file: %d\n", numsSum)

	fmt.Println("--------------------------")

	var seqStart, seqStep, seqCount int
	fmt.Print("Start, step and count for sequence of numbers for sub-task 3: ")
	_, err = fmt.Scanf("%d %d %d", &seqStart, &seqStep, &seqCount)
	if err != nil {
		fmt.Printf("Failed to sequence parameters: %s\n", err)
		return
	}

	err = task1.CreateSequenceFile(filename, seqStart, seqStep, seqCount)
	if err != nil {
		fmt.Printf("Failed to write sequence numbers to a file: %s\n", err)
		return
	}

	numbers, _ = task1.ReadNumbersFromFile(filename)
	fmt.Printf("Numbers: %v\n", numbers)

	numsSum, err = task1.SumFromFile(filename)
	if err != nil {
		fmt.Printf("Failed to find sum of numbers in a file: %s\n", err)
		return
	}
	fmt.Printf("Sum of all numbers in a file: %d\n", numsSum)

	//
}
