package task1

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Applicant struct {
	FirstName   string
	LastName    string
	MiddleName  string
	YearOfBirth int16
	ExamScores  [3]int8
	AvgGrade    int8
}

func writeLine(writer *bufio.Writer, data string) error {
	var err error
	_, err = writer.Write([]byte(data))
	if err != nil {
		return err
	}
	_, err = writer.WriteRune('\n')
	if err != nil {
		return err
	}

	return nil
}

func WriteApplicantsToFile(filename string, applicants []Applicant) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, applicant := range applicants {
		if err = writeLine(writer, applicant.FirstName); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		if err = writeLine(writer, applicant.LastName); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		if err = writeLine(writer, applicant.MiddleName); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		yearOfBirth := fmt.Sprintf("%d", applicant.YearOfBirth)
		if err = writeLine(writer, yearOfBirth); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		examScores := fmt.Sprintf("%d %d %d", applicant.ExamScores[0], applicant.ExamScores[1], applicant.ExamScores[2])
		if err = writeLine(writer, examScores); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		avgGrade := fmt.Sprintf("%d", applicant.AvgGrade)
		if err = writeLine(writer, avgGrade); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		if err = writeLine(writer, ""); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return writer.Flush()
}

func ReadApplicantsFromFile(filename string) ([]Applicant, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var applicants []Applicant

	for scanner.Scan() {
		var applicant Applicant

		applicant.FirstName = scanner.Text()
		if applicant.FirstName == "" {
			break
		}

		if !scanner.Scan() {
			return nil, fmt.Errorf("failed to read last name from a file: %w", scanner.Err())
		}
		applicant.LastName = scanner.Text()

		if !scanner.Scan() {
			return nil, fmt.Errorf("failed to read middle name from a file: %w", scanner.Err())
		}
		applicant.MiddleName = scanner.Text()

		if !scanner.Scan() {
			return nil, fmt.Errorf("failed to read year of birth from a file: %w", scanner.Err())
		}
		_, err = fmt.Sscanf(scanner.Text(), "%d", &applicant.YearOfBirth)
		if err != nil {
			return nil, fmt.Errorf("failed to parse year of birth from a file: %w", err)
		}

		if !scanner.Scan() {
			return nil, fmt.Errorf("failed to read exam scores from a file: %w", scanner.Err())
		}
		_, err = fmt.Sscanf(scanner.Text(), "%d %d %d", &applicant.ExamScores[0], &applicant.ExamScores[1], &applicant.ExamScores[2])
		if err != nil {
			return nil, fmt.Errorf("failed to parse exam scores from a file: %w", err)
		}

		if !scanner.Scan() {
			return nil, fmt.Errorf("failed to read average grade from a file: %w", scanner.Err())
		}
		_, err = fmt.Sscanf(scanner.Text(), "%d", &applicant.AvgGrade)
		if err != nil {
			return nil, fmt.Errorf("failed to parse average grade from a file: %w", err)
		}

		if !scanner.Scan() {
			return nil, fmt.Errorf("failed to read empty line from a file: %w", scanner.Err())
		}
		if scanner.Text() != "" {
			return nil, fmt.Errorf("invalid file data: expected empty line, got \"%s\"", scanner.Text())
		}

		applicants = append(applicants, applicant)
	}

	return applicants, nil
}

func readExactlyOneLine(r io.Reader) (string, error) {
	buffer := make([]byte, 1)
	result := ""
	readBytes, err := r.Read(buffer)
	for readBytes > 0 && err == nil && buffer[0] != '\n' {
		result += string(buffer[0])
		readBytes, err = r.Read(buffer)
	}

	if err != nil && err != io.EOF {
		return "", err
	}

	return result, nil
}

func RemoveApplicantFromFile(filename string, applicantIndex uint) error {
	file, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var startPos, endPos int64

	startPos, err = file.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("failed to get current position: %w", err)
	}

	line, err := readExactlyOneLine(file)

	for line != "" {
		line, err = readExactlyOneLine(file)
		if err != nil {
			return fmt.Errorf("failed to read file (last name): %w", err)
		}
		if line == "" {
			return fmt.Errorf("invalid file structure")
		}

		line, err = readExactlyOneLine(file)
		if err != nil {
			return fmt.Errorf("failed to read file (middle name): %w", err)
		}
		if line == "" {
			return fmt.Errorf("invalid file structure")
		}

		line, err = readExactlyOneLine(file)
		if err != nil {
			return fmt.Errorf("failed to read file (year of birth): %w", err)
		}
		if line == "" {
			return fmt.Errorf("invalid file structure")
		}

		line, err = readExactlyOneLine(file)
		if err != nil {
			return fmt.Errorf("failed to read file (exam scores): %w", err)
		}
		if line == "" {
			return fmt.Errorf("invalid file structure")
		}

		line, err = readExactlyOneLine(file)
		if err != nil {
			return fmt.Errorf("failed to read file (average grade): %w", err)
		}
		if line == "" {
			return fmt.Errorf("invalid file structure")
		}

		line, err = readExactlyOneLine(file)
		if err != nil {
			return fmt.Errorf("failed to read file (empty line): %w", err)
		}
		if line != "" {
			return fmt.Errorf("invalid file structure: expected empty line")
		}

		endPos, err = file.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("failed to get current position: %w", err)
		}

		if applicantIndex == 0 {
			break
		}

		applicantIndex--
		startPos = endPos

		line, err = readExactlyOneLine(file)
		if err != nil {
			return fmt.Errorf("failed to read file (first name): %w", err)
		}
		if line == "" {
			return fmt.Errorf("invalid file structure")
		}
	}

	structLen := endPos - startPos
	buffer := make([]byte, structLen)

	writePos := startPos
	readPos := endPos

	_, err = file.Seek(readPos, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to set read position: %w", err)
	}
	readBytes, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file: %w", err)
	}
	for readBytes > 0 {
		_, err = file.Seek(writePos, io.SeekStart)
		if err != nil {
			return fmt.Errorf("failed to set write position: %w", err)
		}

		_, err = file.Write(buffer[:readBytes])
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		readPos += int64(readBytes)
		writePos += int64(readBytes)

		_, err = file.Seek(readPos, io.SeekStart)
		if err != nil {
			return fmt.Errorf("failed to set read position: %w", err)
		}

		readBytes, err = file.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read file: %w", err)
		}
	}

	err = file.Truncate(writePos)
	if err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}

	return nil
}

func AddApplicantToFile(filename string, applicantIndex int) error {
	//

	return nil
}
