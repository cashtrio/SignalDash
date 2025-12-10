package sderror

import "errors"

var (
	ErrNoSuchRecord   = errors.New("Can't find the record from database")
	ErrEmptyStorage   = errors.New("No records to fetch in database")
	ErrInternalServer = errors.New("Something went wrong in server side")
)
