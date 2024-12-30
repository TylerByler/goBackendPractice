package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func homepage(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("method:", r.Method)
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func entryPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "POST" {
		r.ParseForm()

		// <----SQL---->
		db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/millertech_engraving?charset=utf8")
		checkErr(err)

		// Create order in order table
		stmt, err := db.Prepare("INSERT INTO orders (invoice_number, salesperson_name) VALUES (?, ?);")
		checkErr(err)

		invoice_number := r.Form["invoice_number"][0]
		salesperson_name := r.Form["salesperson_name"][0]

		res, err := stmt.Exec(invoice_number, salesperson_name)
		stmt.Close()
		stmt = nil
		checkErr(err)

		affect, err := res.RowsAffected()
		checkErr(err)

		fmt.Println("Rows affected by insert into order: ", affect)

		// Enter each engraving into the Discrete Orders array
		stmt, err = db.Prepare("INSERT INTO discrete_engravings (invoice_number, product_number, product_desc, color, design_number, font, engraving_desc) VALUES (?, ?, ?, ?, ?, ?, ?);")
		checkErr(err)
		fmt.Println("Invoice Number: ", r.Form["invoice_number"], reflect.TypeOf(r.Form["invoice_number"]))
		numEntries, _ := strconv.Atoi(r.FormValue("num_entries"))
		for i := 0; i < numEntries; i++ {
			product_number := r.Form["product_number[]"][i]
			product_desc := r.Form["product_desc[]"][i]
			color := r.Form["color[]"][i]
			design_number := r.Form["design_number[]"][i]
			font := r.Form["font[]"][i]
			engraving_desc := r.Form["engraving_desc[]"][i]

			_, err = stmt.Exec(invoice_number, product_number, product_desc, color, design_number, font, engraving_desc)
			checkErr(err)
			fmt.Println("Product Number: #", i, " ", r.Form["product_number[]"][i], " ", reflect.TypeOf(r.Form["product_number[]"][i]))
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		t, _ := template.ParseFiles("entry.html")
		t.Execute(w, nil)
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

/* func login(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		// create a unique token for every form to prevent duplicate submissions
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("login.html")
		t.Execute(w, token)
	} else {
		r.ParseForm()
		token := r.Form.Get("token")
		if token != "" {
			// validate token
		} else {
			// give error if no token
		}
		fmt.Println("username length", len(r.Form["username"][0]))
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username")))
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username")))
	}
} */

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.HandleFunc("/", homepage)
	/* http.HandleFunc("/login", login) */
	http.HandleFunc("/entry", entryPage)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
