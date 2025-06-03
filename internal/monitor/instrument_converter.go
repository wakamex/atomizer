package monitor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// InstrumentConverter converts between different exchange naming conventions
type InstrumentConverter struct {
	// Regex patterns for different formats
	derivePattern  *regexp.Regexp // ETH-20250531-2700-C
	deribitPattern *regexp.Regexp // ETH-31MAY25-2700-C
}

func NewInstrumentConverter() *InstrumentConverter {
	return &InstrumentConverter{
		derivePattern:  regexp.MustCompile(`^([A-Z]+)-(\d{8})-(\d+)-([CP])$`),
		deribitPattern: regexp.MustCompile(`^([A-Z]+)-(\d{1,2})([A-Z]{3})(\d{2})-(\d+)-([CP])$`),
	}
}

// ConvertForExchange converts an instrument name to the appropriate format for the given exchange
func (ic *InstrumentConverter) ConvertForExchange(instrument, exchange string) string {
	// Handle perpetuals and futures (no conversion needed)
	if strings.Contains(instrument, "PERPETUAL") || strings.Contains(instrument, "PERP") {
		return instrument
	}

	// Try to parse as Derive format first
	if matches := ic.derivePattern.FindStringSubmatch(instrument); matches != nil {
		if exchange == "derive" {
			return instrument // Already in correct format
		}
		// Convert to Deribit format
		return ic.deriveToDeribit(matches[1], matches[2], matches[3], matches[4])
	}

	// Try to parse as Deribit format
	if matches := ic.deribitPattern.FindStringSubmatch(instrument); matches != nil {
		if exchange == "deribit" {
			return instrument // Already in correct format
		}
		// Convert to Derive format
		return ic.deribitToDerive(matches[1], matches[2], matches[3], matches[4], matches[5], matches[6])
	}

	// Return original if no pattern matches (might be a simple pattern like "ETH")
	return instrument
}

// deriveToDeribit converts from Derive format (ETH-20250531-2700-C) to Deribit format (ETH-31MAY25-2700-C)
func (ic *InstrumentConverter) deriveToDeribit(asset, dateStr, strike, optionType string) string {
	// Parse YYYYMMDD
	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "" // Invalid date
	}

	// Format as DDMMMYY
	deribitDate := strings.ToUpper(date.Format("2Jan06"))

	return fmt.Sprintf("%s-%s-%s-%s", asset, deribitDate, strike, optionType)
}

// deribitToDerive converts from Deribit format (ETH-31MAY25-2700-C) to Derive format (ETH-20250531-2700-C)
func (ic *InstrumentConverter) deribitToDerive(asset, day, month, year, strike, optionType string) string {
	// Convert 2-digit year to 4-digit
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return ""
	}
	fullYear := 2000 + yearInt

	// Parse the date
	dateStr := fmt.Sprintf("%s%s%d", day, month, fullYear)
	date, err := time.Parse("2Jan2006", dateStr)
	if err != nil {
		return "" // Invalid date
	}

	// Format as YYYYMMDD
	deriveDate := date.Format("20060102")

	return fmt.Sprintf("%s-%s-%s-%s", asset, deriveDate, strike, optionType)
}

// ConvertInstrumentList converts a list of instruments for a specific exchange
func (ic *InstrumentConverter) ConvertInstrumentList(instruments []string, exchange string) []string {
	converted := make([]string, 0, len(instruments))
	seen := make(map[string]bool)

	for _, inst := range instruments {
		convertedInst := ic.ConvertForExchange(inst, exchange)
		if convertedInst != "" && !seen[convertedInst] {
			converted = append(converted, convertedInst)
			seen[convertedInst] = true
		}
	}

	return converted
}
