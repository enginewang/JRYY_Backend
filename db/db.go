package db

import (
	"github.com/globalsign/mgo"
)

//
//var GlobalRedis *redis.Client
//
//func InitRedisServer() error {
//	rdb := redis.NewClient(&redis.Options{
//		Addr:     utils.REDIS_URL,
//		Password: "",
//		DB:       0,
//	})
//	GlobalRedis = rdb
//	return nil
//}

type Database struct {
	session  *mgo.Session
	database string
}

var GlobalDB *Database

const (
	CParticipant  = "participant"
	CAdmin        = "admin"
	CNotification = "notification"
)

func (d *Database) DB() (*mgo.Database, func()) {
	conn := d.session.Copy()
	return conn.DB(d.database), func() {
		conn.Close()
	}
}

func InitGlobalDatabase(url string, database string) error {
	d, err := NewDatabase(url, database)
	if err != nil {
		return err
	}
	GlobalDB = d
	return nil
}

func NewDatabase(url string, database string) (*Database, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	cred := &mgo.Credential{
		Username: "engine",
		Password: "Qwert@789",
	}
	err = session.Login(cred)
	if err != nil {
		panic(err)
	}
	d := &Database{
		session:  session,
		database: database,
	}
	return d, nil
}

func (d *Database) Admin() (collection *mgo.Collection, closeConn func()) {
	database, closeConn := d.DB()
	c := database.C(CAdmin)
	return c, closeConn
}

func (d *Database) Participant() (collection *mgo.Collection, closeConn func()) {
	database, closeConn := d.DB()
	c := database.C(CParticipant)
	return c, closeConn
}

func (d *Database) Notification() (collection *mgo.Collection, closeConn func()) {
	database, closeConn := d.DB()
	c := database.C(CNotification)
	return c, closeConn
}
