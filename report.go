// Package artspeople provides a parser for the Arts People Line Item Reconciliation Report.  To create an instantiated LineItemReconReport struct, provide the text of the report's CSV representation into the NewLineItemReconReport() function.
package artspeople

import (
	"encoding/csv"
	"fmt"
	"io"
)

// A LineItemReconReport is a parsed and type-normalized version of the Line Item Reconciliation Report downloaded from Arts People.
type LineItemReconReport struct {
	rawLines [][]string
	Orders   map[int]*Order
}

// NewLineItemReconReport accepts an Arts People Line Item Reconciliation Report as an io.Reader and returns a fully parsed and type-normalized LineItemReconReport.
func NewLineItemReconReport(reportCsv io.Reader) (*LineItemReconReport, error) {
	lirReport := LineItemReconReport{}

	err := lirReport.readInput(reportCsv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input CSV: %v", err)
	}

	err = lirReport.parseRawLines()
	if err != nil {
		return nil, fmt.Errorf("failed to parse lines: %v", err)
	}

	return &lirReport, nil
}

func (l *LineItemReconReport) readInput(reportCsv io.Reader) error {
	lines, err := csv.NewReader(reportCsv).ReadAll()
	if err != nil {
		return err
	}

	// Throw away the first line, which contains the field headers.
	l.rawLines = lines[1:]
	return nil
}

func (l *LineItemReconReport) parseRawLines() error {
	if l.Orders == nil {
		l.Orders = map[int]*Order{}
	}

	for _, rl := range l.rawLines {
		// Use the raw line to create a LineItem.
		li, err := NewLineItem(rl)
		if err != nil {
			return fmt.Errorf("failed to parse line [%s]: %v", rl, err)
		}

		// If this is the first Line for a given order, instantiate a new Order struct.
		if _, ok := l.Orders[li.OrderID]; !ok {
			l.Orders[li.OrderID] = NewOrder()
		}

		// Add the LineItem to the Order with the LineItem's OrderID.
		l.Orders[li.OrderID].AddLineItem(li)
	}

	return nil
}

// GetItems returns a slice of all item names mentioned in the report.
func (l *LineItemReconReport) GetItems() []string {
	is := map[string]bool{}
	for _, o := range l.Orders {
		for _, i := range o.GetItems() {
			is[i] = true
		}
	}

	retItems := []string{}
	for i := range is {
		retItems = append(retItems, i)
	}

	return retItems
}

// GetCustomers returns a slice of all unique customers who have an order in the report.
func (l *LineItemReconReport) GetCustomers() []string {
	cs := map[string]bool{}
	for _, o := range l.Orders {
		cs[o.GetCustomer()] = true
	}

	retCusts := []string{}
	for c := range cs {
		retCusts = append(retCusts, c)
	}

	return retCusts
}
