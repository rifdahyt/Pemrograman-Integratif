package main

import (
		"database/sql"
		"encoding/json"
		"fmt"
		"log"
		"net/http"
		//"github.com/julienschmidt/httprouter"
		//"strconv"
		//"strings"
		//"time"
		_"github.com/go-sql-driver/mysql"
)

type event struct {
		ID_event 		int
		NamaEvent		string
		TanggalEvent	string
		TempatEvent		string
		HostEvent		string
}

type myEvent struct {
		ID_event		int		`json: "ID_event, omitempty"`
		NamaEvent		string	`json: "NamaEvent, omitempty"`
		TanggalEvent	string	`json: "TanggalEvent, omitempty"`
		TempatEvent		string	`json: "TempatEvent, omitempty"`
		HostEvent		string	`json: "HostEvent, omitempty"`
}


func main() {
		port := 8181
		
		http.HandleFunc("/insert", func(w http.ResponseWriter, r* http.Request) {
			http.ServeFile(w,r,"insert.html")
		})
		
		http.HandleFunc("/event/", func(w http.ResponseWriter, r* http.Request) {
			switch r.Method{
				case "GET":
					s := r.URL.Path[len("/event/"):]
					if s !="" {
						if s == "today" {
							GetAllTodaysEvent(w,r)
						} else if s == "tomorrow" {
							GetAllTomorrowEvent(w,r)
						} else if s == "upcoming" {
							GetAllUpcomingEvent(w,r)
						} else {
							GetEvent(w,r,s)
						}	
					}else{
						GetAllEvent(w,r)
					}
					
				case "POST":
					InsertEvent(w,r)
					
				default:
					http.Error(w, "Invalid Request method", 405)
			}
		
		})
		
		http.HandleFunc("/host/", func(w http.ResponseWriter, r* http.Request) {
			switch r.Method{
				case "GET":
					s := r.URL.Path[len("/host/"):]
					if s !="" {
							GetHostEvent(w,r,s)
					}else{
						GetAllEvent(w,r)
					}
				default:
					http.Error(w, "Invalid Request Method", 405)
			}
		})
		
		http.HandleFunc("/tanggal/", func(w http.ResponseWriter, r* http.Request) {
			switch r.Method{
				case "GET":
					s := r.URL.Path[len("/tanggal/"):]
					if s !="" {
							GetEventDate(w,r,s)
					}
				default:
					http.Error(w, "Invalid Request Method", 405)
			}
		})
		
		http.HandleFunc("/tempat/", func(w http.ResponseWriter, r* http.Request) {
			switch r.Method{
				case "GET":
					s := r.URL.Path[len("/tempat/"):]
					if s !="" {
							GetEventPlace(w,r,s)
					}
				default:
					http.Error(w, "Invalid Request Method", 405)
			}
		})
		
		log.Printf("Server starting on port %v\n", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port),nil)) 
}

//getAllEvent
//munculin semua event beserta host (query 1)
func GetAllEvent(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("mysql", 
			"root:@tcp(127.0.0.1:3306)/eventobi_db")
		if err != nil {
			log.Fatal(err)
		}
		
		defer db.Close()
		
		Event := event{}
		
		rows, err := db.Query("SELECT NamaEvent, TanggalEvent, TempatEvent, HostEvent FROM event")
		if err != nil {
			log.Fatal(err)
		}
		
		defer rows.Close()
		
		for rows.Next() {
				err := rows.Scan(&Event.NamaEvent, &Event.TanggalEvent, &Event.TempatEvent, &Event.HostEvent)
				if err != nil {
					log.Fatal(err)
				}
			json.NewEncoder(w).Encode(&Event)
		}						
		err = rows.Err();
		
							
}
							

