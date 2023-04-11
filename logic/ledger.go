package logic

import "math"

var Loans = map[string]*Loan{}

type Result struct {
	AmountPaid float64
	RemEMi     int64
}

type Loan struct {
	Id           string
	BankName     string
	BorrowerName string
	Principal    float64
	Tenure       int64
	Rate         float64
	TotalAmount  float64
	NoOfEMI      int64
	EmiAmount    float64
	Repayments   []Repayment
}

type Repayment struct {
	Amount       float64
	PostEMICount int64
}

func (loan *Loan) init() {
	loan.calculateCollectionAmount()
	loan.calculateEmi()
}

func (loan *Loan) calculateCollectionAmount() {
	loan.TotalAmount = loan.Principal + ((loan.Principal * loan.Rate * (float64(loan.Tenure))) / 100)
}

func (loan *Loan) calculateEmi() {
	loan.NoOfEMI = loan.Tenure * 12
	loan.EmiAmount = math.Ceil(loan.TotalAmount / float64(loan.NoOfEMI))
}

func (loan *Loan) ProcessRepayment(amount float64, emi int64) {
	repayment := Repayment{
		Amount:       amount,
		PostEMICount: emi,
	}

	loan.Repayments = append(loan.Repayments, repayment)
}

func (loan *Loan) getLumpSumTotalTillEmi(emi int64) float64 {
	lumpSumTotal := float64(0)

	for _, repayment := range loan.Repayments {
		if repayment.PostEMICount <= emi {
			lumpSumTotal += repayment.Amount
		}
	}

	return lumpSumTotal
}

func (loan *Loan) getRemainingEMIs(AmountPaid float64) int64 {
	remAmount := loan.TotalAmount - AmountPaid
	RemEMis := math.Ceil(remAmount / loan.EmiAmount)

	return int64(RemEMis)
}

func (loan *Loan) Balance(emi int64) *Result {
	result := &Result{
		AmountPaid: 0,
		RemEMi:     0,
	}
	if emi < 0 {
		return nil
	} else if emi > loan.NoOfEMI {
		result.AmountPaid = loan.TotalAmount
		result.RemEMi = 0
	} else {
		AmountPaidForNEMIs := loan.EmiAmount * float64(emi)
		AmountPaidByLumpSum := loan.getLumpSumTotalTillEmi(emi)

		result.AmountPaid = AmountPaidForNEMIs + AmountPaidByLumpSum
		result.RemEMi = loan.getRemainingEMIs(result.AmountPaid)
	}

	return result
}

func GetLoan(bank string, borrower string) *Loan {
	loanId := bank + "_" + borrower

	loan, ok := Loans[loanId]
	if !ok {
		return nil
	}

	return loan
}

func CreateLoan(bank string, borrower string, amount float64, tenure int64, rate float64) {
	newLoan := &Loan{
		Id:           bank + "_" + borrower,
		BankName:     bank,
		BorrowerName: borrower,
		Principal:    amount,
		Tenure:       tenure,
		Rate:         rate,
		TotalAmount:  0,
		NoOfEMI:      0,
		EmiAmount:    0,
		Repayments:   nil,
	}

	newLoan.init()
	Loans[newLoan.Id] = newLoan
}
