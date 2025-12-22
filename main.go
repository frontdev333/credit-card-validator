package main

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Bank struct {
	Name    string
	BinFrom int
	BinTo   int
}

func loadBankData(path string) ([]Bank, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var banks []Bank

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ",")
		name := data[0]
		from, err := strconv.Atoi(data[1])
		if err != nil {
			return nil, err
		}

		to, err := strconv.Atoi(data[2])
		if err != nil {
			return nil, err
		}

		banks = append(banks, Bank{
			Name:    name,
			BinFrom: from,
			BinTo:   to,
		})
	}

	return banks, nil
}

func extractBIN(cardNumber string) (int, error) {
	bin, err := strconv.Atoi(cardNumber[:6])
	if err != nil {
		return 0, err
	}

	return bin, nil
}

func identifyBank(bin int, banks []Bank) (string, error) {
	for _, bank := range banks {
		if bin >= bank.BinFrom && bin <= bank.BinTo {
			return bank.Name, nil
		}
	}
	return "", errors.New("bin is not recognized")
}

func validateLuhn(cardNumber string) bool {
	var sum int
	digits, err := strNumToIntSliceNum(cardNumber)
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	double := false

	for i := len(digits) - 1; i >= 0; i-- {
		digit := digits[i]
		if double {
			digit = digit * 2
			if digit > 9 {
				digit = digit/10 + digit%10
			}
		}

		double = !double
		sum = digit + sum
	}

	return sum%10 == 0
}

func strNumToIntSliceNum(num string) ([]int, error) {
	num = strings.TrimSpace(num)
	digits := make([]int, len(num))

	for i, v := range num {
		if v > '9' || v < '0' {
			return nil, errors.New(fmt.Sprintf("card number contains non-numeric character: %#v", v))
		}
		digits[i] = int(v - '0')
	}

	return digits, nil

}

func getUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if text == "\n" {
		return text, nil
	}
	if err != nil {
		return "", err
	}

	return text, nil
}

func validateInput(card string) bool {
	card = strings.TrimSpace(card)
	cardlen := len(card)

	if cardlen < 13 || cardlen > 19 {
		return false
	}

	for _, ch := range card {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	return true

}

func main() {
	fmt.Println("Добро пожаловать в программу валидации карт!")
	banks, err := loadBankData("banks.txt")
	if err != nil {
		panic(err)
	}

	for {
		fmt.Print("Please enter a credit card number (or press Enter to quit): ")
		cardNumber, err := getUserInput()
		if err != nil {
			panic(err)
		}

		if cardNumber == "\n" {
			fmt.Println("Программа завершена\n")
			break
		}
		fmt.Println("You entered:", cardNumber)

		if !validateInput(cardNumber) {
			fmt.Println("Ошибка формата")
			continue
		}

		if !validateLuhn(cardNumber) {
			fmt.Println("Невалидный номер")
			continue
		}

		bin, err := extractBIN(cardNumber)
		if err != nil {
			panic(err)
		}
		fmt.Println("Номер карты валиден")

		res, err := identifyBank(bin, banks)
		if err != nil {
			fmt.Println("Банк не определен")
		}
		fmt.Printf("Банк: %s", res)
	}
}
