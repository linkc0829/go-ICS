// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"fmt"
	"io"
	"strconv"
)

type CostCategory string

const (
	CostCategoryInvestment CostCategory = "INVESTMENT"
	CostCategoryDaily      CostCategory = "DAILY"
	CostCategoryLearning   CostCategory = "LEARNING"
	CostCategoryCharity    CostCategory = "CHARITY"
	CostCategoryOthers     CostCategory = "OTHERS"
)

var AllCostCategory = []CostCategory{
	CostCategoryInvestment,
	CostCategoryDaily,
	CostCategoryLearning,
	CostCategoryCharity,
	CostCategoryOthers,
}

func (e CostCategory) IsValid() bool {
	switch e {
	case CostCategoryInvestment, CostCategoryDaily, CostCategoryLearning, CostCategoryCharity, CostCategoryOthers:
		return true
	}
	return false
}

func (e CostCategory) String() string {
	return string(e)
}

func (e *CostCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CostCategory(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CostCategory", str)
	}
	return nil
}

func (e CostCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type IncomeCategory string

const (
	IncomeCategoryInvestment IncomeCategory = "INVESTMENT"
	IncomeCategoryParttime   IncomeCategory = "PARTTIME"
	IncomeCategorySalary     IncomeCategory = "SALARY"
	IncomeCategoryOthers     IncomeCategory = "OTHERS"
)

var AllIncomeCategory = []IncomeCategory{
	IncomeCategoryInvestment,
	IncomeCategoryParttime,
	IncomeCategorySalary,
	IncomeCategoryOthers,
}

func (e IncomeCategory) IsValid() bool {
	switch e {
	case IncomeCategoryInvestment, IncomeCategoryParttime, IncomeCategorySalary, IncomeCategoryOthers:
		return true
	}
	return false
}

func (e IncomeCategory) String() string {
	return string(e)
}

func (e *IncomeCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IncomeCategory(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IncomeCategory", str)
	}
	return nil
}

func (e IncomeCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PortfolioCategory string

const (
	PortfolioCategoryInvestment PortfolioCategory = "INVESTMENT"
	PortfolioCategoryParttime   PortfolioCategory = "PARTTIME"
	PortfolioCategorySalary     PortfolioCategory = "SALARY"
	PortfolioCategoryDaily      PortfolioCategory = "DAILY"
	PortfolioCategoryLearning   PortfolioCategory = "LEARNING"
	PortfolioCategoryCharity    PortfolioCategory = "CHARITY"
	PortfolioCategoryOthers     PortfolioCategory = "OTHERS"
)

var AllPortfolioCategory = []PortfolioCategory{
	PortfolioCategoryInvestment,
	PortfolioCategoryParttime,
	PortfolioCategorySalary,
	PortfolioCategoryDaily,
	PortfolioCategoryLearning,
	PortfolioCategoryCharity,
	PortfolioCategoryOthers,
}

func (e PortfolioCategory) IsValid() bool {
	switch e {
	case PortfolioCategoryInvestment, PortfolioCategoryParttime, PortfolioCategorySalary, PortfolioCategoryDaily, PortfolioCategoryLearning, PortfolioCategoryCharity, PortfolioCategoryOthers:
		return true
	}
	return false
}

func (e PortfolioCategory) String() string {
	return string(e)
}

func (e *PortfolioCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PortfolioCategory(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PortfolioCategory", str)
	}
	return nil
}

func (e PortfolioCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Privacy string

const (
	PrivacyPublic  Privacy = "PUBLIC"
	PrivacyFriend  Privacy = "FRIEND"
	PrivacyPrivate Privacy = "PRIVATE"
)

var AllPrivacy = []Privacy{
	PrivacyPublic,
	PrivacyFriend,
	PrivacyPrivate,
}

func (e Privacy) IsValid() bool {
	switch e {
	case PrivacyPublic, PrivacyFriend, PrivacyPrivate:
		return true
	}
	return false
}

func (e Privacy) String() string {
	return string(e)
}

func (e *Privacy) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Privacy(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Privacy", str)
	}
	return nil
}

func (e Privacy) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

var AllRole = []Role{
	RoleAdmin,
	RoleUser,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleAdmin, RoleUser:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
