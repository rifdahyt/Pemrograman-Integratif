package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"io"
	//"github.com/julienschmidt/httprouter"
	"strconv"
	//"os"
	//"strings"
	//"time"
	_ "github.com/go-sql-driver/mysql"
	//"bytes"
	"html/template"
)

var tpl *template.Template

type event struct {
	ID_event     int
	NamaEvent    string
	TanggalEvent string
	TempatEvent  string
	HostEvent    string
}

type myEvent struct {
	ID_event     int    `json:"ID_event"`
	NamaEvent    string `json:"NamaEvent"`
	TanggalEvent string `json:"TanggalEvent"`
	TempatEvent  string `json:"TempatEvent"`
	HostEvent    string `json:"HostEvent"`
}

func main() {
	port := 8181
	
	http.HandleFunc("/eventobi/", func(w http.ResponseWriter, r *http.Request) {
		
			http.ServeFile(w, r, "eventobi.html")
	})
	
	http.HandleFunc("/getbyname/", func(w http.ResponseWriter, r *http.Request) {
		
			http.ServeFile(w, r, "getNama.html")
	})
	
	http.HandleFunc("/getbydate/", func(w http.ResponseWriter, r *http.Request) {
		
			http.ServeFile(w, r, "getTanggal.html")
	})
	
	http.HandleFunc("/getbyplace/", func(w http.ResponseWriter, r *http.Request) {
		
			http.ServeFile(w, r, "getTempat.html")
	})

	http.HandleFunc("/insert/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "insert.html")
	})
	
	http.HandleFunc("/getbyhost/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "getHost.html")
	})
	
	http.HandleFunc("/submitEvent/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			InsertEvent(w, r)
			GetAllEvent(w,r)	
		}
	})

	http.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			s := r.URL.Path[len("/event/"):]
			if s != "" {
				if s == "today" {
					GetAllTodaysEvent(w, r)
				} else if s == "tomorrow" {
					GetAllTomorrowEvent(w, r)
				} else if s == "upcoming" {
					GetAllUpcomingEvent(w, r)
				} else {}
			} else {
				GetAllEvent(w,r)
			}

		default:
			http.Error(w, "Invalid Request method", 405)
		}

	})
	
	http.HandleFunc("/nama/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if r.URL.Query().Get("nama") != "" {
				nama := r.URL.Query().Get("nama")
				GetEvent(w,r,nama)
			}
				
		default:
			http.Error(w, "Invalid Request Method", 405)
		}
	})

	http.HandleFunc("/host/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if r.URL.Query().Get("host") != "" {
				host := r.URL.Query().Get("host")
				GetHostEvent(w,r,host)
			}
				
		default:
			http.Error(w, "Invalid Request Method", 405)
		}
	})

	http.HandleFunc("/tanggal/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if r.URL.Query().Get("tanggal") != "" {
				tanggal := r.URL.Query().Get("tanggal")
				GetEventDate(w,r,tanggal)
			}
				
		default:
			http.Error(w, "Invalid Request Method", 405)
		}
	})

	http.HandleFunc("/tempat/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if r.URL.Query().Get("tempat") != "" {
				tempat := r.URL.Query().Get("tempat")
				GetEventPlace(w,r,tempat)
			}
				
		default:
			http.Error(w, "Invalid Request Method", 405)
		}
	})

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

