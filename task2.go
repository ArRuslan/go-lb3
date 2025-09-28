package main

import (
	"fmt"
	task1 "go-lb3/task2"
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
	applicants := []task1.Applicant{
		task1.Applicant{
			FirstName:   "firstName 0",
			LastName:    "lastName 0",
			MiddleName:  "middleName 0",
			YearOfBirth: 2000,
			ExamScores:  [...]int8{1, 2, 3},
			AvgGrade:    5,
		},
		task1.Applicant{
			FirstName:   "firstName 1",
			LastName:    "lastName 1",
			MiddleName:  "middleName 1",
			YearOfBirth: 2001,
			ExamScores:  [...]int8{4, 5, 6},
			AvgGrade:    4,
		},
		task1.Applicant{
			FirstName:   "firstName 2",
			LastName:    "lastName 2",
			MiddleName:  "middleName 2",
			YearOfBirth: 2002,
			ExamScores:  [...]int8{7, 8, 9},
			AvgGrade:    5,
		},
	}

	err := task1.WriteApplicantsToFile("applicants.txt", applicants)
	if err != nil {
		panic(err)
	}

	applicantsRead, err := task1.ReadApplicantsFromFile("applicants.txt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", applicantsRead)

	err = task1.RemoveApplicantFromFile("applicants.txt", 2)
	if err != nil {
		panic(err)
	}
}
