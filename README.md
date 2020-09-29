# This is the readme file for the key value storage server.

# To build and run the server:

* go get github.com/syndtr/goleveldb/leveldb
* go run main.go

# To send a CSV file to the server:
# NOTE: The records of the csv files sent to the server are expected to have exactly 3 values.
* curl --data-binary @your_file.csv http://localhost:1321/promotions/upload/1

# To get the object by the ID:
* curl http://localhost:1321/promotions/1?ID=record_id

