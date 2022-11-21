package oddslib

import (
	"errors"
	"fmt"
)

// ErrOddslib : Base error for this library
var ErrOddslib = errors.New("oddslib error")

// DBError : Error raised while trying to interract with a database
type DBError struct {
	err error
}
func (DBError) Unwrap() error {
	return ErrOddslib
}
func (e DBError) Error() string {
	return fmt.Sprintf("An error occurred while working with the db : %v", e.err)
}


// ErrFile : Error raised while trying to interract with a database
type ErrFile struct {
	err error
}
func (ErrFile) Unwrap() error {
	return ErrOddslib
}
func (e ErrFile) Error() string {
	return fmt.Sprintf("An error occurred while working with a file : %v", e.err)
}


// ErrSolving : Error raised while trying to solve the problem
var ErrSolving = errors.New("error during solving")
