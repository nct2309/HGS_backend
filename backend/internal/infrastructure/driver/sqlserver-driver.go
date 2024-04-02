package driver

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// SqlServerDB is a struct as Singleton to hold the connection of sql server on Azure Cloud
type SqlServerDB struct {
	DB *gorm.DB
}

var SqlServer *SqlServerDB = nil

func ConnectSqlServerDB() *gorm.DB {
	if SqlServer == nil {
		SqlServer = &SqlServerDB{}
		//first-azure-sql.database.windows.net
		//username: Tin
		//password: ClashOfClan123
		//the database is default
		dsn := "sqlserver://Tin:ClashOfClan123@first-azure-sql.database.windows.net?database=Multidisciplinary-Project"
		db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

		if err != nil {
			panic(err)
		}

		SqlServer.DB = db
	}
	return SqlServer.DB
}

// close function for the only instance of SqlServerDB
func CloseSqlServerDB() {
	db, err := SqlServer.DB.DB()
	if err != nil {
		panic(err)
	}
	db.Close()
}
