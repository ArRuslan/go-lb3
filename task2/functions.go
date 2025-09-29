package task2

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

func readExactlyOneLine(r io.Reader) (string, bool, error) {
	buffer := make([]byte, 1)
	result := ""
	readBytes, err := r.Read(buffer)
	for readBytes > 0 && err == nil && buffer[0] != '\n' {
		result += string(buffer[0])
		readBytes, err = r.Read(buffer)
	}

	if err != nil {
		if err == io.EOF {
			return result, true, nil
		}
		return "", false, err
	}

	return result, false, nil
}

func skipLine(r io.Reader, lineDesc string, empty bool, allowNonEmpty bool, failOnEof bool) (error, bool) {
	line, eof, err := readExactlyOneLine(r)
	if eof && failOnEof {
		return fmt.Errorf("unexpected eof while reading %s", lineDesc), false
	}
	if err != nil {
		return fmt.Errorf("failed to read file (%s): %w", lineDesc, err), false
	}
	if (line != "" && empty) || (line == "" && !empty && !allowNonEmpty) {
		return fmt.Errorf("invalid file structure (when reading %s)", lineDesc), false
	}

	return nil, eof
}

func getApplicantStartEndPos(file *os.File, applicantIndex int) (int64, int64, error) {
	if applicantIndex < 0 {
		return 0, 0, fmt.Errorf("applicant index must be greater or equal to 0, got %d", applicantIndex)
	}

	var startPos, endPos int64

	for applicantIndex >= 0 {
		startPos = endPos

		if err, eof := skipLine(file, "first name", false, true, false); err != nil || eof {
			if eof {
				return 0, 0, fmt.Errorf("applicant with this number does not exist")
			}
			return 0, 0, err
		}

		if err, _ := skipLine(file, "last name", false, false, true); err != nil {
			return 0, 0, err
		}

		if err, _ := skipLine(file, "middle name", false, false, true); err != nil {
			return 0, 0, err
		}

		if err, _ := skipLine(file, "year of birth", false, false, true); err != nil {
			return 0, 0, err
		}

		if err, _ := skipLine(file, "exam scores", false, false, true); err != nil {
			return 0, 0, err
		}

		if err, _ := skipLine(file, "average grade", false, false, true); err != nil {
			return 0, 0, err
		}

		if err, _ := skipLine(file, "empty line", true, false, true); err != nil {
			return 0, 0, err
		}

		var err error
		endPos, err = file.Seek(0, io.SeekCurrent)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to get current position: %w", err)
		}

		applicantIndex--
	}

	return startPos, endPos, nil
}

func RemoveApplicantFromFile(filename string, applicantIndex int) error {
	file, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	startPos, endPos, err := getApplicantStartEndPos(file, applicantIndex)

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
	file, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	startPos, endPos, err := getApplicantStartEndPos(file, applicantIndex)

	structLen := endPos - startPos
	buffer := make([]byte, structLen)

	_, err = file.Seek(startPos, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to set read position: %w", err)
	}
	readBytes, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file: %w", err)
	}
	if int64(readBytes) != structLen {
		return fmt.Errorf("expected to read %d bytes, actualle read %d", structLen, readBytes)
	}

	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to set write position: %w", err)
	}

	_, err = file.Write(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
