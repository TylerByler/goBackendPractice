package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type SQLInfo struct {
	Address  string `json:"address"`
	DBname   string `json:"dbname"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type OrderList struct {
	Empty  bool
	Orders []Order
}

type Order struct {
	Invoice         float32
	SalespersonName string
	Date            string
	Empty           bool
	Engravings      []Engraving
}

type Engraving struct {
	Invoice              float32
	ProductNumber        string
	ProductDescription   string
	ProductColor         string
	DesignNumber         string
	Font                 string
	EngravingDescription string
}

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
		db, err := openDB()
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

func orderListPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		db, err := openDB()
		checkErr(err)

		rows, err := db.Query("SELECT * FROM orders ORDER BY `orders`.`creation_date` DESC LIMIT 50;")
		checkErr(err)

		var orderList OrderList
		orderList.Empty = true
		tempSlice := make([]Order, 0)

		for rows.Next() {
			var order Order
			err := rows.Scan(&order.Invoice, &order.SalespersonName, &order.Date)
			checkErr(err)
			orderList.Empty = false

			tempSlice = append(tempSlice, order)
		}

		orderList.Orders = tempSlice

		fmt.Println(orderList)
		fmt.Println("Type Of \"orderlist\":", reflect.TypeOf(orderList).Kind())
		fmt.Println("Type Of \"orderlist.Orders\":", reflect.TypeOf(orderList.Orders).Kind())

		t, _ := template.ParseFiles("orderList.html")
		t.Execute(w, orderList)
	}
}

func selectedOrderPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		params, _ := url.ParseQuery(r.URL.RawQuery)
		fmt.Println("Parameters: ", params)

		invoice_number := params.Get("invoice")
		fmt.Println("INVOICE_NUMBER: ", invoice_number)

		db, err := openDB()
		checkErr(err)

		rows, err := db.Query("SELECT * FROM discrete_engravings WHERE `invoice_number` =" + invoice_number + ";")
		checkErr(err)

		orderInfo, err := db.Query("SELECT * FROM orders WHERE `invoice_number` =" + invoice_number + ";")
		checkErr(err)

		var order Order
		if orderInfo != nil {
			orderInfo.Scan(&order.Invoice, &order.SalespersonName, &order.Date)
		}

		for rows.Next() {
			var engraving Engraving
			err := rows.Scan(
				&engraving.Invoice,
				&engraving.ProductNumber,
				&engraving.ProductDescription,
				&engraving.ProductColor,
				&engraving.DesignNumber,
				&engraving.Font,
				&engraving.EngravingDescription)
			checkErr(err)

			order.Engravings = append(order.Engravings, engraving)
		}

		fmt.Println(order)
		fmt.Println("Type Of \"orderlist\":", reflect.TypeOf(order).Kind())
		fmt.Println("Type Of \"orderlist.Orders\":", reflect.TypeOf(order.Engravings).Kind())

		t, _ := template.ParseFiles("selectedOrder.html")
		t.Execute(w, order)
	}
}

func openDB() (*sql.DB, error) {
	jsonFile, err := os.Open("sql.info.json")
	checkErr(err)

	jsonValue, err := io.ReadAll(jsonFile)
	checkErr(err)

	var sqlinfo SQLInfo
	json.Unmarshal(jsonValue, &sqlinfo)

	sqlLogin := sqlinfo.User + ":" + sqlinfo.Password + "@tcp(" + sqlinfo.Address + ")/" + sqlinfo.DBname + "?charset=utf8"
	fmt.Println("Sql Login info: ", sqlLogin)

	return sql.Open("mysql", sqlLogin)
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
	http.HandleFunc("/orders", orderListPage)
	http.HandleFunc("/selectedorder", selectedOrderPage)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
