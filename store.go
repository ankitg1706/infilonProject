package main

import (
	"database/sql"
	"log"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", "ankit:password@tcp(localhost:3306)/infilon")
	if err != nil {
		log.Fatal(err)
	}
}

func getPersonByID(personID int) (Person, error) {
	// Query to fetch person's info, phone number, and address using JOINs
	query := `
        SELECT p.name, ph.number AS phone_number, a.city, a.state, a.street1, a.street2, a.zip_code
        FROM person p
        JOIN phone ph ON p.id = ph.person_id
        JOIN address_join aj ON p.id = aj.person_id
        JOIN address a ON aj.address_id = a.id
        WHERE p.id = ?
    `
	var personInfo Person
	err := db.QueryRow(query, personID).Scan(
		&personInfo.Name,
		&personInfo.PhoneNumber,
		&personInfo.City,
		&personInfo.State,
		&personInfo.Street1,
		&personInfo.Street2,
		&personInfo.ZipCode,
	)
	return personInfo, err
}

func insertPerson(person Person) error {
    tx, err := db.Begin()
    if err != nil {
        log.Println("Error beginning transaction:", err)
        return err
    }

    result, err := tx.Exec(`
        INSERT INTO person(name) VALUES (?)
    `, person.Name)
    if err != nil {
        tx.Rollback()  
        log.Println("Error inserting into person table:", err)
        return err
    }

    personID, err := result.LastInsertId()
    if err != nil {
        tx.Rollback()  
        log.Println("Error getting last insert ID:", err)
        return err
    }

    _, err = tx.Exec(`
        INSERT INTO phone(person_id, number) VALUES (?, ?)
    `, personID, person.PhoneNumber)
    if err != nil {
        tx.Rollback()  
        log.Println("Error inserting into phone table:", err)
        return err
    }

    result, err = tx.Exec(`
        INSERT INTO address(city, state, street1, street2, zip_code) VALUES (?, ?, ?, ?, ?)
    `, person.City, person.State, person.Street1, person.Street2, person.ZipCode)
    if err != nil {
        tx.Rollback()  
        log.Println("Error inserting into address table:", err)
        return err
    }

    
    addressID, err := result.LastInsertId()
    if err != nil {
        tx.Rollback()  
        log.Println("Error getting last insert ID:", err)
        return err
    }

    _, err = tx.Exec(`
        INSERT INTO address_join(person_id, address_id) VALUES (?, ?)
    `, personID, addressID)
    if err != nil {
        tx.Rollback()  
        log.Println("Error inserting into address_join table:", err)
        return err
    }

    if err := tx.Commit(); err != nil {
        log.Println("Error committing transaction:", err)
        return err
    }

    return nil
}

