package logic

import (
	"math"
	"testing"
)

func TestCalculateCollectionAmount(t *testing.T) {
	loan := Loan{
		Principal: 100000,
		Rate:      10,
		Tenure:    5,
	}

	loan.calculateCollectionAmount()

	// Expected value after calling calculateCollectionAmount()
	expectedTotalAmount := loan.Principal + ((loan.Principal * loan.Rate * float64(loan.Tenure)) / 100)

	if loan.TotalAmount != expectedTotalAmount {
		t.Errorf("Expected TotalAmount to be %.2f, but got %.2f", expectedTotalAmount, loan.TotalAmount)
	}
}

func TestCalculateEmi(t *testing.T) {
	loan := Loan{
		TotalAmount: 100000,
		Tenure:      5,
	}

	loan.calculateEmi()

	// Expected values after calling calculateEmi()
	expectedNoOfEMI := 60
	expectedEmiAmount := math.Ceil(100000 / float64(60))

	if int(loan.NoOfEMI) != expectedNoOfEMI {
		t.Errorf("Expected NoOfEMI to be %d, but got %d", expectedNoOfEMI, loan.NoOfEMI)
	}

	if loan.EmiAmount != expectedEmiAmount {
		t.Errorf("Expected EmiAmount to be %.2f, but got %.2f", expectedEmiAmount, loan.EmiAmount)
	}
}

func TestProcessRepayment(t *testing.T) {
	loan := Loan{}
	amount := 5000.0
	emi := int64(3)

	loan.ProcessRepayment(amount, emi)

	// Check the length of Repayments slice after calling ProcessRepayment()
	if len(loan.Repayments) != 1 {
		t.Errorf("Expected Repayments length to be 1, but got %d", len(loan.Repayments))
	}

	repayment := loan.Repayments[0]
	if repayment.Amount != amount {
		t.Errorf("Expected Repayment Amount to be %.2f, but got %.2f", amount, repayment.Amount)
	}
	if repayment.PostEMICount != emi {
		t.Errorf("Expected Repayment PostEMICount to be %d, but got %d", emi, repayment.PostEMICount)
	}
}

func TestGetLumpSumTotalTillEmi(t *testing.T) {
	loan := Loan{}
	repayments := []Repayment{
		{Amount: 1000.0, PostEMICount: 1},
		{Amount: 2000.0, PostEMICount: 2},
		{Amount: 3000.0, PostEMICount: 4},
	}
	loan.Repayments = repayments

	emi := int64(3)
	expectedLumpSumTotal := float64(3000.0)

	lumpSumTotal := loan.getLumpSumTotalTillEmi(emi)

	if lumpSumTotal != expectedLumpSumTotal {
		t.Errorf("Expected LumpSumTotal to be %.2f, but got %.2f", expectedLumpSumTotal, lumpSumTotal)
	}
}

func TestGetRemainingEMIs(t *testing.T) {
	loan := Loan{}
	loan.TotalAmount = 10000.0
	loan.EmiAmount = 2400.0

	amountPaid := 5000.0
	// Expected value after calling getRemainingEMIs()
	expectedRemainingEMIs := int64(math.Ceil((10000.0 - 5000.0) / 2400.0))

	remainingEMIs := loan.getRemainingEMIs(amountPaid)

	if remainingEMIs != expectedRemainingEMIs {
		t.Errorf("Expected RemainingEMIs to be %d, but got %d", expectedRemainingEMIs, remainingEMIs)
	}
}

func TestGetLoan(t *testing.T) {
	// Create test data
	bank := "BankA"
	borrower := "JohnDoe"
	loanId := bank + "_" + borrower
	expectedLoan := &Loan{Id: loanId}

	// Add the loan to Loans map for testing
	Loans[loanId] = expectedLoan

	// Test the GetLoan function
	loan := GetLoan(bank, borrower)

	//Check if returned loan with the expected value
	if loan == nil {
		t.Errorf("Expected loan to be not nil, but got nil")
	} else if loan.Id != expectedLoan.Id {
		t.Errorf("Expected loan ID to be %s, but got %s", expectedLoan.Id, loan.Id)
	}
}

func TestCreateLoan(t *testing.T) {
	// Create test data
	bank := "BankA"
	borrower := "JohnDoe"
	amount := 1000.0
	tenure := int64(12)
	rate := 10.0

	// Call the CreateLoan function
	CreateLoan(bank, borrower, amount, tenure, rate)

	// Assert the loan is created in Loans map
	loanId := bank + "_" + borrower
	loan, ok := Loans[loanId]
	if !ok {
		t.Errorf("Expected loan to be created in Loans map, but not found")
	}

	// Assert the loan properties are set correctly
	if loan == nil {
		t.Errorf("Expected loan to be not nil, but got nil")
	} else {
		if loan.Id != loanId {
			t.Errorf("Expected loan ID to be %s, but got %s", loanId, loan.Id)
		}
		if loan.BankName != bank {
			t.Errorf("Expected bank name to be %s, but got %s", bank, loan.BankName)
		}
		if loan.BorrowerName != borrower {
			t.Errorf("Expected borrower name to be %s, but got %s", borrower, loan.BorrowerName)
		}
		if loan.Principal != amount {
			t.Errorf("Expected principal amount to be %f, but got %f", amount, loan.Principal)
		}
		if loan.Tenure != tenure {
			t.Errorf("Expected tenure to be %d, but got %d", tenure, loan.Tenure)
		}
		if loan.Rate != rate {
			t.Errorf("Expected rate to be %f, but got %f", rate, loan.Rate)
		}
		if loan.Repayments != nil {
			t.Errorf("Expected repayments to be nil, but got %+v", loan.Repayments)
		}
	}
}

