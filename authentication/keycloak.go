package authentication

import (
	"fmt"

	"github.com/Nerzal/gocloak/v5"
)

// AuthConfig contains the configuration for the IdP.
var AuthConfig map[string]string

// HasAuth checks the passed token against the IdP, and returns true
// if the IdP can return the user's info, false if not.
func HasAuth(accessToken string) (success bool) {
	client := gocloak.NewClient(AuthConfig["provider_url"])
	_, err := client.GetUserInfo(accessToken, AuthConfig["realm"])
	if err != nil {
		return false
	}
	return true
}

// GetLoginLink generates a redirect link to the IdP login page.
func GetLoginLink() (url string) {
	baseString := "%v/auth/realms/%v/protocol/openid-connect/auth?client_id=account&response_mode=fragment&response_type=token&login=true&redirect_uri=%v/api/auth/callback"
	authURL := fmt.Sprintf(baseString, AuthConfig["provider_url"], AuthConfig["realm"], AuthConfig["redirect_base_url"])
	return authURL
}
