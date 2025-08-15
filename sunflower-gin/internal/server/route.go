package server

import (
	"net/http"

	"sunflower-gin/internal/handler/auth"
	"sunflower-gin/internal/handler/checkin"
	"sunflower-gin/internal/handler/points"
	"sunflower-gin/internal/handler/user"
	"sunflower-gin/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SetupRoutes(cfg *viper.Viper) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowHeaders = append(corsCfg.AllowHeaders, "Authorization")
	corsCfg.AllowAllOrigins = true // 允许所有跨域请求，不建议在生产环境使用
	r.Use(cors.New(corsCfg))       // CORS 跨域中间件，简单粗暴，直接放行所有跨域请求
	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/users", user.CreateHandler)     // 创建用户
		apiV1.POST("/auth/login", auth.LoginHandler) // 用户登录
		apiV1.POST("/auth/refresh", auth.RefreshHandler)

		apiV1.Use(middleware.Auth()) // 注册认证中间件
		// 在这个Auth中间件后面的都需要认证通过才能访问
		apiV1.GET("/users/me", user.ProfileHandler) // 获取当前用户信息

		// checkin api group
		checkinGroup := apiV1.Group("/checkins")
		{
			checkinGroup.POST("", checkin.DailyHandler)
			checkinGroup.GET("/calendar", checkin.CalendarHandler)
			checkinGroup.POST("/retroactive", checkin.RetroactiveHandler)
		}
		// points api group
		pointsGroup := apiV1.Group("/points")
		{
			pointsGroup.GET("/summary", points.SummaryHandler)
			pointsGroup.GET("/records", points.RecordsHandler)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
