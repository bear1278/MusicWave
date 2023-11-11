package handlers

import (
	"github.com/bear1278/MusicWave/pkg/service"
	"github.com/gin-gonic/gin"
	_ "html/template"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("public/*html")
	router.Static("/static", "./front-end/")

	main := router.Group("/")
	{
		main.GET("/", h.mainGet)
		auth := main.Group("/auth")
		{
			auth.GET("/sign-up", h.signUpGet)
			auth.GET("/sign-in", h.signInGet)
			auth.GET("/recommendation", h.recommendationGet)
			auth.POST("/recommendation", h.userIdentity, h.recommendationPost) //test this endpoint
			auth.POST("/sign-up", h.signUp)
			auth.POST("/sign-in", h.signIn)
		}

		api := main.Group("/api", h.userIdentity)
		{
			playlists := api.Group("/playlists")
			{
				playlists.POST("/", h.NewPlaylist)
				playlists.DELETE("/:id", h.DeletePlaylist)
				playlists.GET("/", h.GetAllPlaylists)
				playlists.GET("/:id", h.GetById)
				playlists.PUT("/:id", h.UpdatePlaylist)
				playlists.DELETE("/exclude/:id", h.ExcludePlaylist)
				playlists.POST("/add/:id", h.AddPlaylist)

				tracks := playlists.Group("/:id/tracks")
				{
					tracks.POST("/", h.AddTrack)
					tracks.DELETE("/:id")
					tracks.GET("/")
					tracks.GET("/:id")
				}
			}

		}
	}
	return router
}
