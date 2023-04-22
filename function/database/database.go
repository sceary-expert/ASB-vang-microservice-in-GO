// Copyright (c) 2021 Amirhossein Movahedi (@qolzam)
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package database

import (
	"context"
	"fmt"

	"example.com/core/data/mongodb"
)

var Db interface{}

// Connect open database connection
func Connect(ctx context.Context) error {
	fmt.Println("at database.go at line 20")
	// coreConfig := config.AppConfig
	// fmt.Println("coreConfig got the value from config.AppConfig at line 22")
	// fmt.Println("coreConfig  db type : ", coreConfig.DBType)
	// if coreConfig.DBType == nil {
	// 	fmt.Println("DB TYPE is NULL")
	// 	return errors.New("coreConfig.DBType is null")
	// }
	// *coreConfig.DBType = config.DB_MONGO
	// switch *coreConfig.DBType {
	// case config.DB_MONGO:
	// fmt.Println("switched to case 1 at line 26")
	mongoClient, err := mongodb.NewMongoClient(ctx, "User-Rels", "mongodb+srv://asb:asbdesigngo@cluster0.ninbdmn.mongodb.net/?retryWrites=true&w=majority")
	fmt.Println("mongoClient got assigned at line 28")
	Db = mongoClient
	fmt.Println("db is assigned to mongo client at line 30")
	return err
	// }
	// fmt.Println("did not switch at line 33")
	// return fmt.Errorf("Please set valid database type in confing file!")
}
