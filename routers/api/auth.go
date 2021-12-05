package api

import (
	"JRYY/db"
	"JRYY/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"time"
)


func AdminLogin(c *gin.Context) {
	collection, closeConn := db.GlobalDB.Admin()
	defer closeConn()
	name := c.PostForm("name")
	password := c.PostForm("password")
	var admin model.Admin
	err := collection.Find(bson.M{"name": name}).One(&admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
	}
	if password == admin.Password {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["type"] = 0
		claims["sid"] = -1
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintln(err),
			})
		}
		c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln("账号或密码错误"),
		})
	}
}

func UserLogin(c *gin.Context) {
	collection, closeConn := db.GlobalDB.Participant()
	defer closeConn()
	name := c.PostForm("name")
	sid := c.PostForm("sid")
	var participant model.Participant
	err := collection.Find(bson.M{"name": name, "sid": sid}).One(&participant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["type"] = 1
	claims["sid"] = sid
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"token": t,
		"_id": participant.Id.Hex(),
		"sid": participant.Sid,
	})
}
