/* webserver.go
 *
 * This is the webserver handler for Pogo, and handles
 * all incoming connections, including authentication.
 */

package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/gmemstr/nas/files"
	"github.com/go-yaml/yaml"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gmemstr/nas/router"
)

// Main function that defines routes
func main() {
	if _, err := os.Stat(".lock"); os.IsNotExist(err) {
		createDatabase()
		createLockFile()
	}

	file, err := ioutil.ReadFile("providers.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &files.Providers)
	if err != nil {
		panic(err)
	}
	fmt.Println(files.Providers)

	r := router.Init()
	fmt.Println("Your NAS instance is live on port :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func createDatabase() {
	fmt.Println("Initializing the database")
	os.Create("assets/config/users.db")

	db, err := sql.Open("sqlite3", "assets/config/users.db")
	if err != nil {
		fmt.Println("Problem opening database file! %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `users` ( `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, `username` TEXT UNIQUE, `hash` TEXT, `token` TEXT, `permissions` INTEGER )")
	if err != nil {
		fmt.Println("Problem creating database! %v", err)
	}

	text, err := GenerateRandomString(12)
	if err != nil {
		fmt.Println("Error randomly generating password", err)
	}
	fmt.Println("Admin password: ", text)
	hash, err := bcrypt.GenerateFromPassword([]byte(text), 4)
	if err != nil {
		fmt.Println("Error generating hash", err)
	}
	if bcrypt.CompareHashAndPassword(hash, []byte(text)) == nil {
		fmt.Println("Password hashed")
	}
	_, err = db.Exec("INSERT INTO users(id,username,hash,permissions) VALUES (0,'admin','" + string(hash) + "',2)")
	if err != nil {
		fmt.Println("Problem creating database! %v", err)
	}
	defer db.Close()
}

func createLockFile() {
	lock, err := os.Create(".lock")
	if err != nil {
		fmt.Println("Error: %v", err)
	}
	lock.Write([]byte("This file left intentionally empty"))
	defer lock.Close()
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}


// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}