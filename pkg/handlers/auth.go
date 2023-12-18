package handlers

import (
	music "github.com/bear1278/MusicWave"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

const (
	body      = `<p>To change the password click the link: </p>`
	linkBegin = `<a href="http://localhost:8000/auth/reset-pass/`
	linkEnd   = `">Reset</a>`
)

func (h *Handler) signUp(ctx *gin.Context) {
	var user music.User
	if err := ctx.BindJSON(&user); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.service.CreateUser(user)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.service.GenerateToken(user.UserName, user.Password)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":    id,
		"token": token,
	})
}

type signInInput struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(ctx *gin.Context) {
	var user signInInput
	if err := ctx.BindJSON(&user); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.service.GenerateToken(user.UserName, user.Password)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})

}

func (h *Handler) signUpGet(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signUp.html", nil)
}

func (h *Handler) signInGet(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signIn.html", nil)
}

func (h *Handler) mainGet(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

func (h *Handler) recommendationGet(ctx *gin.Context) {
	tmpl, err := template.ParseFiles("public/recommendation.html")
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	genres, err := h.service.FillHtml()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	err = tmpl.Execute(ctx.Writer, genres)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)

}

type genreIds struct {
	Genres []int64 `json:"genres" binding:"required"`
}

func (h *Handler) recommendationPost(ctx *gin.Context) {
	var genresId genreIds
	if err := ctx.ShouldBindJSON(&genresId); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userId := ctx.GetInt64(userCtx)
	genres := make([]music.Genre, len(genresId.Genres))
	for key, genreId := range genresId.Genres {
		genres[key].Id = genreId
	}
	if err := h.service.InsertRecommendation(genres, userId); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

type resetEmailInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (h *Handler) GetEmailForReset(ctx *gin.Context) {
	var input resetEmailInput
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	user, token, err := h.service.Authorization.GenerateTokenForReset(input.Username, input.Email)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if err := SendMail(user.Email, body+linkBegin+token+linkEnd); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

type passwordInput struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

func (h *Handler) SetNewPassword(ctx *gin.Context) {
	var input passwordInput
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.service.Authorization.ParseToken(input.Token)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.service.Authorization.SetNewPassword(id, input.Password); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) resetEmailGet(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "reset_email.html", nil)
}

func (h *Handler) GetPasswordReset(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "reset_pass.html", nil)
}
