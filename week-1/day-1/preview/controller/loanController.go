package controller

import (
	"net/http"
	"preview-w1/dto"
	"preview-w1/entity"
	"preview-w1/helpers"
	"preview-w1/repository"
	"preview-w1/utils"

	"github.com/labstack/echo/v4"
)

type LoanController struct {
	repository.DbHandler
}

func NewLoanController(dbHandler repository.DbHandler) LoanController {
	return LoanController{
		DbHandler: dbHandler,
	}
}

func (lc LoanController) ProposeLoan(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}
	
	var loanDataTmp dto.Loan
	if err := c.Bind(&loanDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.Details(err.Error()))
	}

	if err := c.Validate(&loanDataTmp); err != nil {
		return err
	}

	if loanDataTmp.Loan > loanDataTmp.Salary * 0.3 {
		return echo.NewHTTPError(utils.ErrBadRequest.Details("Salary is inproportional with the requested loan"))
	}

	loanData := entity.Loan{
		UserID: user.ID,
		Salary: loanDataTmp.Salary,
		Loan: loanDataTmp.Loan,
	}

	if err := lc.DbHandler.AddLoan(&loanData); err != nil {
		return err
	}
	
	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Loan created",
		Data: loanDataTmp,
	})
}

func (lc LoanController) Limit(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}
	
	limit, err := lc.DbHandler.CheckLimit(user.ID)
	if err != nil {
		return err
	}
	
	return c.JSON(http.StatusOK, dto.Response{
		Message: "Current limit",
		Data: limit,
	})
}

func (lc LoanController) Withdraw(c echo.Context) error {
	
	
	return nil
}

func (lc LoanController) Deposit(c echo.Context) error {
	return nil
}