package server

import (
	"JRYY/routers"
	"JRYY/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	Addr string
	g    *gin.Engine
}

func NewServer(addr string) *Server {
	return &Server{
		Addr: addr,
		g:    gin.New(),
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (s *Server)Init() (err error) {
	s.g.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	s.g.Static("/image", utils.BASE_PATH + "image")
	s.g.Use(CORSMiddleware())
	r := routers.NewRouter()
	err = r.InitRouter(s.g)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) StartServer() error {
	err := s.g.Run(s.Addr)
	if err != nil {
		return err
	}
	return nil
}