package api

import (
	"JRYY/db"
	"JRYY/model"
	"JRYY/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//func PublishAdminNoti(c *gin.Context) {
//	err := db.GlobalRedis.Publish(c, "admin-noti", "这是一条信息").Err()
//	if err != nil {
//		c.String(http.StatusInternalServerError, "通知发布异常")
//		return
//	}
//}

//func SubscribeAdminNoti(c *gin.Context) {
//	//c := context.Background()
//	subscriber := db.GlobalRedis.Subscribe(c, "admin-notification")
//	for {
//		msg, err := subscriber.ReceiveMessage(c)
//		if err != nil {
//			c.String(http.StatusInternalServerError, "通知接受异常")
//			return
//		}
//		fmt.Printf("获取了通知：%v", msg)
//		c.String(http.StatusOK, msg.String())
//	}
//}

func PublishAllNoti(c *gin.Context) {
	collection, closeConn := db.GlobalDB.Notification()
	defer closeConn()
	var newNoti model.Notification
	newNoti.Id = bson.NewObjectId()
	newNoti.Type = 0
	newNoti.Alive = true
	newNoti.Time = time.Now().Add(time.Hour * 8)
	newNoti.Level, _ = strconv.Atoi(c.PostForm("level"))
	newNoti.Title = c.PostForm("title")
	newNoti.Content = c.PostForm("content")
	err := collection.Insert(&newNoti)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
	}
	c.String(http.StatusOK, "公告发送成功！")
}

func PublishSomeNoti(c *gin.Context) {
	collection, closeConn := db.GlobalDB.Notification()
	defer closeConn()
	// 存储一个notification到数据库中
	var newNoti model.Notification
	newNoti.Id = bson.NewObjectId()
	newNoti.Title = c.PostForm("title")
	newNoti.Content = c.PostForm("content")
	newNoti.Time = time.Now().Add(time.Hour * 8)
	newNoti.Type = 1
	newNoti.Alive = true
	newNoti.IdList = c.PostForm("idList")
	err := collection.Insert(&newNoti)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
	}
	// 遍历每个participant
	collection2, closeConn2 := db.GlobalDB.Participant()
	defer closeConn2()
	pidList := strings.Split(newNoti.IdList, ",")
	//var participants []model.Participant

	var epList []utils.EmailParam
	for _, id := range pidList {
		content := newNoti.Content
		var participant model.Participant
		err = collection2.FindId(bson.ObjectIdHex(id)).One(&participant)
		//participants = append(participants, participant)
		content = strings.Replace(content, "$name$", participant.Name, -1)
		content = strings.Replace(content, "$sid$", participant.Sid, -1)
		content = strings.Replace(content, "$pid$", participant.Pid, -1)
		content = strings.Replace(content, "$wechat$", participant.Wechat, -1)
		content = strings.Replace(content, "$email$", participant.Email, -1)
		content = strings.Replace(content, "$pairPid$", participant.PairPid, -1)
		var ep utils.EmailParam
		ep.EmailAddr = participant.Email
		ep.Title = newNoti.Title
		ep.Content = content
		epList = append(epList, ep)
		//if err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"error": fmt.Sprintln(err),
		//	})
		//}
	}
	go utils.SendEmailList(epList)
	c.String(http.StatusOK, "通知发送成功！")
}

// 获取admin发出的所有noti
func GetAdminNoti(c *gin.Context) {
	collection, closeConn := db.GlobalDB.Notification()
	defer closeConn()
	var results []model.Notification
	err := collection.Find(bson.M{"type": bson.M{"$lt": 2}}).Sort("-time").All(&results)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
	}
	c.JSON(http.StatusOK, results)
}

// 获取最后一个全员公告
func GetLastAdminAllNoti(c *gin.Context) {
	collection, closeConn := db.GlobalDB.Notification()
	defer closeConn()
	var result model.Notification
	err := collection.Find(bson.M{"type": 0}).Sort("-time").One(&result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
	}
	c.JSON(http.StatusOK, result)
}

//func ParticipantGetNoti(c *gin.Context)  {
//	collection, closeConn := db.GlobalDB.Notification()
//	defer closeConn()
//}
