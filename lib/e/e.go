package e

import "fmt"

func Wrap(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("error message: %s, error: %w", msg, err)
	}
	return err
}
