package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Grade string

const (
	A Grade = "A"
	B Grade = "B"
	C Grade = "C"
	F Grade = "F"
)

type student struct {
	firstName, lastName, university                string
	test1Score, test2Score, test3Score, test4Score int
}

type studentStat struct {
	student
	finalScore float32
	grade      Grade
}

func parseCSV(filePath string) ([]student, error) {
	var students []student
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	strData := string(data)
	//fmt.Println(strData)
	strData = strings.ReplaceAll(strData, "\r\n", "\n")
	records := strings.Split(strings.TrimSpace(strData), "\n")
	records = records[1:]
	//fmt.Println(records)

	for i := 0; i < len(records); i++ {
		fields := strings.Split(records[i], ",")
		test1Score, errT1 := strconv.Atoi(fields[3])
		test2Score, errT2 := strconv.Atoi(fields[4])
		test3Score, errT3 := strconv.Atoi(fields[5])
		test4Score, errT4 := strconv.Atoi(fields[6])

		if errT1 != nil || errT2 != nil || errT3 != nil || errT4 != nil {
			return nil, errors.New("error parsing string")
		}

		students = append(students, student{
			firstName:  fields[0],
			lastName:   fields[1],
			university: fields[2],
			test1Score: test1Score,
			test2Score: test2Score,
			test3Score: test3Score,
			test4Score: test4Score,
		})
	}
	return students, nil
}

func calculateGrade(students []student) []studentStat {
	var stStat []studentStat
	for _, st := range students {
		avg := float32(st.test1Score+st.test2Score+st.test3Score+st.test4Score) / 4
		grade := findGrade(float64(avg))
		stStat = append(stStat, studentStat{
			student:    st,
			finalScore: float32(avg),
			grade:      grade,
		})
	}
	return stStat
}

func findOverallTopper(gradedStudents []studentStat) studentStat {
	sort.Slice(gradedStudents, func(i, j int) bool {
		return gradedStudents[i].finalScore > gradedStudents[j].finalScore
	})
	return gradedStudents[0]
}

func findTopperPerUniversity(gs []studentStat) map[string]studentStat {
	mp := make(map[string]studentStat)

	for _, stat := range gs {
		existing, ok := mp[stat.university]
		if !ok || stat.finalScore > existing.finalScore {
			mp[stat.university] = stat
		}
	}
	return mp
}

func main() {
	students, err := parseCSV("C:\\Users\\ASUS\\GolangProjects\\ankur-tdd-go\\student-grading-go\\grades.csv")
	if err != nil {
		log.Fatal("Error opening file")
	}
	stStat := calculateGrade(students)
	fmt.Println(stStat)
	fmt.Println(findOverallTopper(stStat))
	fmt.Println(findTopperPerUniversity(stStat))
}

func findGrade(grade float64) Grade {
	if grade < 35 {
		return "F"
	} else if grade < 50 {
		return "C"
	} else if grade < 70 {
		return "B"
	} else {
		return "A"
	}
}
