package routers

import (
	"JRYY/middleware"
	"JRYY/routers/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

func (r Router)InitRouter(g *gin.Engine) error {
	//go v1.Manager.Start()
	base := g.Group("api")
	base.POST("/upload", api.UploadFile)
	base.POST("/uploadImage", api.UploadImage)
	base.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello api")
	})
	auth := base.Group("/auth")
	auth.POST("/adminLogin", api.AdminLogin)
	auth.POST("/userLogin", api.UserLogin)
	public := base.Group("/public")
	public.Use(middleware.JWTAuthMiddleware())
	public.GET("/info/:gender", middleware.LoginAuth(), api.GetPublicInfo)
	public.GET("/info/:gender/:page", middleware.LoginAuth(), api.GetPublicInfoPage)
	//public.GET("/girlInfo", api.GetPublicInfoGirl)
	participant := base.Group("/participant")
	participant.Use(middleware.JWTAuthMiddleware())
	participant.GET("/bySid/:sid", middleware.UserAuth("sid"), api.GetParticipantBySid)
	participant.PUT("/bySid/:sid", middleware.UserAuth("sid"), api.ModifyParticipantBySid)
	participant.GET("/all", middleware.AdminAuth(), api.GetAllParticipant)
	participant.GET("/initDataByCsv/:path", middleware.AdminAuth(), api.InitParticipantByCsv)
	participant.GET("/addDataByCsv/:path", middleware.AdminAuth(), api.AddParticipantByCsv)
	participant.GET("/addPairByCsv/:path", middleware.AdminAuth(), api.AddParticipantPairByCsv)
	participant.POST("/removePair", middleware.AdminAuth(), api.RemoveParticipantPair)
	participant.POST("/addPair", middleware.AdminAuth(), api.AddParticipantPair)
	participant.GET("/downloadReport", middleware.AdminAuth(), api.GetReportDownload)
	notification := base.Group("/noti")
	notification.Use(middleware.JWTAuthMiddleware())
	notification.POST("/publishAllNoti",middleware.AdminAuth(), api.PublishAllNoti)
	notification.GET("/getAdminNoti",middleware.AdminAuth(), api.GetAdminNoti)
	notification.GET("/getLastAdminAllNoti", api.GetLastAdminAllNoti)
	notification.POST("/publishSomeNoti",middleware.AdminAuth(), api.PublishSomeNoti)

	//ws := apiv1.Group("/ws")
	//ws.GET("/",v1.WsHandler)
	//ws.GET("/pong", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "pong",
	//	})
	//})


	return nil
}