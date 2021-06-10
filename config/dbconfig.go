package config

import "fmt"

var (
        Dbconfig map[string]interface{}
)

func init() {
        Dbconfig = make(map[string]interface{})
        Dbconfig["host"] = "localhost"
        Dbconfig["port"] = 5432
        Dbconfig["user"] = "postgres"
        Dbconfig["password"] = "kjunn2000"
        Dbconfig["dbname"] = "straper_db"
}

func NewConnString() string {
        return fmt.Sprintf("host=%s port=%d user=%s "+
                "password=%s dbname=%s sslmode=disable",
                Dbconfig["host"], Dbconfig["port"], Dbconfig["user"], Dbconfig["password"], Dbconfig["dbname"])
}