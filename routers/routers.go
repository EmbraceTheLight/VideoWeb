package routers

import (
	"VideoWeb/Utilities/WebSocket"
	_ "VideoWeb/docs"
	"VideoWeb/logic"
	"VideoWeb/service"
	"github.com/gin-gonic/gin"
	swaggoFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initWebSocket() {
	WebSocket.Hub = WebSocket.NewServerHub()
	go WebSocket.Hub.Run()
	router := gin.Default()
	router.GET("/ws/:UserID", logic.CreateWebSocket)
	router.Run(":51234")
}

func CollectRouter(r *gin.Engine) {
	//配置websocket管理器
	go initWebSocket()

	//swagger配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggoFiles.Handler))

	//路由规则
	//公有方法
	r.POST("/send-code", service.SendCode)
	r.POST("/user/login", service.Login)
	r.POST("/user/register", service.Register)

	//用户相关接口
	userInfo := r.Group("/user")
	{
		userInfo.GET("/user-detail", service.GetUserDetail)
		userInfo.POST("/fans/follows", service.FollowOtherUser)
		userInfo.POST("/face/upload/Avatar", service.UploadUserAvatar)
		userInfo.POST("/AddSearchHistory", service.AddSearchHistory)
		userInfo.POST("/AddVideoHistory", service.AddVideoHistory)
		userInfo.PUT("/ModifySignature", service.ModifyUserSignature)
		userInfo.PUT("/ModifyEmail", service.ModifyUserEmail)
		userInfo.PUT("/ModifyPassword", service.ModifyPassword)
		userInfo.PUT("/ModifyUserName", service.ModifyUserName)
		userInfo.PUT("/ForgetPassword", service.ForgetPassword)
		userInfo.DELETE("/Logout", service.Logout)
		userInfo.DELETE("/fans/unfollows", service.UnFollowOtherUser)
		Favorites := userInfo.Group("/favorites")
		{
			Favorites.POST("/create", service.CreateFavorites)
			Favorites.PUT("/modify", service.ModifyFavorites)
			Favorites.DELETE("/delete", service.DeleteFavorites)
		}
	}

	//视频相关接口
	video := r.Group("/video")
	{
		video.GET("getVideoDetail", service.GetVideoInfo)
		video.GET("/download", service.DownloadVideo)
		video.GET("/StreamTransmission", service.StreamTransmission)
		video.POST("/upload", service.UploadVideo)
		video.POST("/:VideoID/AddBarrage", service.AddBarrage)
		video.DELETE("/delete", service.DeleteVideo)
	}

	//评论相关接口
	comment := r.Group("/comment")
	{
		comment.POST("/ToVideo", service.CommentToVideo)
		comment.POST("/ToUser", service.CommentToOtherUser)
	}
	//管理员私有方法

}
