package clients

import "errors"

type Client struct {
	ID       string
	Capacity float64
	RPS      float64
}

func (c *Client) Validate() error {
	var errs []error
	if c.Capacity < 1 {
		errs = append(errs, ErrCapacityValidation)
	}
	if c.RPS <= 0 {
		errs = append(errs, ErrRPSValidation)
	}

	return errors.Join(errs...)
}
