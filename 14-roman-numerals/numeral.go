package romannumerals

import "strings"

func ConvertToArabic(roman string) (total uint16) {
	for _, symbols := range widowedRoman(roman).Symbols() {
		total += allRomansNumerals.ValueOf(symbols...)
	}

	return
}

func ConvertToRoman(arabic uint16) string {
	var result strings.Builder

	for _, numeral := range allRomansNumerals {
		for arabic >= numeral.Value {
			result.WriteString(numeral.Symbol)
			arabic -= numeral.Value
		}
	}

	return result.String()
}

type romanNumeral struct {
	Value  uint16
	Symbol string
}

type romanNumerals []romanNumeral

func (r romanNumerals) ValueOf(symbols ...byte) uint16 {
	symbol := string(symbols)

	for _, s := range r {
		if s.Symbol == symbol {
			return s.Value
		}
	}

	return 0
}

func (r romanNumerals) Exists(symbols ...byte) bool {
	symbol := string(symbols)

	for _, s := range r {
		if s.Symbol == symbol {
			return true
		}
	}

	return false
}

var allRomansNumerals = romanNumerals{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

type widowedRoman string

func (w widowedRoman) Symbols() (symbols [][]byte) {
	for i := 0; i < len(w); i += 1 {
		symbol := w[i]
		notAtEnd := i+1 < len(w)

		if notAtEnd && isSubtractive(symbol) && allRomansNumerals.Exists(symbol, w[i+1]) {
			symbols = append(symbols, []byte{symbol, w[i+1]})
			i += 1
		} else {
			symbols = append(symbols, []byte{symbol})
		}
	}

	return
}

func isSubtractive(currentSymbol uint8) bool {
	return currentSymbol == 'I' || currentSymbol == 'X' || currentSymbol == 'C'
}
