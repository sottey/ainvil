package api

import (
	"fmt"
	"strings"

	"github.com/sottey/ainvil/lib/api/limitless"
	"github.com/sottey/ainvil/lib/api/omi"
	"github.com/sottey/ainvil/lib/clientiface"
)

// GetClient returns an initialized API client based on the name and token.
func GetClient(name, token string) (clientiface.APIClient, error) {
	switch strings.ToLower(name) {
	case "omi":
		return omi.NewClient(token), nil
	case "limitless":
		return limitless.NewClient(token), nil
	default:
		return nil, fmt.Errorf("unsupported API: %s", name)
	}
}
