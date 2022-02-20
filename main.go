package main

import (
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
	_, err := db.Query("INSERT INTO posts(post, user_id) VALUES (?, ?)", p.post, p.user_id)

	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
		return err
	}
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

func unblockFriend(db *sql.DB, id int, receiver int) error {
	_, err := db.Query("DELETE FROM blocks WHERE (sender = ? AND receiver = ?)", id, receiver)

	if err != nil {
		panic(err.Error())
	}
	return nil
}

func searchForUsers(db *sql.DB, id int, first_name string, last_name string) {

	res, err := db.Query("SELECT id, first_name, last_name FROM users WHERE id != ? AND (first_name = ? OR last_name = ?)", id, first_name, last_name)

	if err != nil {
		panic(err.Error())
	}

	for res.Next() {
		var users users
		err = res.Scan(&users.id, &users.first_name, &users.last_name)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(users.id, users.first_name, users.last_name)
	}
}

func removeFriend(db *sql.DB, id int, user_id int) error {
	_, err := db.Query("DELETE FROM friendships WHERE (sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)", id, user_id, user_id, id)

	if err != nil {
		panic(err.Error())
	}
	return nil
}

func rejectRequest(db *sql.DB, id int, user_id int) error {
	_, err := db.Query("DELETE FROM friendships WHERE sender = ? AND receiver = ?", id)

	if err != nil {
		panic(err.Error())
	}
	return nil
}

func insertPicture(db *sql.DB, u users, id int) error {
	_, err := db.Query("UPDATE users  SET picture=? WHERE users.id = ?", u.picture, id)

	if err != nil {
		log.Printf("Error %s when inserting picture into users table", err)
		return err
	}
	return nil
}

func getMyPosts(db *sql.DB, id int) {

	res, err := db.Query("SELECT posts.post, posts.timestamp FROM posts WHERE posts.user_id = ?", id)

	if err != nil {
		panic(err.Error())
	}

	for res.Next() {
		var posts posts
		err = res.Scan(&posts.post, &posts.timestamp)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(posts.post, posts.timestamp)
	}
}

func getPostLikes(db *sql.DB, post_id int) {
	res, err := db.Query("SELECT COUNT(like_id) FROM likes WHERE post_id = ?", post_id)

	if err != nil {
		panic(err.Error())
	}

	for res.Next() {
		var count int
		err = res.Scan(&count)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Number of like are ", count)
	}
}

func getFriends(db *sql.DB, id int) {

	res, err := db.Query("SELECT id, first_name, last_name, picture FROM users INNER JOIN friendships ON users.id = friendships.sender OR users.id = friendships.receiver LEFT JOIN blocks ON  users.id = blocks.receiver OR users.id = blocks.sender WHERE (friendships.sender = ? OR friendships.receiver = ?) AND friendships.accepted = 1 AND id != ? AND id NOT IN (SELECT blocks.sender FROM blocks WHERE blocks.receiver = ?) AND id NOT IN (SELECT blocks.receiver FROM blocks WHERE blocks.sender = ?)", id, id, id, id, id)

	if err != nil {
		panic(err.Error())
	}

	for res.Next() {
		var users users

		err = res.Scan(&users.id, &users.first_name, &users.last_name, &users.picture)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(users.id, users.first_name, users.last_name, users.picture)
	}
}

func getFriendRequests(db *sql.DB, id int) {

	res, err := db.Query("SELECT id, first_name, last_name, picture FROM users INNER JOIN friendships ON users.id = friendships.sender OR users.id = friendships.receiver WHERE friendships.receiver = ? AND friendships.accepted = 0 AND id != ?", id, id)

	if err != nil {
		panic(err.Error())
	}

	for res.Next() {
		var users users

		err = res.Scan(&users.id, &users.first_name, &users.last_name, &users.picture)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(users.id, users.first_name, users.last_name, users.picture)
	}
}

func getBlockedUsers(db *sql.DB, id int) {

	res, err := db.Query("SELECT id, first_name, last_name, picture FROM users INNER JOIN blocks ON users.id = blocks.receiver WHERE blocks.sender = ?", id)

	if err != nil {
		panic(err.Error())
	}

	for res.Next() {
		var users users

		err = res.Scan(&users.id, &users.first_name, &users.last_name, &users.picture)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(users.id, users.first_name, users.last_name, users.picture)
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
	//getUserData(db, 145)
	//p := posts{
	//	post:    "ZEINAB",
	//	user_id: 10,
	//}

	// perform a db.Query insert
	//err = insertPost(db, p)

	// if there is an error inserting, handle it
	//if err != nil {
	//	panic(err.Error())
	//}

	//unblockFriend(db, 143, 148)
	//searchForUsers(db, 144, "Mohammad", "Badreddine")
	//removeFriend(db, 143, 145)
	//getPostLikes(db, 253)
	//getFriends(db, 143)
	//getFriendRequests(db, 147)
	getBlockedUsers(db, 143)
}
