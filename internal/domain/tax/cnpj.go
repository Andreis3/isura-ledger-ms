package tax

import (
	"regexp"
	"slices"
	"strconv"

	"github.com/andreis3/isura-ledger-ms/internal/domain/validator"
)

var blackListCNPJ = []string{
	"00000000000000",
	"11111111111111",
	"22222222222222",
	"33333333333333",
	"44444444444444",
	"55555555555555",
	"66666666666666",
	"77777777777777",
	"88888888888888",
	"99999999999999",
}

type CNPJ struct {
	value string
}

// NewCNPJ tries to create the Value Object, validating it through its Evaluator.
func NewCNPJ(rawCNPJ string) (*CNPJ, validator.Evaluator) {
	var eval validator.Evaluator

	regex := regexp.MustCompile("[^0-9]")
	cleanCNPJ := regex.ReplaceAllString(rawCNPJ, "")

	eval.CheckField(validator.NotBlank(cleanCNPJ), "tax_id", "cannot be blank")
	eval.CheckField(len(cleanCNPJ) == 14, "tax_id", "must have exactly 14 characters")

	if len(cleanCNPJ) == 14 {
		eval.CheckField(!slices.Contains(blackListCNPJ, cleanCNPJ), "tax_id", "must be a valid CNPJ number")
		eval.CheckField(validateCNPJ(cleanCNPJ), "tax_id", "invalid digits calculated with module 11 algorithm")
	}

	if len(eval) > 0 {
		return nil, eval
	}

	return &CNPJ{value: cleanCNPJ}, eval
}

func (c *CNPJ) String() string {
	return c.value
}

func validateCNPJ(cnpj string) bool {
	size := len(cnpj) - 2
	numbers := cnpj[:size]
	digits := cnpj[size:]
	sum := 0
	pos := size - 7
	for i := size; i >= 1; i-- {
		convertNumber, _ := strconv.Atoi(string(numbers[size-i]))
		sum += convertNumber * pos
		pos--
		if pos < 2 {
			pos = 9
		}
	}
	result := 0
	if rest := sum % 11; rest < 2 {
		result = 0
	} else {
		result = 11 - (sum % 11)
	}
	if strconv.Itoa(result) != string(digits[0]) {
		return false
	}
	size++
	numbers = cnpj[:size]
	sum = 0
	pos = size - 7
	for i := size; i >= 1; i-- {
		convertNumber, _ := strconv.Atoi(string(numbers[size-i]))
		sum += convertNumber * pos
		pos--
		if pos < 2 {
			pos = 9
		}
	}
	if rest := sum % 11; rest < 2 {
		result = 0
	} else {
		result = 11 - (sum % 11)
	}
	if strconv.Itoa(result) != string(digits[1]) {
		return false
	}
	return true
}