//getEvent berdasarkan nama
//search event berdasarkan nama event (query 5)
func GetEvent(w http.ResponseWriter, r *http.Request, s string) {
		
		db, err := sql.Open("mysql",
				"root:@tcp(127.0.0.1:3306)/eventobi_db")
		if err != nil {
			log.Fatal(err)
		}
		
		defer db.Close()
		
		Event := event{}
		
		rows, err := db.Query("SELECT NamaEvent, TanggalEvent, TempatEvent, HostEvent FROM event WHERE NamaEvent LIKE ?", "%"+s+"%")
		if err != nil {
			log.Fatal(err)
		}
		
		defer db.Close()
		for rows.Next() {
			err := rows.Scan(&Event.NamaEvent, &Event.TanggalEvent, &Event.TempatEvent, &Event.HostEvent)
			if err != nil {
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode(&Event)	
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
}

//getAllTodaysEvent 
//menampilkan semua event hari ini (query 2)
func GetAllTodaysEvent(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("mysql", 
			"root:@tcp(127.0.0.1:3306)/eventobi_db")
		if err != nil {
			log.Fatal(err)
		}
		
		defer db.Close()
		
		Event := event{}
		
		rows, err := db.Query("SELECT NamaEvent, TanggalEvent, TempatEvent, HostEvent FROM event WHERE TanggalEvent = CURDATE()")
		if err != nil {
			log.Fatal(err)
		}
		
		defer rows.Close()
		
		for rows.Next() {
				err := rows.Scan(&Event.NamaEvent, &Event.TanggalEvent, &Event.TempatEvent, &Event.HostEvent)
				if err != nil {
					log.Fatal(err)
				}
			json.NewEncoder(w).Encode(&Event)
		}					
		err = rows.Err();
								
}

//getAllTomorrowEvent
//Menampilkan semua event besok (query 3)
func GetAllTomorrowEvent(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("mysql", 
			"root:@tcp(127.0.0.1:3306)/eventobi_db")
		if err != nil {
			log.Fatal(err)
		}
		
		defer db.Close()
		
		Event := event{}
		
		rows, err := db.Query("SELECT NamaEvent, TanggalEvent, TempatEvent, HostEvent FROM event WHERE TanggalEvent BETWEEN CURDATE() + INTERVAL 1 DAY AND CURDATE() + INTERVAL 1 DAY")
		if err != nil {
			log.Fatal(err)
		}
		
		defer rows.Close()
		
		for rows.Next() {
				err := rows.Scan(&Event.NamaEvent, &Event.TanggalEvent, &Event.TempatEvent, &Event.HostEvent)
				if err != nil {
					log.Fatal(err)
				}
			json.NewEncoder(w).Encode(&Event)	
		}						
		err = rows.Err();
								
}

//getAllUpcomingEvent
//Menampilkan upcoming event (query 4)
func GetAllUpcomingEvent(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("mysql", 
			"root:@tcp(127.0.0.1:3306)/eventobi_db")
		if err != nil {
			log.Fatal(err)
		}
		
		defer db.Close()
		
		Event := event{}
		
		rows, err := db.Query("SELECT NamaEvent, TanggalEvent, TempatEvent, HostEvent FROM event WHERE TanggalEvent BETWEEN CURDATE() + INTERVAL 1 DAY AND CURDATE() + INTERVAL 90 DAY")
		if err != nil {
			log.Fatal(err)
		}
		
		defer rows.Close()
		
		for rows.Next() {
				err := rows.Scan(&Event.NamaEvent, &Event.TanggalEvent, &Event.TempatEvent, &Event.HostEvent)
				if err != nil {
					log.Fatal(err)
				}
			json.NewEncoder(w).Encode(&Event)	
		}						
		err = rows.Err();
								
}

//getHostEvent
//menampilkan event berdasarkan nama host (query 6)
func GetHostEvent(w http.ResponseWriter, r *http.Request, s string) {
		
		db, err := sql.Open("mysql",
				"root:@tcp(127.0.0.1:3306)/eventobi_db")
		if err != nil {
			log.Fatal(err)
		}
		
		defer db.Close()
		
		Event := event{}
		
		rows, err := db.Query("SELECT NamaEvent, TanggalEvent, TempatEvent, HostEvent FROM event WHERE HostEvent LIKE ?", "%"+s+"%")
		if err != nil {
			log.Fatal(err)
		}
		
		defer db.Close()
		for rows.Next() {
			err := rows.Scan(&Event.NamaEvent, &Event.TanggalEvent, &Event.TempatEvent, &Event.HostEvent)
			if err != nil {
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode(&Event)	
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
}

//getEventDate
//menampilkan event berdasarkan tanggal (query 7)
func GetEventDate(w http.ResponseWriter, r *http.Request, s string) {
		
		db, err := sql.Open("mysql",
				"root:@tcp(127.0.0.1:3306)/eventobi_db")
		if err != nil {
			log.Fatal(err)
		}
		
		defer db.Close()
		
		Event := event{}
		
		rows, err := db.Query("SELECT NamaEvent, TanggalEvent, TempatEvent, HostEvent FROM event WHERE TanggalEvent LIKE ?", "%"+s+"%")
		if err != nil {
			log.Fatal(err)
		}
		
		defer db.Close()
		for rows.Next() {
			err := rows.Scan(&Event.NamaEvent, &Event.TanggalEvent, &Event.TempatEvent, &Event.HostEvent)
			if err != nil {
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode(&Event)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
}

//getEventPlace
//menampilkan event berdasarkan tempat (query 8)
func GetEventPlace(w http.ResponseWriter, r *http.Request, s string) {
		
		db, err := sql.Open("mysql",
				"root:@tcp(127.0.0.1:3306)/eventobi_db")
		if err != nil {
			log.Fatal(err)
		}
		
		defer db.Close()
		
		Event := event{}
		
		rows, err := db.Query("SELECT NamaEvent, TanggalEvent, TempatEvent, HostEvent FROM event WHERE TempatEvent LIKE ?", "%"+s+"%")
		if err != nil {
			log.Fatal(err)
		}
		
		defer db.Close()
		for rows.Next() {
			err := rows.Scan(&Event.NamaEvent, &Event.TanggalEvent, &Event.TempatEvent, &Event.HostEvent)
			if err != nil {
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode(&Event)	
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
}

//InsertEvent
//Menambahkan event
func InsertEvent (w http.ResponseWriter, r *http.Request) {
	var Event myEvent
	dec := json.NewDecoder(r.Body)
	err:=dec.Decode(&Event)
	if err != nil{
		log.Fatal(err)
	}
	defer r.Body.Close()

	
	db, err := sql.Open("mysql",
				"root:@tcp(127.0.0.1:3306)/eventobi_db")
	if err != nil {
		log.Fatal(err)
	}
	
	stmt, err := db.Prepare("INSERT INTO event (ID_event, NamaEvent, TanggalEvent, TempatEvent, HostEvent) VALUES (?,?,?,?,?)")
	if err != nil{
		log.Fatal(err)
	}
	
	_, err = stmt.Exec(Event.ID_event, Event.NamaEvent, Event.TanggalEvent, Event.TempatEvent, Event.HostEvent)
	
}

