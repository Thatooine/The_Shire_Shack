package users

import (
	"fmt"
	"strings"
)

func (r *CreateUserRequest) Validate() error {
	var reasons []string

	if r.Name == "" {
		reasons = append(reasons, "Name is required")
	}

	if r.Email == "" {
		reasons = append(reasons, "Email is required")
	}

	if len(reasons) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(reasons, "; "))
	}

	return nil
}
