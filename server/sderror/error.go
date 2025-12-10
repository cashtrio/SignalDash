package sderror

import "errors"

var (
	NoSuchRecord = errors.New("Can't find the record from database")
	EmptyStorage = errors.New("No records to fetch in database")
)
