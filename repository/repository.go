package repository

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"password_generator/internal/cypher"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func CreateTable(db *sql.DB) {
	// Check if the table already exists
	var tableExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT FROM information_schema.tables WHERE table_name = 'passwords');").Scan(&tableExists)
	if err != nil {
		panic(fmt.Sprintf("Error checking if table exists: %v", err))
	}

	if !tableExists {
		_, err = db.Exec(`
            CREATE TABLE passwords (
                id SERIAL PRIMARY KEY,
                password TEXT UNIQUE NOT NULL,
                service TEXT NOT NULL
            );
        `)
		if err != nil {
			panic(fmt.Sprintf("Error creating table: %v", err))
		}
	}
}

func CheckConnection(db *sql.DB) {
	start := time.Now()
	for {
		if time.Since(start) > 20*time.Second {
			fmt.Println("Database connection timeout")
			break
		}
		err := db.Ping()
		if err == nil {
			break
		}
		fmt.Println("Waiting for database connection...")
		time.Sleep(2 * time.Second)
	}
}

func InsertIntoTable(db *sql.DB, note string, service string) {
	// Check if the password exists in the database
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM passwords WHERE password = $1)", note).Scan(&exists)
	if err != nil {
		panic(fmt.Sprintf("Error checking password existence: %v", err))
	}

	if exists {
		fmt.Println("Password already exists in the database")
	} else {
		// Insert a new encrypted record into the table
		encryptedPass, err := cypher.EncryptPassword(note, cypher.GetEncryptionKey())
		if err != nil {
			panic(fmt.Sprintf("Error encrypting password: %v", err))
		}
		_, err = db.Exec("INSERT INTO passwords (password, service) VALUES ($1, $2)", encryptedPass, service)
		if err != nil {
			panic(fmt.Sprintf("Error adding record: %v", err))
		}
	}
}

func ClearTable(db *sql.DB) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you really want to clear the table? (y/n): ")

	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	if text == "y" {
		_, err := db.Exec("DELETE FROM passwords")
		if err != nil {
			return err
		}

		_, err = db.Exec("ALTER SEQUENCE passwords_id_seq RESTART WITH 1")
		if err != nil {
			return err
		}

		fmt.Println("Table cleared")
	} else {
		fmt.Println("Operation canceled")
	}
	return nil
}

func PrintTableContent(db *sql.DB) {
	// Querying the table and printing its content
	rows, err := db.Query("SELECT * FROM passwords")
	if err != nil {
		panic(fmt.Sprintf("Error querying table: %v", err))
	}
	defer rows.Close()

	fmt.Println("Table 'passwords' content:")
	for rows.Next() {
		var id int
		var password string
		var service string

		err := rows.Scan(&id, &password, &service)
		if err != nil {
			panic(fmt.Sprintf("Error scanning row: %v", err))
		}
		// Print out all passwords if all of them match a decryption key
		decryptedPass, err := cypher.DecryptPassword(password, cypher.GetEncryptionKey())
		if err != nil {
			panic(fmt.Sprintf("Error decrypting password: %v", err))
		}

		fmt.Printf("ID: %d, Password: %s, Service: %s\n", id, decryptedPass, service)
	}

	if err := rows.Err(); err != nil {
		panic(fmt.Sprintf("Error during iteration: %v", err))
	}
}

func DeletePassword(db *sql.DB, primaryKey string) error {
	// Check if the password exists in the database
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM passwords WHERE id = $1)", primaryKey).Scan(&exists)
	if err != nil {
		panic(fmt.Sprintf("Error checking password existence: %v", err))
	}

	if !exists {
		fmt.Printf("Password with ID %s does not exist in the database", primaryKey)
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you really want to delete this password ", primaryKey, "? (y/n): ")

	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	if text == "y" {
		_, err := db.Exec("DELETE FROM passwords WHERE id = $1", primaryKey)
		if err != nil {
			panic(fmt.Sprintf("Error deleting password: %v", err))
		}
		fmt.Printf("Password %s deleted", primaryKey)
	} else {
		fmt.Println("Operation cancelled")
	}

	return nil
}

func DeleteLastPasswords(db *sql.DB, count string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you really want to delete the last ", count, " passwords? (y/n): ")

	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	if text == "y" {
		_, err := db.Exec("DELETE FROM passwords WHERE id IN (SELECT id FROM passwords ORDER BY id DESC LIMIT $1)", count)
		return err
	} else {
		fmt.Println("Operation canceled")
		return nil
	}
}

func GetPasswordsByService(db *sql.DB, service string) error {
	rows, err := db.Query("SELECT password FROM passwords WHERE service = $1", service)
	if err != nil {
		return err
	}
	defer rows.Close()

	var encryptedPasswords []string
	for rows.Next() {
		var encryptedPassword string
		err := rows.Scan(&encryptedPassword)
		if err != nil {
			return err
		}
		encryptedPasswords = append(encryptedPasswords, encryptedPassword)
	}

	var decryptedPasswords []string
	for _, encryptedPassword := range encryptedPasswords {
		decryptedPass, err := cypher.DecryptPassword(encryptedPassword, cypher.GetEncryptionKey())
		if err != nil {
			return fmt.Errorf("error decrypting password: %v", err)
		}
		decryptedPasswords = append(decryptedPasswords, decryptedPass)
	}

	fmt.Printf("Passwords for service %s:\n", service)
	for _, decryptedPassword := range decryptedPasswords {
		fmt.Println(decryptedPassword)
	}

	return nil
}

func GetPasswordById(db *sql.DB, id int) error {
	var password string
	err := db.QueryRow("SELECT password FROM passwords WHERE id = $1", id).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Password with ID %d not found\n", id)
			return nil
		}
		return fmt.Errorf("error querying password: %v", err)
	}

	decryptedPass, err := cypher.DecryptPassword(password, cypher.GetEncryptionKey())
	if err != nil {
		panic(fmt.Sprintf("Error decrypting password: %v", err))
	}

	fmt.Println("Password:", decryptedPass)
	return nil
}
