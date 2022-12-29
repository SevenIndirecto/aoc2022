package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func deSNAFUfy(snafu string) int {
	m := map[rune]int{
		'=': -2,
		'-': -1,
		'0': 0,
		'1': 1,
		'2': 2,
	}

	sum := 0
	for i, r := range snafu {
		sum += intPow(5, len(snafu) - i - 1) * m[r]
	}
	return sum
}

func SNAFUfy(n int) string {
	s := ""
	firstDigit, firstValue, numDigits := getFirstDigit(n)
	s += firstDigit
	next := n - firstValue

	for p := numDigits-1; p > 0; p-- {
		secondDigit, value := getDigitAtPower(next, p-1)
		s += secondDigit
		next = next - value
	}

	return s
}

func getFirstDigit(n int) (firstDigit string, value int, numDigits int) {
	digits := len(strconv.Itoa(n))

	for i := digits ;; i++ {
		pv := intPow(5, i-1)
		limit := pv / 2

		if 2 * pv - limit <= n && n <= 2 * pv + limit {
			return "2", 2 * pv, i
		}

		if pv - limit <= n && n <= pv + limit {
			return "1", pv, i
		}
	}
}

func getDigitAtPower(n int, power int) (string, int) {
	pv := intPow(5, power)
	limit := pv / 2

	if 2*pv-limit <= n && n <= 2*pv+limit {
		return "2", 2 * pv
	}
	if pv-limit <= n && n <= pv+limit {
		return "1", pv
	}
	if (-1) * limit <= n && n <= limit {
		return "0", 0
	}
	if (-1)*pv - limit <= n && n <= (-1)*pv + limit {
		return "-", -1 * pv
	}
	return "=", -2 * pv
}


func PartOne(lines []string) string {
	sum := 0
	for _, l := range lines {
		num := deSNAFUfy(l)
		sum += num
	}

	return SNAFUfy(sum)
}

func PartTwo(lines []string) int {
	return 0
}

func intPow(n, m int) int {
	if m == 0 {
		return 1
	}
	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}

func LoadLines(path string) ([]string, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	txt := string(dat)
	lines := strings.Split(txt, "\n")
	return lines[:len(lines)-1], nil
}

func main() {
	lines, _ := LoadLines("input.txt")
	fmt.Printf("Part one %v\n", PartOne(lines))
	fmt.Printf("Part two %v\n", PartTwo(lines))
}
