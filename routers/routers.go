package routers

import (
	"VideoWeb/Utilities/WebSocket"
	_ "VideoWeb/docs"
	"VideoWeb/middlewares"
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
		userInfo := user.Group("", middlewares.CheckIfUserLogin()) //添加登录中间件
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
			favorites := userInfo.Group("/Favorites")
			{
				favorites.POST("/Create", service.CreateFavorites)
				favorites.PUT("/Modify", service.ModifyFavorites)
				favorites.DELETE("/Delete", service.DeleteFavorites)
			}
		}

	}

	//视频相关接口
	video := r.Group("/video")
	{
		video.GET("/OfferMpd", service.OfferMpd)
		video.GET("/DASHStreamTransmission", service.DASHStreamTransmission)
		video.GET("/VideoList", middlewares.CheckIfUserLogin(), service.GetVideoList)
		video.POST("/Upload", middlewares.CheckIfUserLogin(), service.UploadVideo)

		videoInfo := video.Group("/:ID", middlewares.CheckIfUserLogin()) //ID为视频ID
		{
			videoInfo.POST("/AddBarrage", service.AddBarrage)
			videoInfo.POST("/Bookmark", service.BookmarkVideo)
			videoInfo.DELETE("/Bookmark", service.UnBookmarkVideo)
			videoInfo.PUT("/LikeOrUndoLike", service.LikeOrUndoLike)
			videoInfo.PUT("/ThrowShell", service.ThrowShell)
			videoInfo.GET("/StreamTransmission", service.StreamTransmission)
			videoInfo.GET("/Comments", service.GetVideoComments)
			videoInfo.GET("/Download", service.DownloadVideo)
			videoInfo.GET("/VideoDetail", service.GetVideoInfo)
			//videoInfo.GET("/test", service.Test)
			videoInfo.DELETE("/Delete", service.DeleteVideo)
		}
	}

	//评论相关接口
	comment := r.Group("/Comment/:VideoID", middlewares.CheckIfUserLogin())
	{
		comment.POST("/Comment", service.PostComment)
		comment.PUT("/LikeOrDislikeComment", service.LikeOrDislikeComment)
		comment.PUT("/UndoLikeOrDislikeComment", service.UndoLikeOrDislikeComment)
	}

}
