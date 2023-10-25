package file

import (
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"
)

var (
	ColumnTitleName = "name"
)

func GetAllNamesFromSheet(fileName, sheetName string) ([]string, error) {

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", fileName, err)
	}

	defer f.Close()

	// get all rows from sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s sheet: %v", sheetName, err)
	}

	// check the rows is empty or not
	if len(rows) == 0 {
		return nil, errors.New("sheet doesn't contain any rows")
	}
	// then check the first row(titles) have values or not
	if len(rows[0]) == 0 {
		return nil, errors.New("sheet doesn't contain titles")
	}

	// find the column 'name' index
	columnNameIdx := -1
	for i := range rows[0] { // range the title and find the name column title index
		if rows[0][i] == ColumnTitleName {
			columnNameIdx = i
			break
		}
	}

	// check the column index found or not
	if columnNameIdx == -1 {
		return nil, fmt.Errorf("the sheet doesn't contain column title '%s'", ColumnTitleName)
	}

	// make string slice of row length(exclude the title)
	names := make([]string, len(rows)-1)

	for i := 1; i < len(rows); i++ {
		// set name the names slice
		names[i-1] = rows[i][columnNameIdx]
	}

	return names, nil
}
