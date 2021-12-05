package api

import (
	"JRYY/db"
	"JRYY/model"
	"JRYY/utils"
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetAllParticipant(c *gin.Context) {
	collection, closeConn := db.GlobalDB.Participant()
	defer closeConn()
	var results []model.Participant
	err := collection.Find(nil).All(&results)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	c.JSON(http.StatusOK, results)
}

type ParticipantBindContact struct {
	model.Participant
	PairWechat string `bson:"pairWechat" json:"pairWechat"`
	PairEmail  string `bson:"pairEmail" json:"pairEmail"`
}

func GetParticipantBySid(c *gin.Context) {
	collection, closeConn := db.GlobalDB.Participant()
	defer closeConn()
	sid := c.Param("sid")
	var p model.Participant
	err := collection.Find(bson.M{"sid": sid}).One(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	if p.HasPair {
		var other model.Participant
		err = collection.Find(bson.M{"pid": p.PairPid}).One(&other)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintln(err),
			})
			return
		}
		var result = ParticipantBindContact{p, other.Wechat, other.Email}
		c.JSON(http.StatusOK, result)
	} else {
		var result = ParticipantBindContact{p, "", ""}
		c.JSON(http.StatusOK, result)
	}
}

func GetParticipant(c *gin.Context) {
	collection, closeConn := db.GlobalDB.Participant()
	defer closeConn()
	id := c.Param("id")
	var p model.Participant
	err := collection.FindId(bson.ObjectIdHex(id)).One(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	c.JSON(http.StatusOK, p)
}

// 初始化全部的参与者，最开始只调用一次
func InitParticipantByCsv(c *gin.Context) {
	path := utils.BASE_PATH + "file/" + c.Param("path")
	collection, closeConn := db.GlobalDB.Participant()
	_, err := collection.RemoveAll(nil)
	if err != nil {
		return
	}
	defer closeConn()
	data := utils.ReadCsvFile(path)
	for i, line := range data {
		if i != 0 {
			var newParticipant model.Participant
			newParticipant.Id = bson.NewObjectId()
			newParticipant.Pid = line[0]
			newParticipant.Sid = line[1]
			newParticipant.Name = line[2]
			newParticipant.Gender = line[3]
			newParticipant.Wechat = line[4]
			newParticipant.Email = line[5]
			newParticipant.NativePlace = line[6]
			newParticipant.Height, err = strconv.ParseFloat(line[7], 32)
			newParticipant.Weight, err = strconv.ParseFloat(line[8], 32)
			newParticipant.Birthday = line[9]
			newParticipant.School = line[10]
			if line[11] == "不接受" {
				newParticipant.AcceptDiffSchool = false
			} else {
				newParticipant.AcceptDiffSchool = true
			}
			newParticipant.Major = line[12]
			newParticipant.FuturePlace = line[13]
			newParticipant.HopeOther = strings.Join([]string{line[14], line[15], line[16], line[17], line[18]}, "$")
			newParticipant.Declaration = line[19]
			newParticipant.ModifyTimes = 3
			newParticipant.HasPair = false
			newParticipant.PairPid = ""
			err = collection.Insert(&newParticipant)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintln(err),
				})
				return
			}
		}
	}
	c.String(http.StatusOK, "导入成功！")
}

// 通过csv文件增加参与者
func AddParticipantByCsv(c *gin.Context) {
	path := utils.BASE_PATH + "file/" + c.Param("path")
	collection, closeConn := db.GlobalDB.Participant()
	defer closeConn()
	data := utils.ReadCsvFile(path)
	for i, line := range data {
		if i != 0 {
			var newParticipant model.Participant
			newParticipant.Id = bson.NewObjectId()
			newParticipant.Pid = line[0]
			newParticipant.Sid = line[1]
			newParticipant.Name = line[2]
			newParticipant.Gender = line[3]
			newParticipant.Wechat = line[4]
			newParticipant.Email = line[5]
			newParticipant.NativePlace = line[6]
			newParticipant.Height, _ = strconv.ParseFloat(line[7], 32)
			newParticipant.Weight, _ = strconv.ParseFloat(line[8], 32)
			newParticipant.Birthday = line[9]
			newParticipant.School = line[10]
			if line[11] == "不接受" {
				newParticipant.AcceptDiffSchool = false
			} else {
				newParticipant.AcceptDiffSchool = true
			}
			newParticipant.Major = line[12]
			newParticipant.FuturePlace = line[13]
			newParticipant.HopeOther = strings.Join([]string{line[14], line[15], line[16], line[17], line[18]}, "$")
			newParticipant.Declaration = line[19]
			newParticipant.ModifyTimes = 3
			newParticipant.HasPair = false
			newParticipant.PairPid = ""
			err := collection.Insert(&newParticipant)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintln(err),
				})
				return
			}
		}
	}
	c.String(http.StatusOK, "导入成功！")
}

