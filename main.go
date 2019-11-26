package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var tableNames []string
var db *sql.DB
var err error
var exiting bool = false
var tcpClosed = false
var udpClosed = false

func main() {

	InitDb("sqlite3", "./suds.db")
	defer db.Close()

	// setup the udp listener
	ServerConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 10001, Zone: ""})
	defer ServerConn.Close()

	// spawn a thread for the listener
	go HandleUDPData(ServerConn)

	http.HandleFunc("/dump/", handleConnection)

	// setup the tcp listener
	http.ListenAndServe(":8080", nil)

	// intercept the KILL signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		fmt.Println("Caught CTRL-C")
		exiting = true

		ServerConn.Close()

	}()

	// loop infinitely while not exiting. Probably a better way of doing this.
	for !exiting {
		time.Sleep(1)
	}

	// When we get here, we are exiting which will signal the UDP and TCP sockets to close
	fmt.Print("Waiting to close sockets.")
	for i := 0; i < 10 && (tcpClosed == false && udpClosed == false); i++ {
		fmt.Println("i=", i, "; tcpClosed: ", tcpClosed, "; udpClosed=", udpClosed)
		time.Sleep(10 * time.Second)
	}

	/*
		var sqlStmt string

		sqlStmt = `SELECT name FROM sqlite_master WHERE type='table'`

		sqlStmt = `
		create table foo (id integer not null primary key, name text);
		delete from foo;
		`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		for i := 0; i < 100; i++ {
			_, err = stmt.Exec(i, fmt.Sprintf("???????%03d", i))
			if err != nil {
				log.Fatal(err)
			}
		}
		tx.Commit()

		rows, err := db.Query("select id, name from foo")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var name string
			err = rows.Scan(&id, &name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id, name)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		stmt, err = db.Prepare("select name from foo where id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		var name string
		err = stmt.QueryRow("3").Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(name)

		_, err = db.Exec("delete from foo")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
		if err != nil {
			log.Fatal(err)
		}

		rows, err = db.Query("select id, name from foo")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var name string
			err = rows.Scan(&id, &name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id, name)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	*/

}
