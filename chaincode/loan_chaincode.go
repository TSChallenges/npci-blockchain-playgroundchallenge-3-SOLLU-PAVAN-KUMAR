package main
 
import (
	"encoding/json"
	"errors"
	"fmt"
 
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)
 
type LoanContract struct {
	contractapi.Contract
}
 
type Loan struct {
	LoanID        string    `json:"loanID"`
	ApplicantName string    `json:"applicantName"`
	LoanAmount    float64   `json:"loanAmount"`
	TermMonths    int       `json:"termMonths"`
	InterestRate  float64   `json:"interestRate"`
	Outstanding   float64   `json:"outstanding"`
	Status        string    `json:"status"`
	Repayments    []float64 `json:"repayments"`
}
 
// TODO: Implement ApplyForLoan
func (c *LoanContract) ApplyForLoan(ctx contractapi.TransactionContextInterface, loanID, applicantName string, loanAmount float64, termMonths int, interestRate float64) error {
	loanInfo, err := c.CheckLoanBalance(ctx, loanID)
	if err != nil && err.Error() != "Loan not found" {
		fmt.Println("error while fetching the loan info", err)
		return err
	}
	if loanInfo != nil {
		return errors.New("loan application already exists")
	}
 
	loanData := Loan{
		LoanID:        loanID,
		ApplicantName: applicantName,
		LoanAmount:    loanAmount,
		TermMonths:    termMonths,
		InterestRate:  interestRate,
		Outstanding:   loanAmount,
		Status:        "Applied",
		Repayments:    make([]float64, 0),
	}
	byteInfo, err := json.Marshal(loanData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}
	err = ctx.GetStub().PutState(loanID, byteInfo)
	if err != nil {
		return fmt.Errorf("failed to store data: %v", err)
	}
	return nil
}
 
// TODO: Implement ApproveLoan
func (c *LoanContract) ApproveLoan(ctx contractapi.TransactionContextInterface, loanID string, status string) error {
	loanInfo, err := c.CheckLoanBalance(ctx, loanID)
	if err != nil {
		fmt.Println("error while fetching the loan info", err)
		return err
	}
	if loanInfo.Status == "Approved" {
		return errors.New("loan application already approved")
	}
	loanInfo.Status = status
	byteInfo, err := json.Marshal(loanInfo)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}
	err = ctx.GetStub().PutState(loanID, byteInfo)
	if err != nil {
		return fmt.Errorf("failed to store data: %v", err)
	}
	return nil
}
 
// TODO: Implement MakeRepayment
func (c *LoanContract) MakeRepayment(ctx contractapi.TransactionContextInterface, loanID string, repaymentAmount float64) error {
	loanInfo, err := c.CheckLoanBalance(ctx, loanID)
	if err != nil {
		fmt.Println("error while fetching the loan info", err)
		return err
	}
	if loanInfo.Outstanding <= 0 {
		return errors.New("loan amount already paid")
	}
	loanInfo.Outstanding -= repaymentAmount
	loanInfo.Repayments = append(loanInfo.Repayments, repaymentAmount)
 
	byteInfo, err := json.Marshal(loanInfo)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}
	err = ctx.GetStub().PutState(loanID, byteInfo)
	if err != nil {
		return fmt.Errorf("failed to store data: %v", err)
	}
 
	return nil
}
 
// TODO: Implement CheckLoanBalance
func (c *LoanContract) CheckLoanBalance(ctx contractapi.TransactionContextInterface, loanID string) (*Loan, error) {
	data, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve data: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("Loan not found")
	}
	loadInfo := &Loan{}
	err = json.Unmarshal([]byte(data), loadInfo)
	if err != nil {
		return nil, err
	}
	return loadInfo, nil
}
 
func main() {
	chaincode, err := contractapi.NewChaincode(new(LoanContract))
	if err != nil {
		fmt.Printf("Error creating loan chaincode: %s", err)
		return
	}
 
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting loan chaincode: %s", err)
	}
}