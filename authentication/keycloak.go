package authentication

import (
	"fmt"
	"github.com/Nerzal/gocloak/v5"
)

var AuthConfig map[string]string

func HasAuth(accessToken string) (success bool) {
	client := gocloak.NewClient(AuthConfig["provider_url"])
	_, err := client.GetUserInfo(accessToken, AuthConfig["realm"])
	if err != nil {
		return false
	}
	return true
}

func GetLoginLink() (url string) {
	baseString := "%v/auth/realms/%v/protocol/openid-connect/auth?client_id=account&response_mode=fragment&response_type=token&login=true&redirect_uri=%v/api/auth/callback"
	authUrl := fmt.Sprintf(baseString, AuthConfig["provider_url"], AuthConfig["realm"], AuthConfig["redirect_base_url"])
	return authUrl
}