/* webserver.go
 *
 * This is the webserver handler for Pogo, and handles
 * all incoming connections, including authentication.
 */

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gmemstr/nas/router"
)

// Main function that defines routes
func main() {
	// Define routes
	// We're live
	r := router.Init()
	fmt.Println("Your NAS instance is live on port :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
