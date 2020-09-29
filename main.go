package main

import (
    "encoding/csv"
    "log"
    "io"
	"net/http"
    "strings"

    "github.com/syndtr/goleveldb/leveldb"
)

// The handler for CSV file POST requests
type CSVRequestHandler struct {
    db *leveldb.DB
}

func (this *CSVRequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    reader := csv.NewReader(req.Body)
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("Error while reading the CSV file: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if len(record) != 3 {
            log.Printf("Error while reading the CSV file: found %d values, while expecting 3 in each record.", len(record))
			http.Error(w, "Each record in the CSV file must contain 3 values.", http.StatusBadRequest)
            return
        }

        // Persist the record into the key value store
        err = this.db.Put([]byte(record[0]), []byte(record[1] + "," + record[2]), nil)
        if err != nil {
            log.Printf("Error persisting a record: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    }
}

// The handler for GET
type GetRequestHandler struct {
    db *leveldb.DB
}

func (this *GetRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("ID")
    if id == "" {
        log.Printf("Error: No ID parameter is specified.")
        http.Error(w, "Error: No ID parameter is specified.", http.StatusBadRequest)
        return
    }

    // Get the record from the key value store
    data, err := this.db.Get([]byte(id), nil)
    
    if err == leveldb.ErrNotFound {
        w.Write([]byte("Not found\n"))
        return
    }
    if err != nil {
        log.Printf("Error while getting a record from the DB: %s", err)
    }

    // Only 2 comma-separated values are expected here.
    values := strings.Split(string(data), ",")
    if len(values) != 2 {
        log.Printf("Error while parsing the values of the record: Expected 2 values, received %d", len(values))
        http.Error(w, "Internal server error.", http.StatusInternalServerError)
    }

    w.Write([]byte("{\"id\":\"" + id + "\", \"price\":\"" + values[0] + "\", \"expiration_date\":\"" + values[1] + "\"}\n"))
}

func main() {
    // Create or connect to the database
    db, err := leveldb.OpenFile("data/storage.db", nil)
    defer db.Close()
    if err != nil {
        log.Printf("Error opening the DB: %s", err)
        return
    }

	http.Handle("/promotions/upload/1", &CSVRequestHandler{db})
	http.Handle("/promotions/1", &GetRequestHandler{db})

	http.ListenAndServe(":1321", nil)
}