// 通过csv文件增加pair
func AddParticipantPairByCsv(c *gin.Context) {
	path := utils.BASE_PATH + "file/" + c.Param("path")
	collection, closeConn := db.GlobalDB.Participant()
	defer closeConn()
	data := utils.ReadCsvFile(path)
	for i, line := range data {
		if i != 0 {
			var p1, p2 model.Participant
			err := collection.Find(bson.M{"pid": line[0]}).One(&p1)
			err = collection.Find(bson.M{"pid": line[1]}).One(&p2)
			p1.HasPair = true
			p1.PairPid = p2.Pid
			p2.HasPair = true
			p2.PairPid = p1.Pid
			err = collection.Update(bson.M{"sid": p1.Sid}, p1)
			err = collection.Update(bson.M{"sid": p2.Sid}, p2)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintln(err),
				})
				return
			}
		}
	}
	c.String(http.StatusOK, "添加成功！")
}

// 删除一对配对
func RemoveParticipantPair(c *gin.Context) {
	pid1 := c.PostForm("pid1")
	pid2 := c.PostForm("pid2")
	if err := deletePidPair(pid1); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	if err := deletePidPair(pid2); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	c.String(http.StatusOK, "删除成功！")
}

// 添加一对配对
func AddParticipantPair(c *gin.Context) {
	collection, closeConn := db.GlobalDB.Participant()
	defer closeConn()
	var p1, p2 model.Participant
	pid1 := c.PostForm("pid1")
	pid2 := c.PostForm("pid2")
	err := collection.Find(bson.M{"pid": pid1}).One(&p1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	err = collection.Find(bson.M{"pid": pid2}).One(&p2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	// 如果已经有配对了，就删除
	if p1.PairPid != "" {
		err = deletePidPair(p1.PairPid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintln(err),
			})
			return
		}
	}
	if p2.PairPid != "" {
		err = deletePidPair(p2.PairPid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintln(err),
			})
			return
		}
	}
	p1.HasPair = true
	p1.PairPid = pid2
	p2.HasPair = true
	p2.PairPid = pid1
	err = collection.Update(bson.M{"sid": p1.Sid}, p1)
	err = collection.Update(bson.M{"sid": p2.Sid}, p2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	c.String(http.StatusOK, "添加成功！")
}

// 删除pid的配对信息
func deletePidPair(pid string) error {
	var p model.Participant
	collection, closeConn := db.GlobalDB.Participant()
	defer closeConn()
	err := collection.Find(bson.M{"pid": pid}).One(&p)
	p.HasPair = false
	p.PairPid = ""
	err = collection.Update(bson.M{"sid": p.Sid}, p)
	return err
}

