package routers

import (
	_ "VideoWeb/docs"
	"VideoWeb/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CollectRouter(r *gin.Engine) {

	//swagger配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//路由规则
	//私有方法
	r.GET("/userInfo/user-IPInfo", service.GetUserIpInfo)
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
		video.POST("/upload", service.UploadVideo)
		video.GET("/download", service.DownloadVideo)
		video.DELETE("/delete", service.DeleteVideo)
	}

	//管理员私有方法

}
