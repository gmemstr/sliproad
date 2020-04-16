package main

import (
	"fmt"
	"github.com/gmemstr/nas/authentication"
	"github.com/gmemstr/nas/files"
	"github.com/gmemstr/nas/router"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"net/http"
)

// Main function that defines routes
func main() {
	// Initialize file providers.
	file, err := ioutil.ReadFile("providers.yml")
	if err != nil {
		fmt.Println("Unable to read providers.yml file, does it exist?")
	}
	err = yaml.Unmarshal(file, &files.ProviderConfig)
	if err != nil {
		fmt.Println("Unable to parse providers.yml file, is it correct?")
	}
	files.SetupProviders()

	// Initialize auth if set up.
	authConfig, err := ioutil.ReadFile("auth.yml")
	if err != nil {
		fmt.Println("!! No Keycloack configuration found !!\n!! Requests will be unauthenticated !!")
		router.AuthEnabled = false
	} else {
		err = yaml.Unmarshal(authConfig, &authentication.AuthConfig)
		if err != nil {
			fmt.Println("Unable to parse auth.yml file, is it correct?")
			router.AuthEnabled = false
		}
		fmt.Println("Keycloak configured")
	}

	r := router.Init()
	fmt.Println("Your NAS instance is live on port :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
