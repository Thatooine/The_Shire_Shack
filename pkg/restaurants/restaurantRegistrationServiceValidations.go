package restaurants

import (
	"fmt"
	"strings"
)

func (r *RegisterRestaurantRequest) Validate() error {
	var reasons []string

	if r.Name == "" {
		reasons = append(reasons, "Name is required")
	}

	if r.City == "" {
		reasons = append(reasons, "City is required")
	}

	if len(reasons) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(reasons, "; "))
	}

	return nil
}
