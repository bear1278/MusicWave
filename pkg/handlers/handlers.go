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
	router.LoadHTMLGlob("web/templates/*html")
	router.Static("/static", "./web")
	main := router.Group("/")
	{
		main.GET("/login", h.LoginHandler)
		main.GET("/relocate", h.GetRelocatePage)
		main.GET("/relocate/token", h.GetClient)
		main.GET("/", h.mainGet)
		auth := main.Group("/auth")
		{
			auth.GET("/sign-up", h.signUpGet)
			auth.GET("/sign-in", h.signInGet)
			auth.GET("/recommendation", h.recommendationGet)
			auth.POST("/recommendation", h.userIdentity, h.recommendationPost) //work
			auth.POST("/sign-up", h.signUp)                                    //work
			auth.POST("/sign-in", h.signIn)                                    //work
			auth.GET("/reset-pass-email", h.resetEmailGet)
			auth.POST("/reset-pass-email", h.GetEmailForReset)
			auth.GET("/reset-pass/:token", h.GetPasswordReset)
			auth.POST("/reset-pass", h.SetNewPassword)
		}
		main.GET("/main", h.GetMainPage)
		main.GET("/recommendations", h.userIdentity, h.spotifyIdentity, h.GetUserRecommendation)
		main.GET("/library", h.GetLibraryPage)
		main.GET("/api/tracks/page/:id", h.GetTrackPage)
		main.GET("/api/albums/page/:id", h.GetAlbumPage)
		main.GET("/api/artists/page/:id", h.GetArtistPage)
		main.GET("/api/playlists/page/:id", h.GetPlaylistPage)
		api := main.Group("/api", h.userIdentity)
		{
			playlists := api.Group("/playlists")
			{
				playlists.POST("/", h.NewPlaylist)           //work
				playlists.DELETE("/:id", h.DeletePlaylist)   //work
				playlists.GET("/", h.GetAllPlaylistsForUser) //work
				playlists.GET("/my", h.GetUsersPlaylists)
				playlists.GET("/:id", h.GetById)                    //work
				playlists.PUT("/:id", h.UpdatePlaylist)             //work
				playlists.DELETE("/exclude/:id", h.ExcludePlaylist) //work
				playlists.POST("/add/:id", h.AddPlaylist)           //work
				playlists.POST("/spotify/:id", h.spotifyIdentity, h.CreatePlaylistInSpotify)

				tracks := playlists.Group("/:id/tracks")
				{
					tracks.POST("/", h.spotifyIdentity, h.AddTrack) //work but need to think about exceptions
					tracks.DELETE("/:id_track", h.ExcludeTrack)     //work
					tracks.GET("/", h.GetAllTracks)                 //work
				}

			}
			api.GET("/tracks/:id_track", h.spotifyIdentity, h.GetTrackById) //work
			api.POST("/tracks/:id", h.spotifyIdentity, h.AddTrackToSpotifyFav)
			albums := api.Group("/albums") //
			{
				albums.POST("/", h.spotifyIdentity, h.AddAlbumToFav)               //work but need to think about exceptions and need optimization
				albums.DELETE("/:id", h.ExcludeAlbum)                              //work
				albums.GET("/:id/tracks", h.spotifyIdentity, h.GetTracksFromAlbum) //work but need to handle situation when albums doesn't have all tracks in db
				albums.GET("/:id", h.spotifyIdentity, h.GetAlbumById)              //work but need to think about exceptions
				albums.GET("/", h.spotifyIdentity, h.GetAllAlbumsForUser)          //work but need to think about exceptions
			}
			artists := api.Group("/artists")
			{
				artists.POST("/", h.spotifyIdentity, h.AddArtistToFav)             //work but need optimization
				artists.DELETE("/:id", h.ExcludeArtist)                            //work
				artists.GET("/:id/albums", h.spotifyIdentity, h.GetAlbumsByArtist) //work
				artists.GET("/:id/tracks", h.spotifyIdentity, h.GetTopTracksByArtist)
				artists.GET("/:id", h.spotifyIdentity, h.GetArtistById) //work
				artists.GET("/", h.GetAllArtistForUser)                 //work but may be should add genres
			}
		}
		search := main.Group("/search")
		{
			search.GET("/all/:string", h.SearchTrack, h.SearchAlbum, h.SearchArtist, h.SearchPlaylist)
			search.GET("/track/:string/:page", h.spotifyIdentity, h.SearchTrack)   //work
			search.GET("/album/:string/:page", h.spotifyIdentity, h.SearchAlbum)   //work
			search.GET("/artist/:string/:page", h.spotifyIdentity, h.SearchArtist) //work
			search.GET("/playlists/:string", h.SearchPlaylist)                     //work
			search.GET("/", h.GetSearchPage)
		}
		profile := main.Group("/profile")
		{
			profile.PATCH("/change-name", h.userIdentity, h.ChangeUsername)   //work
			profile.PATCH("/change-pass", h.userIdentity, h.ChangePassword)   //work
			profile.PATCH("/change-picture", h.userIdentity, h.ChangePicture) //work
			profile.PATCH("/change-email", h.userIdentity, h.ChangeEmail)     //work
			profile.GET("/", h.GetProfilePage)
			profile.GET("/user", h.userIdentity, h.GetUserForProfile)
			profile.GET("/spotify", h.spotifyIdentity, h.GetSpotifyProfileInfo)
			profile.GET("/info", h.userIdentity, h.GetProfileInfo)
		}

		main.GET("/admin/user/page", h.GetAdminUserPage)
		main.GET("/admin/artist/page", h.GetAdminArtistPage)
		main.GET("/admin/playlist/page", h.GetAdminPlaylistPage)
		main.GET("/admin/report/page", h.GetAdminReportPage)
		main.GET("/admin/analytic/page", h.GetAdminAnalyticsPage)
		admin := main.Group("/admin", h.userIdentity, h.CheckAdmin)
		{
			admin.GET("/user", h.GetAllUsers)                      //work
			admin.GET("/playlist", h.GetAllPlaylistsForAdmin)      //work
			admin.GET("/artist", h.GetAllArtists)                  //work
			admin.DELETE("/user/:id", h.DeleteUser)                //work
			admin.DELETE("/playlist/:id", h.DeletePlaylistByAdmin) //work
			admin.DELETE("/artist/:id", h.DeleteArtist)            //work
			admin.GET("/history", h.GetHistory)                    //work
			admin.GET("/history/:format", h.GetReport)
			admin.GET("/export", h.ExportJSON)
			admin.POST("/import", h.ImportJSON)
			admin.GET("/genre", h.GetGenrePopularity)
			admin.GET("/genre/diversity", h.GetGenreDiversity)
			admin.GET("/artist/popularity", h.GetArtistPopularity)
		}
	}
	return router
}
