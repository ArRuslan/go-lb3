package main

import (
	"bufio"
	"fmt"
	"go-lb3/task2"
	"os"
	"strings"
)

/* Написати пакет функцій для роботи з файлами, які містять текстові дані. Врахувати ситуацію помилок та панік */
/* 1. Створити файл та записати в нього структуровані дані */
/* 2. Вивести створений файл на екран */
/* 3. Видалити з файлу дані відповідно до варіанта */
/* 4. Додати в файл дані відповідно до варіанта */
/* 5. Вивести змінений файл на екран */
/* Структура даних: Структура "Абітурієнт":
- прізвище ім'я по батькові;
- рік народження;
- оцінки вступних іспитів (3);
- Середній бал атестата */
/* Видалення - Видалити елемент із зазначеним номером */
/* Додавання - Додати елемент із зазначеним номером */
func main() {
	/*applicants := []task2.Applicant{
		task2.Applicant{
			FirstName:   "firstName 0",
			LastName:    "lastName 0",
			MiddleName:  "middleName 0",
			YearOfBirth: 2000,
			ExamScores:  [...]int8{1, 2, 3},
			AvgGrade:    5,
		},
		task2.Applicant{
			FirstName:   "firstName 1",
			LastName:    "lastName 1",
			MiddleName:  "middleName 1",
			YearOfBirth: 2001,
			ExamScores:  [...]int8{4, 5, 6},
			AvgGrade:    4,
		},
		task2.Applicant{
			FirstName:   "firstName 2",
			LastName:    "lastName 2",
			MiddleName:  "middleName 2",
			YearOfBirth: 2002,
			ExamScores:  [...]int8{7, 8, 9},
			AvgGrade:    5,
		},
	}*/

	var err error
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("File name: ")
	if !scanner.Scan() {
		fmt.Printf("Failed to read stdin\n")
		return
	}
	filename := scanner.Text()

	var applicantsCount int
	fmt.Print("Applicants count: ")
	_, err = fmt.Scan(&applicantsCount)
	if err != nil {
		fmt.Printf("Failed to read applicants count: %s\n", err)
		return
	}

	applicants := make([]task2.Applicant, applicantsCount)

	for i := 0; i < applicantsCount; i++ {
		applicant := task2.Applicant{}

		fmt.Printf("[%d] First name: ", i)
		if !scanner.Scan() {
			fmt.Printf("Failed to read stdin\n")
			return
		}
		applicant.FirstName = scanner.Text()

		fmt.Printf("[%d] Last name: ", i)
		if !scanner.Scan() {
			fmt.Printf("Failed to read stdin\n")
			return
		}
		applicant.LastName = scanner.Text()

		fmt.Printf("[%d] Middle name: ", i)
		if !scanner.Scan() {
			fmt.Printf("Failed to read stdin\n")
			return
		}
		applicant.MiddleName = scanner.Text()

		fmt.Printf("[%d] Year of birth: ", i)
		_, err = fmt.Scanf("%d", &applicant.YearOfBirth)
		if err != nil {
			fmt.Printf("Failed to read yead of birth: %s\n", err)
			return
		}

		fmt.Printf("[%d] Exam scores: ", i)
		_, err = fmt.Scanf("%d %d %d", &applicant.ExamScores[0], &applicant.ExamScores[1], &applicant.ExamScores[2])
		if err != nil {
			fmt.Printf("Failed to read exam scores: %s\n", err)
			return
		}

		fmt.Printf("[%d] Average grade: ", i)
		_, err = fmt.Scanf("%d", &applicant.AvgGrade)
		if err != nil {
			fmt.Printf("Failed to read yead of birth: %s\n", err)
			return
		}

		applicants[i] = applicant
	}

	err = task2.WriteApplicantsToFile(filename, applicants)
	if err != nil {
		fmt.Printf("Failed to write applicants to a file: %s\n", err)
		return
	}

	contents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read applicants file: %s\n", err)
		return
	}

	sep := strings.Repeat("=", 32)
	fmt.Printf("%s\nRead file\n%s\n%s\n%s\n", sep, sep, string(contents), sep)

	var applicantRemove int
	fmt.Print("Index of applicant to remove: ")
	_, err = fmt.Scan(&applicantRemove)
	if err != nil {
		fmt.Printf("Failed to read applicant index: %s\n", err)
		return
	}

	err = task2.RemoveApplicantFromFile(filename, applicantRemove)
	if err != nil {
		fmt.Printf("Failed to remove applicant: %s\n", err)
		return
	}

	var applicantAdd int
	fmt.Print("Index of applicant to add: ")
	_, err = fmt.Scan(&applicantAdd)
	if err != nil {
		fmt.Printf("Failed to read applicant index: %s\n", err)
		return
	}

	err = task2.AddApplicantToFile(filename, applicantAdd)
	if err != nil {
		fmt.Printf("Failed to add applicant: %s\n", err)
		return
	}

	contents, err = os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read applicants file: %s\n", err)
		return
	}

	fmt.Printf("%s\nRead file\n%s\n%s\n%s\n", sep, sep, string(contents), sep)
}

/*
test data

firstName 0
lastName 0
middleName 0
2000
1 2 3
3
firstName 1
lastName 1
middleName 1
2001
4 5 6
4
firstName 2
lastName 2
middleName 2
2002
7 8 9
5
*/
