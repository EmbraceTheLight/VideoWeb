package routers

import (
	"VideoWeb/Utilities/WebSocket"
	_ "VideoWeb/docs"
	"VideoWeb/service"
	"github.com/gin-gonic/gin"
	swaggoFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initWebSocket() {
	hub := WebSocket.NewServerHub()
	go hub.Run()
	router := gin.Default()
	router.GET("/ws/:UserID", service.CreateWebSocket)
	router.Run(":51234")
}

func CollectRouter(r *gin.Engine) {
	//配置websocket管理器
	go initWebSocket()

	//swagger配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggoFiles.Handler))

	//------------路由规则-------------//
	/*公有方法*/
	//验证码相关接口
	captcha := r.Group("/Captcha")
	{
		captcha.GET("/Send-code", service.SendCode)
		captcha.GET("/GenerateGraphicCaptcha", service.GenerateGraphicCaptcha)
		//captcha.POST("/CheckGraphicCaptcha", service.CheckGraphicCaptcha)
	}

	//用户相关接口
	user := r.Group("/User")
	{
		user.POST("/Register", service.Register)
		user.POST("/Login", service.Login)
		userInfo := user.Group("/:UserID")
		{
			userInfo.GET("/User-detail", service.GetUserDetail)
			userInfo.POST("/Fans/Follows", service.FollowOtherUser)
			userInfo.POST("/AddSearchHistory", service.AddSearchHistory)
			userInfo.POST("/AddVideoHistory", service.AddVideoHistory)
			userInfo.PUT("/ModifySignature", service.ModifyUserSignature)
			userInfo.PUT("/Face/Upload/Avatar", service.UploadUserAvatar)
			userInfo.PUT("/ModifyEmail", service.ModifyUserEmail)
			userInfo.PUT("/ModifyPassword", service.ModifyPassword)
			userInfo.PUT("/ModifyUserName", service.ModifyUserName)
			userInfo.PUT("/ForgetPassword", service.ForgetPassword)
			userInfo.DELETE("/Logout", service.Logout)
			userInfo.DELETE("/Fans/Unfollows", service.UnFollowOtherUser)

			//用户收藏夹相关接口
			Favorites := userInfo.Group("/Favorites")
			{
				Favorites.POST("/Create", service.CreateFavorites)
				Favorites.PUT("/Modify", service.ModifyFavorites)
				Favorites.DELETE("/Delete", service.DeleteFavorites)
			}
		}

	}

	//视频相关接口
	video := r.Group("/video")
	{
		video.GET("/OfferMpd", service.OfferMpd)
		video.GET("/DASHStreamTransmission", service.DASHStreamTransmission)

		videoInfo := video.Group("/:ID") //ID在UploadVideo方法中代表用户ID，在其他方法中为视频ID
		{
			videoInfo.POST("/upload", service.UploadVideo)
			videoInfo.POST("/AddBarrage", service.AddBarrage)
			videoInfo.PUT("/likeOrUndoLike", service.LikeOrUndoLike)
			videoInfo.PUT("/throwShell", service.ThrowShell)
			videoInfo.GET("/StreamTransmission", service.StreamTransmission)
			videoInfo.GET("/download", service.DownloadVideo)
			videoInfo.GET("/getVideoDetail", service.GetVideoInfo)
			//videoInfo.GET("/test", service.Test)
			videoInfo.DELETE("/delete", service.DeleteVideo)
		}
	}

	//评论相关接口
	comment := r.Group("/comment")
	{
		comment.POST("/ToVideo", service.CommentToVideo)
		comment.POST("/ToUser", service.CommentToOtherUser)
	}

}
