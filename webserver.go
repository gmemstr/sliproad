package main

import (
	"fmt"
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
		panic(err)
	}
	err = yaml.Unmarshal(file, &files.ProviderConfig)
	if err != nil {
		panic(err)
	}
	files.SetupProviders()

	r := router.Init()
	fmt.Println("Your NAS instance is live on port :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