//getAllEvent
//munculin semua event beserta host (query 1)
func GetAllEvent(w http.ResponseWriter, r *http.Request) {
	
	db, err := sql.Open("mysql",
		"root:@tcp(127.0.0.1:3306)/eventobi_db")
	if err != nil {
		log.Fatal(err)
	}
	
	tpl = template.Must(template.ParseGlob("templates/*"))
	
	rows, err := db.Query(
       `SELECT NamaEvent, 
       TanggalEvent, 
       TempatEvent, 
       HostEvent 
       FROM event;`)
    if err != nil {
		log.Fatal(err)
	}

    Event := make([]event, 0)
    for rows.Next() {
      myevent := event{}
      rows.Scan(&myevent.NamaEvent, &myevent.TanggalEvent, &myevent.TempatEvent, &myevent.HostEvent)
      Event = append(Event, myevent)
    }
	
	//log.Println(Event)
	tpl.ExecuteTemplate(w, "getAll.html", Event)

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
	
	tpl = template.Must(template.ParseGlob("templates/*"))
	
	rows, err := db.Query(
       `SELECT NamaEvent, TanggalEvent, TempatEvent, HostEvent FROM event WHERE TanggalEvent = CURDATE();`)
    if err != nil {
		log.Fatal(err)
	}

    Event := make([]event, 0)
    for rows.Next() {
      myevent := event{}
      rows.Scan(&myevent.NamaEvent, &myevent.TanggalEvent, &myevent.TempatEvent, &myevent.HostEvent)
      Event = append(Event, myevent)
    }
	
	//log.Println(Event)
	tpl.ExecuteTemplate(w, "getAll.html", Event)

}

//getAllTomorrowEvent
//Menampilkan semua event besok (query 3)
func GetAllTomorrowEvent(w http.ResponseWriter, r *http.Request) {
	
	db, err := sql.Open("mysql",
		"root:@tcp(127.0.0.1:3306)/eventobi_db")
	if err != nil {
		log.Fatal(err)
	}
	
	tpl = template.Must(template.ParseGlob("templates/*"))
	
	rows, err := db.Query(
       `SELECT NamaEvent, TanggalEvent, TempatEvent, HostEvent FROM event WHERE TanggalEvent BETWEEN CURDATE() + INTERVAL 1 DAY AND CURDATE() + INTERVAL 1 DAY;`)
    if err != nil {
		log.Fatal(err)
	}

    Event := make([]event, 0)
    for rows.Next() {
      myevent := event{}
      rows.Scan(&myevent.NamaEvent, &myevent.TanggalEvent, &myevent.TempatEvent, &myevent.HostEvent)
      Event = append(Event, myevent)
    }
	
	//log.Println(Event)
	tpl.ExecuteTemplate(w, "getAll.html", Event)

}

//getAllUpcomingEvent
//Menampilkan upcoming event (query 4)
func GetAllUpcomingEvent(w http.ResponseWriter, r *http.Request) {
	
	db, err := sql.Open("mysql",
		"root:@tcp(127.0.0.1:3306)/eventobi_db")
	if err != nil {
		log.Fatal(err)
	}
	
	tpl = template.Must(template.ParseGlob("templates/*"))
	
	rows, err := db.Query(
       `SELECT NamaEvent, TanggalEvent, TempatEvent, HostEvent FROM event WHERE TanggalEvent BETWEEN CURDATE() + INTERVAL 1 DAY AND CURDATE() + INTERVAL 30 DAY;`)
    if err != nil {
		log.Fatal(err)
	}

    Event := make([]event, 0)
    for rows.Next() {
      myevent := event{}
      rows.Scan(&myevent.NamaEvent, &myevent.TanggalEvent, &myevent.TempatEvent, &myevent.HostEvent)
      Event = append(Event, myevent)
    }
	
	//log.Println(Event)
	tpl.ExecuteTemplate(w, "getAll.html", Event)

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
func InsertEvent(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == http.MethodPost {
		Event := event{}
		
		idevent, _ :=strconv.Atoi(r.FormValue("ID_event"))
		
		Event.ID_event = idevent
		Event.NamaEvent = r.FormValue("NamaEvent")
		Event.TanggalEvent = r.FormValue("TanggalEvent")
		Event.TempatEvent = r.FormValue("TempatEvent")
		Event.HostEvent = r.FormValue("HostEvent")
		
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/eventobi_db")
			if err != nil {
				log.Fatal(err)
			}
		
		_,err = db.Exec("INSERT INTO event (ID_event,NamaEvent,TanggalEvent,TempatEvent,HostEvent) VALUES (?,?,?,?,?)", Event.ID_event, Event.NamaEvent, Event.TanggalEvent, Event.TempatEvent, Event.HostEvent)
		if err != nil {
		log.Fatal(err)
		}
	}

}