func TestBalanceWithNegativeEMI(t *testing.T) {
	// Create a loan object
	loan := &Loan{
		NoOfEMI:     int64(12),
		EmiAmount:   100.0,
		Repayments:  nil,
		TotalAmount: 1200.0,
	}

	// Call the Balance function with negative EMI
	result := loan.Balance(-1)

	// Check if result is nil
	if result != nil {
		t.Errorf("Expected result to be nil, but got %+v", result)
	}
}

func TestBalanceWithZeroEMI(t *testing.T) {
	// Create a loan object
	loan := &Loan{
		NoOfEMI:     int64(12),
		EmiAmount:   100.0,
		Repayments:  nil,
		TotalAmount: 1200.0,
	}

	// Call the Balance function with zero EMI
	result := loan.Balance(0)

	// Check if result has correct values
	if result == nil {
		t.Errorf("Expected result to be not nil, but got nil")
	} else {
		if result.AmountPaid != 0 {
			t.Errorf("Expected amount paid to be 0, but got %f", result.AmountPaid)
		}
		if result.RemEMi != loan.NoOfEMI {
			t.Errorf("Expected remaining EMIs to be %d, but got %d", loan.NoOfEMI, result.RemEMi)
		}
	}

	// Test for 0 emi when there is repayment
	repayment := Repayment{
		Amount:       200,
		PostEMICount: 0,
	}
	loan.Repayments = append(loan.Repayments, repayment)

	// Call the Balance function with zero EMI
	result2 := loan.Balance(0)

	// Check if result has correct values
	if result2 == nil {
		t.Errorf("Expected result to be not nil, but got nil")
	} else {
		if result2.AmountPaid != loan.getLumpSumTotalTillEmi(0) {
			t.Errorf("Expected amount paid to be 0, but got %f", result2.AmountPaid)
		}
		if result2.RemEMi != loan.getRemainingEMIs(200) {
			t.Errorf("Expected remaining EMIs to be %d, but got %d", loan.NoOfEMI, result2.RemEMi)
		}
	}
}

func TestBalanceWithEMIGreaterThanNoOfEMI(t *testing.T) {
	// Create a loan object
	loan := &Loan{
		NoOfEMI:     int64(12),
		EmiAmount:   100.0,
		Repayments:  nil,
		TotalAmount: 1200.0,
	}

	// Call the Balance function with EMI greater than number of EMIs
	result := loan.Balance(15)

	// C that result has correct values
	if result == nil {
		t.Errorf("Expected result to be not nil, but got nil")
	} else {
		if result.AmountPaid != loan.TotalAmount {
			t.Errorf("Expected amount paid to be %f, but got %f", loan.TotalAmount, result.AmountPaid)
		}
		if result.RemEMi != 0 {
			t.Errorf("Expected remaining EMIs to be 0, but got %d", result.RemEMi)
		}
	}
}

func TestBalanceWithValidEMI(t *testing.T) {
	// Create repayment objects to add to Repayments array of loam
	repayment1 := Repayment{
		Amount:       150,
		PostEMICount: 3,
	}

	repayment2 := Repayment{
		Amount:       300,
		PostEMICount: 8,
	}

	// Create a loan object
	loan := &Loan{
		NoOfEMI:     int64(12),
		EmiAmount:   100.0,
		Repayments:  nil,
		TotalAmount: 1200.0,
	}

	loan.Repayments = append(loan.Repayments, repayment1, repayment2)

	// Call the Balance function with valid EMI
	result := loan.Balance(5)

	// Check if result has correct values
	if result == nil {
		t.Errorf("Expected result to be not nil, but got nil")
	} else {
		expectedAmountPaid := loan.EmiAmount*float64(5) + loan.getLumpSumTotalTillEmi(5)
		expectedRemEMI := loan.getRemainingEMIs(expectedAmountPaid)

		if result.AmountPaid != expectedAmountPaid {
			t.Errorf("Expected amount paid to be %f, but got %f", expectedAmountPaid, result.AmountPaid)
		}

		if result.RemEMi != expectedRemEMI {
			t.Errorf("Expected remaining EMIs to be %d, but got %d", expectedRemEMI, result.RemEMi)
		}
	}
}