func ModifyParticipantBySid(c *gin.Context) {
	collection, closeConn := db.GlobalDB.Participant()
	defer closeConn()
	var updatedParticipant model.Participant
	err := c.BindJSON(&updatedParticipant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	sid := c.Param("sid")
	updatedParticipant.ModifyTimes -= 1
	err = collection.Update(bson.M{"sid": sid}, updatedParticipant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	c.JSON(http.StatusOK, updatedParticipant)
}

func GetReportDownload(c *gin.Context) {
	var reportStream [][]string
	var pList []model.Participant
	collection, closeConn := db.GlobalDB.Participant()
	defer closeConn()
	err := collection.Find(nil).All(&pList)
	if err != nil {
		log.Fatal(err)
	}
	var fileHead []string
	fileHead = append(fileHead, "项目编号", "学号", "姓名", "性别", "微信", "邮箱", "籍贯", "身高", "体重", "生日", "校区", "跨校区", "专业", "未来发展", "理想1", "理想2", "理想3", "理想4", "理想5", "宣言")
	reportStream = append(reportStream, fileHead)
	for _, participant := range pList {
		var line []string
		line = append(line, participant.Pid, participant.Name, participant.Sid, participant.Wechat, participant.Email, participant.NativePlace, fmt.Sprintf("%.1f", participant.Height), fmt.Sprintf("%.1f", participant.Weight), participant.Birthday, participant.School)
		if participant.AcceptDiffSchool {
			line = append(line, "接受")
		} else {
			line = append(line, "不接受")
		}
		line = append(line, participant.Major, participant.FuturePlace)
		hopeList := strings.Split(participant.HopeOther, "$")
		line = append(line, hopeList...)
		line = append(line, participant.Declaration)
		reportStream = append(reportStream, line)
	}
	b := &bytes.Buffer{}
	b.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(b)
	//c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+"名单.csv")
	//c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	for _, record := range reportStream {
		if err = w.Write(record); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintln(err),
			})
			return
		}
		w.Flush()
	}
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}

func GetPublicInfo(c *gin.Context) {
	collection, closeConn := db.GlobalDB.Participant()
	defer closeConn()
	var all []model.Participant
	var publicInfoTable []model.ParticipantPublic
	gender := c.Param("gender")
	g := "男"
	if gender == "girl" {
		g = "女"
	}
	err := collection.Find(bson.M{"gender": g}).All(&all)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	for _, p := range all {
		if p.Pid != "b000" && p.Pid != "g000" {
			var publicInfo = model.ParticipantPublic{
				Pid:              p.Pid,
				Gender:           p.Gender,
				NativePlace:      p.NativePlace,
				Height:           p.Height,
				Weight:           p.Weight,
				Birthday:         p.Birthday,
				School:           p.School,
				AcceptDiffSchool: p.AcceptDiffSchool,
				Major:            p.Major,
				FuturePlace:      p.FuturePlace,
				HopeOther:        p.HopeOther,
				Declaration:      p.Declaration,
				HasPair:          p.HasPair,
				Resume:           p.Resume,
			}
			publicInfoTable = append(publicInfoTable, publicInfo)
		}
	}
	c.JSON(http.StatusOK, publicInfoTable)
}

type ListWithPageNum struct {
	PublicInfoList []model.ParticipantPublic `bson:"publicInfoList" json:"publicInfoList"`
	PageCount      int                       `bson:"pageCount" json:"pageCount"`
}

func GetPublicInfoPage(c *gin.Context) {
	collection, closeConn := db.GlobalDB.Participant()
	defer closeConn()
	var someParticipant []model.Participant
	var publicInfoTable []model.ParticipantPublic
	gender := c.Param("gender")
	g := "男"
	if gender == "girl" {
		g = "女"
	}
	pageNum, _ := strconv.Atoi(c.Param("page"))
	pageSize := 20
	allCount, _ := collection.Find(bson.M{"gender": g}).Count()
	pageCount := int(allCount/pageSize)
	err := collection.Find(bson.M{"gender": g}).Sort("pid").Skip(pageSize * (pageNum - 1)).Limit(pageSize).All(&someParticipant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	for _, p := range someParticipant {
		if p.Pid != "b000" && p.Pid != "g000" {
			var publicInfo = model.ParticipantPublic{
				Pid:              p.Pid,
				Gender:           p.Gender,
				NativePlace:      p.NativePlace,
				Height:           p.Height,
				Weight:           p.Weight,
				Birthday:         p.Birthday,
				School:           p.School,
				AcceptDiffSchool: p.AcceptDiffSchool,
				Major:            p.Major,
				FuturePlace:      p.FuturePlace,
				HopeOther:        p.HopeOther,
				Declaration:      p.Declaration,
				HasPair:          p.HasPair,
				Resume:           p.Resume,
			}
			publicInfoTable = append(publicInfoTable, publicInfo)
		}
	}
	c.JSON(http.StatusOK, ListWithPageNum{PublicInfoList: publicInfoTable, PageCount: pageCount})
}
