package auth

import (
	"net/http"
	"whatsapp-server/internal/database"
	"whatsapp-server/internal/hashing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	TIME       = 1
	MEMORY     = 64 * 1024
	THREADS    = 4
	KEY_LENGTH = 32
	SALT_SIZE  = 8
)

var argonHasher = hashing.NewArgon2idHash(TIME, SALT_SIZE, MEMORY, THREADS, KEY_LENGTH)

type UserCred struct {
	Email  string
	Hash   []byte
	salt   []byte
	UserID []byte
}

type UserCredErr struct {
	msg string
}

func (u *UserCredErr) Error() string {
	return u.msg
}

func fetchUserEmailAndPassword(email string, dbConn database.DBConnection) (*UserCred, error) {
	if email == "" {
		return nil, &UserCredErr{"empty email"}
	}

	var user UserCred
	err := dbConn.DB.QueryRow("SELECT email, hash, salt, user_id FROM user_credentials WHERE email = ?", email).Scan(
		&user.Email, &user.Hash, &user.salt, &user.UserID)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func HandleLogin(dbConn database.DBConnection) func(*gin.Context) {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")

		if email == "" || password == "" {
			//TODO: add better errors
			c.AbortWithStatusJSON(http.StatusBadRequest, "no email or password")
			return
		}

		user, err := fetchUserEmailAndPassword(email, dbConn)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		err = argonHasher.Compare(hashing.HashSalt{Hash: user.Hash, Salt: user.salt}, []byte(password))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		session := sessions.Default(c)
		//TODO: experiment and see if you can abuse this somehow
		session.Set(email, *user)

		c.Redirect(http.StatusPermanentRedirect, "/")
	}
}

func HandleRegister(c *gin.Context) {

}
