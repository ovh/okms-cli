package config

import (
	"errors"
	"net"
	"net/url"
	"os"

	"github.com/ovh/okms-cli/internal/utils"
)

type Validator func(string) error

func (v Validator) AllowEmpty() Validator {
	return func(s string) error {
		if s == "" {
			return nil
		}
		return v(s)
	}
}

var (
	ValidateURL Validator = func(s string) error {
		u, err := url.Parse(s)
		if err != nil {
			return err
		}
		if _, err = net.ResolveIPAddr("ip", u.Hostname()); err != nil {
			return err
		}
		return nil
	}
	ValidateFileExists Validator = func(s string) error {
		s, err := utils.ExpandTilde(s)
		if err != nil {
			return err
		}
		inf, err := os.Stat(s)
		if err != nil {
			return err
		}
		if inf.IsDir() {
			return errors.New("Must be a file but is a directory")
		}
		return nil
	}
	ValidateTCPAddr Validator = func(s string) error {
		_, err := net.ResolveTCPAddr("tcp", s)
		return err
	}
)

func checkValidate(value string, validate ...Validator) error {
	var errs []error
	for _, isValid := range validate {
		if err := isValid(value); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
