package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

func CPFCNPJValidator(input string) bool {
	input = strings.Trim(input, " ")

	switch len(input) {
	case 11:
		return validateCPF(input)
	case 14:
		return validateCNPJ(input)
	default:
		return false
	}
}

func validateCPF(cpf string) bool {
	firstSequency := []int{
		10, 9, 8, 7, 6, 5, 4, 3, 2,
	}

	checkDigit := getCheckDigit(cpf, firstSequency)
	if string(cpf[9]) == checkDigit {
		secondSequency := []int{
			11, 10, 9, 8, 7, 6, 5, 4, 3, 2,
		}

		checkDigit = getCheckDigit(cpf, secondSequency)
		if string(cpf[10]) == checkDigit {
			return true
		}
	}

	return false
}

func validateCNPJ(cnpj string) bool {
	firstSequency := []int{
		5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2,
	}

	checkDigit := getCheckDigit(cnpj, firstSequency)
	if string(cnpj[12]) == checkDigit {
		secondSequency := []int{
			6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2,
		}

		checkDigit = getCheckDigit(cnpj, secondSequency)
		if string(cnpj[13]) == checkDigit {
			return true
		}
	}

	return false
}

func getCheckDigit(input string, sequencyToSum []int) string {
	sum := 0

	for index, value := range sequencyToSum {
		valueAsInt, err := strconv.Atoi(string(rune(input[index])))
		if err != nil {
			return ""
		}

		sum += (value * valueAsInt)
	}

	mod := sum % 11
	var checkDigit string

	if mod < 2 {
		checkDigit = "0"
	} else {
		checkDigit = fmt.Sprint(11 - mod)
	}

	return checkDigit
}
