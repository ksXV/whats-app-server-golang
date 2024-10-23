package main

import (
	"log"
	"whatsapp-server/internal/database"
	"whatsapp-server/internal/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memcached"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter(dbConn database.DBConnection, mc database.MemCacheDB) *gin.Engine {
	r := gin.Default()

	//TODO: replace these with actually something
	store := memcached.NewStore(mc.DB, "", []byte("mysecret"))

	//TODO: add csrf protection either handrolled or library
	r.Use(sessions.Sessions("foobar", store))

	routes.GroupRoutes(r, dbConn)

	return r
}

func main() {
	err := godotenv.Load()

	db := database.NewDBConnection()
	mc := database.NewMemCacheDB()

	if err != nil {
		log.Fatalf("ERROR LOADING .env %s", err.Error())
	}

	r := setupRouter(db, mc)

	r.Run(":8080")
}
