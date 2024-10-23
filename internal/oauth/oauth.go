package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"whatsapp-server/internal/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type IPAddr string
type State string

var tempDB map[IPAddr]State
var conf *oauth2.Config = nil

func generateRandString() State {
	b := make([]byte, 32)
	rand.Read(b)
	return State(base64.RawStdEncoding.EncodeToString(b))
}

func HandleGoogleAuthLink(c *gin.Context) {
	if conf == nil {
		conf = &oauth2.Config{
			ClientID:     os.Getenv("CLIENT_ID_ENV"),
			ClientSecret: os.Getenv("CLIENT_SECRET_ENV"),
			RedirectURL:  "http://localhost:8080/api/auth/callback/google",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
		}
	}
	if tempDB == nil {
		tempDB = make(map[IPAddr]State)
	}

	clientIP := IPAddr(c.ClientIP())
	tempDB[clientIP] = generateRandString()

	c.Redirect(http.StatusTemporaryRedirect, conf.AuthCodeURL(string(tempDB[clientIP])))
}

func HandleGoogleCallBack(dbConn database.DBConnection) func(*gin.Context) {
	return func(c *gin.Context) {
		if conf == nil || tempDB == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "{\"status\":\"conf or tempDB is null\"}")
			return
		}
		token := tempDB[IPAddr(c.ClientIP())]
		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, "{\"status\":\"token is null\"}")
			return
		}

		requestToken := State(c.Query("state"))
		if requestToken == "" || requestToken != token {
			c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("{\"status\":\"bad token %s \"}", token))
			return
		}

		ctx := context.Background()

		tok, err := conf.Exchange(ctx, c.Query("code"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("{\"status\":\"we fucked up %s \"}", err.Error()))
			return
		}

		//store the token into my other tempdb and use that for sessions
		// log.Println(*tok)

		client := conf.Client(ctx, tok)
		response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo?alt=json")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("{\"status\":\"we fucked up \", %s}", err.Error()))
			return
		}
		defer response.Body.Close()

		//TODO: HANDLE ERROR
		data, _ := io.ReadAll(response.Body)

		data2 := base64.RawStdEncoding.EncodeToString(data)

		c.JSON(http.StatusOK, data2)

		//remove this
		err = dbConn.DB.Ping()
		if err != nil {
			log.Fatal(err)
		}
	}

}
