package db

// Helper utilities for extracting a description out of a database

import (
	"database/sql"
	"encoding/json"
	"github.com/spekary/gengen/maps"
	"github.com/spekary/goradd/util"
	"log"
	"strconv"
	"strings"
)

const (
	SQL_TYPE_UNKNOWN  = "Unknown"
	SQL_TYPE_BLOB     = "Blob"
	SQL_TYPE_VARCHAR  = "VarChar"
	SQL_TYPE_CHAR     = "Char"
	SQL_TYPE_TEXT     = "Text"
	SQL_TYPE_INTEGER  = "Int"
	SQL_TYPE_DATETIME = "DateTime"
	SQL_TYPE_DATE     = "Date"
	SQL_TYPE_TIME     = "Time"
	SQL_TYPE_FLOAT    = "Float"
	SQL_TYPE_DOUBLE   = "Double"
	SQL_TYPE_BOOL     = "Bool"
	SQL_TYPE_DECIMAL  = "Decimal" // a fixed point type
)

func getTableDescription(tableName string, tableComment string, db SqlDbI) *TableDescription {
	var ok bool

	options, err := extractOptions(tableComment)

	if err != nil {
		log.Print("Error in table comment for table " + tableName + ": " + err.Error())
	}

	if v, _ := options.LoadBool("NoCodegen"); v {
		return nil
	}

	td := NewTableDescription(tableName)

	if td.EnglishName, ok = options.LoadString("englishName"); options.Has("englishName") && !ok {
		log.Print("Error in table comment for table " + tableName + ": englishName is not a string")
	}

	if td.EnglishPlural, ok = options.LoadString("englishPlural"); options.Has("englishPlural") && !ok {
		log.Print("Error in table comment for table " + tableName + ": EnglishPlural is not a string")
	}

	if td.GoName, ok = options.LoadString("goName"); options.Has("goName") && !ok {
		log.Print("Error in table comment for table " + tableName + ": goName is not a string")
	} else {
		td.GoName = strings.Title(td.GoName)
	}

	if td.GoPlural, ok = options.LoadString("goPlural"); options.Has("goPlural") && !ok {
		log.Print("Error in table comment for table " + tableName + ": goPlural is not a string")
	} else {
		td.GoPlural = strings.Title(td.GoName)
	}

	td.IsType = util.EndsWith(tableName, db.TypeTableSuffix())
	td.IsAssociation = util.EndsWith(tableName, db.AssociationTableSuffix())

	td.Comment = tableComment
	return td
}

// Find the json encoded list of options in the given string
func extractOptions(comment string) (options *maps.SliceMap, err error) {
	var optionString string
	firstIndex := strings.Index(comment, "{")
	lastIndex := strings.LastIndex(comment, "}")
	options = maps.NewSliceMap()

	if firstIndex != -1 &&
		lastIndex != -1 &&
		lastIndex > firstIndex {

		optionString = comment[firstIndex : lastIndex+1]

		err = json.Unmarshal([]byte(optionString), &options)
	}
	return
}

// Given a data definition description of the table, will extract the length from the definition
// If more than one number, returns the first number
// Example:
//	bigint(21) -> 21
// varchar(50) -> 50
// decimal(10,2) -> 10
func getDataDefLength(description string) int {
	var lastPos, lenPos int
	var size string
	if lenPos = strings.Index(description, "("); lenPos != -1 {
		lastPos = strings.LastIndex(description, ")")
		size = description[lenPos+1 : lastPos]
		sizes := strings.Split(size, ",")
		i, _ := strconv.Atoi(sizes[0])
		return i
	}
	return 0
}

// Retrieves a numeric value from the options, which is always going to return a float64
func getNumericOption(o *maps.SliceMap, option string, defaultValue float64) (float64, bool) {
	if v := o.Get(option); v != nil {
		if v2, ok := v.(float64); !ok {
			return defaultValue, false
		} else {
			return v2, true
		}
	} else {
		return defaultValue, true
	}
}

// Retrieves a boolean value from the options
func getBooleanOption(o *maps.SliceMap, option string) (val bool, ok bool) {
	val, ok = o.LoadBool(option)
	return
}

// Extracts a minimum and maximum value from the option map, returning defaults if none was found, and making sure
// the boundaries of anything found are not exceeded
func getMinMax(o *maps.SliceMap, defaultMin float64, defaultMax float64, tableName string, columnName string) (min float64, max float64) {
	var errString string

	if columnName == "" {
		errString = "table " + tableName
	} else {
		errString = "table " + tableName + ":" + columnName
	}

	v, ok := getNumericOption(o, "min", defaultMin)
	if !ok {
		log.Print("Error in min value in comment for " + errString + ". Value is not a valid number.")
		min = defaultMin
	} else {
		if v < defaultMin {
			log.Print("Error in min value in comment for " + errString + ". Value is less than the allowed minimum.")
			min = defaultMin
		} else {
			min = v
		}
	}

	v, ok = getNumericOption(o, "max", defaultMax)
	if !ok {
		log.Print("Error in max value in comment for " + errString + ". Value is not a valid number.")
		max = defaultMax
	} else {
		if v > defaultMax {
			log.Print("Error in max value in comment for " + errString + ". Value is more than the allowed maximum.")
			max = defaultMax
		} else {
			max = v
		}
	}

	return
}

func fkRuleToAction(rule sql.NullString) FKAction {

	if !rule.Valid {
		return FK_ACTION_NONE // This means we will emulate foreign key actions
	}
	switch strings.ToUpper(rule.String) {
	case "NO ACTION":
		fallthrough
	case "RESTRICT":
		return FK_ACTION_RESTRICT
	case "CASCADE":
		return FK_ACTION_CASCADE
	case "SET DEFAULT":
		return FK_ACTION_SET_DEFAULT
	case "SET NULL":
		return FK_ACTION_SET_NULL

	}
	return FK_ACTION_NONE
}
