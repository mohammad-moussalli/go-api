package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type users struct {
	id         int
	first_name string
	last_name  string
	dob        string
	email      string
	password   string
	picture    string
	timestamp  time.Time
	address_id int
}

type posts struct {
	post_id   int
	post      string
	user_id   int
	timestamp *time.Time
}

type likes struct {
	like_id int
	user_id int
	post_id int
}

type friendships struct {
	friendship_id int
	sender        int
	receiver      int
	accepted      int
	timestamp     time.Time
}

type blocks struct {
	block_id int
	sender   int
	receiver int
}

type addresses struct {
	address_id int
	country    string
	city       string
	street     string
}

func insertPost(db *sql.DB, p posts) error {
	insert, err = db.Query("INSERT INTO posts(post, user_id) VALUES (?, ?)")
	query := "INSERT INTO posts(post, user_id) VALUES (?, ?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, p.post, p.user_id)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d post created ", rows)
	return nil
}

func getPosts(db *sql.DB, id int) {

	res, err := db.Query("SELECT DISTINCT posts.post, posts.timestamp, posts.user_id, users.first_name, users.last_name, users.picture FROM posts JOIN users ON posts.user_id = users.id JOIN friendships ON (posts.user_id = friendships.sender OR posts.user_id = friendships.receiver) WHERE(friendships.sender = ? OR friendships.receiver = ?) AND friendships.accepted = 1 AND users.id NOT IN (SELECT blocks.receiver FROM blocks WHERE blocks.receiver = ? OR blocks.sender = ?) ORDER BY timestamp DESC;", id, id, id, id)

	if err != nil {
		panic(err.Error())
	}

	for res.Next() {
		var posts posts
		var users users

		err = res.Scan(&posts.post, &posts.timestamp, &posts.user_id, &users.first_name, &users.last_name, &users.picture)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(posts.post, posts.timestamp, posts.user_id, users.first_name, users.last_name, users.picture)
	}
}

func getUserData(db *sql.DB, id int) {

	res, err := db.Query("SELECT first_name, last_name, dob, email, picture, addresses.country, addresses.city, addresses.street FROM users JOIN addresses ON  users.address_id = addresses.address_id WHERE id = ?", id)

	if err != nil {
		panic(err.Error())
	}

	for res.Next() {
		var users users
		var addresses addresses
		err = res.Scan(&users.first_name, &users.last_name, &users.dob, &users.email, &users.picture, &addresses.country, &addresses.city, &addresses.street)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(users.first_name, users.last_name, users.dob, users.email, users.picture, addresses.country, addresses.city, addresses.street)
	}
}

func main() {
	fmt.Println("Go MySQL Tutorial")

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/facebookdb?parseTime=true")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	//getPosts(db, 145)
	getUserData(db, 145)

}
