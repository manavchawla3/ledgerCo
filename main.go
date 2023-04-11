package main

import (
	"fmt"
	"geektrust/logic"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func ProcessInput(inputString string) {
	inputs := strings.Split(inputString, " ")
	command := inputs[0]
	bank := inputs[1]
	borrower := inputs[2]

	switch true {

	case command == "LOAN":
		amount, ok1 := strconv.ParseFloat(inputs[3], 64)
		tenure, ok2 := strconv.ParseInt(inputs[4], 10, 64)
		rate, ok3 := strconv.ParseFloat(inputs[5], 64)

		if ok1 != nil || ok2 != nil || ok3 != nil {
			fmt.Errorf("%s: %s", inputString, "INVALID PARAMS FOR COMMAND!")
		}

		logic.CreateLoan(bank, borrower, amount, int64(tenure), rate)

	case command == "PAYMENT" || command == "BALANCE":
		loan := logic.GetLoan(bank, borrower)

		if loan == nil {
			fmt.Errorf("%s: %s", inputString, "INVALID PARAMS FOR COMMAND!")
		}

		switch command {
		case "PAYMENT":
			amount, ok4 := strconv.ParseFloat(inputs[3], 64)
			emi, ok5 := strconv.ParseInt(inputs[4], 10, 64)
			if ok4 != nil || ok5 != nil {
				fmt.Errorf("%s: %s", inputString, "INVALID PARAMS FOR COMMAND!")
			}

			loan.ProcessRepayment(amount, int64(emi))

		case "BALANCE":
			emi, ok4 := strconv.ParseInt(inputs[3], 10, 64)
			if ok4 != nil {
				fmt.Errorf("%s: %s", inputString, "INVALID PARAMS FOR COMMAND!")
			}

			output := loan.Balance(int64(emi))
			if output == nil {
				fmt.Errorf("%s: %s", inputString, "INVALID PARAMS FOR COMMAND!")
			}

			fmt.Printf("%v %v %v %v \n", bank, borrower, output.AmountPaid, output.RemEMi)
		}

	default:
		fmt.Errorf("%s: %s", inputString, "INVALID COMMAND!")
	}
}

func main() {
	if len(os.Args) < 2 {
		panic("FILE PATH IS MANDATORY!")
	}

	filePath := os.Args[1]
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
		return
	}

	for _, input := range strings.Split(string(data), "\n") {
		ProcessInput(input)
	}
}
