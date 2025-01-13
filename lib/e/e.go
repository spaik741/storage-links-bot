package e

import "fmt"

func Wrap(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("Error message: %s, error: %w", msg, err)
	}
	return err
}
