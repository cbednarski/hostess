package hostess

import (
	"errors"
)

func Add(args []string, flags map[string]*bool) error {
	if len(args) > 2 {
		return errors.New("Unexpected arguments")
	}
	return nil
}

func Del(args []string, flags map[string]*bool) error {
	return nil
}

func Has(args []string, flags map[string]*bool) error {
	return nil
}

func Off(args []string, flags map[string]*bool) error {
	return nil
}

func On(args []string, flags map[string]*bool) error {
	return nil
}

func Ls(args []string, flags map[string]*bool) error {
	return nil
}

func Fix(args []string, flags map[string]*bool) error {
	return nil
}

func Dump(args []string, flags map[string]*bool) error {
	return nil
}

func Apply(args []string, flags map[string]*bool) error {
	return nil
}
