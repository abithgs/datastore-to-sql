package main

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/syndtr/goleveldb/leveldb/journal"
	"io/ioutil"
	"log"
	"os"
	datastore "./datastore"
	ds "./internal/datastore"
)

var Strict bool = true

// Load the backup into real model
// backupFilePath - the backup file path
// dst - the struct model that represents datastore entity and the model you want to load the data of this backup
// onPreload - callback that will be called before loading each entity
// onResult - callback that will be called with already loaded entity in the model
func Load(backupFilePath string, dst interface{}, onPreload func(dst interface{}), onResult func(dst interface{})) {
	f, err := os.Open(backupFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	journals := journal.NewReader(f, nil, Strict, true)
	for {
		j, err := journals.Next()
		if err != nil {
			// log.Fatal(err)
			break
		}
		b, err := ioutil.ReadAll(j)
		if err != nil {
			// log.Fatal(err)
			break
		}
		pb := &ds.EntityProto{}
		if err := proto.Unmarshal(b, pb); err != nil {
			log.Fatal(err)
			break
		}
		if onPreload != nil {
			onPreload(dst)
		}
		datastore.LoadEntity(dst, pb)
		onResult(dst)
	}
}