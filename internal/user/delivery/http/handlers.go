package http

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/JamesHsu333/go-twitter/config"
	"github.com/JamesHsu333/go-twitter/internal/file"
	"github.com/JamesHsu333/go-twitter/internal/follow"
	"github.com/JamesHsu333/go-twitter/internal/like"
	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/internal/session"
	"github.com/JamesHsu333/go-twitter/internal/tweet"
	"github.com/JamesHsu333/go-twitter/internal/user"
	"github.com/JamesHsu333/go-twitter/pkg/csrf"
	"github.com/JamesHsu333/go-twitter/pkg/httpErrors"
	"github.com/JamesHsu333/go-twitter/pkg/logger"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// User handlers
type UserHandlers struct {
	cfg      *config.Config
	userUC   user.UseCase
	sessUC   session.UCSession
	fileUC   file.UseCase
	followUC follow.UseCase
	likeUC   like.UseCase
	tweetUC  tweet.UseCase
	logger   logger.Logger
}

// NewUserHandlers User handlers constructor
func NewUserHandlers(cfg *config.Config, userUC user.UseCase, sessUC session.UCSession, fileUC file.FileRepository,
	followUC follow.UseCase, likeUC like.UseCase, tweetUC tweet.UseCase, log logger.Logger) user.Handlers {
	return &UserHandlers{
		cfg:      cfg,
		userUC:   userUC,
		sessUC:   sessUC,
		fileUC:   fileUC,
		followUC: followUC,
		likeUC:   likeUC,
		tweetUC:  tweetUC,
		logger:   log,
	}
}

// Register godoc
// @Summary Register new user
// @Description register new user, returns user and token
// @Tags User
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Router /users/register [post]
func (h *UserHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.Register", nil)
		defer span.End()

		user := &models.User{}
		if err := utils.ReadRequest(c, user); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdUser, err := h.userUC.Register(ctx, user)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		sess, err := h.sessUC.CreateSession(ctx, &models.Session{
			UserID: createdUser.User.UserID,
		}, h.cfg.Session.Expire)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		c.SetCookie(utils.CreateSessionCookie(h.cfg, sess))

		return c.JSON(http.StatusCreated, createdUser)
	}
}

// Login godoc
// @Summary Login new user
// @Description login user, returns user and set session
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /users/login [post]
func (h *UserHandlers) Login() echo.HandlerFunc {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.Login", nil)
		defer span.End()

		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		userWithToken, err := h.userUC.Login(ctx, &models.User{
			Email:    login.Email,
			Password: login.Password,
		})
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		sess, err := h.sessUC.CreateSession(ctx, &models.Session{
			UserID: userWithToken.User.UserID,
		}, h.cfg.Session.Expire)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		c.SetCookie(utils.CreateSessionCookie(h.cfg, sess))

		return c.JSON(http.StatusOK, userWithToken)
	}
}

// Logout godoc
// @Summary Logout user
// @Description logout user removing session
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"ok"
// @Router /users/logout [post]
func (h *UserHandlers) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.Logout", nil)
		defer span.End()

		cookie, err := c.Cookie("session-id")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				tracer.AddSpanError(span, err)
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(err))
			}
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(err))
		}

		if err := h.sessUC.DeleteByID(ctx, cookie.Value); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		utils.DeleteSessionCookie(c, h.cfg.Session.Name)

		return c.NoContent(http.StatusOK)
	}
}

// Update godoc
// @Summary Update user
// @Description update existing user
// @Tags User
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 200 {object} models.User
// @Router /users/{id} [patch]
func (h *UserHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.Update", nil)
		defer span.End()

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		user := &models.User{}
		user.UserID = uID

		if err = utils.ReadRequest(c, user); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		updatedUser, err := h.userUC.Update(ctx, user)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}

// GetUserByID godoc
// @Summary get user by id
// @Description get string by ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "user_id"
// @Success 200 {object} models.User
// @Failure 500 {object} httpErrors.RestError
// @Router /users/{id} [get]
func (h *UserHandlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.GetUserByID", nil)
		defer span.End()

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		self, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		cachedUser, _ := h.userUC.GetCacheByID(ctx, self.UserID, uID)
		if cachedUser != nil {
			return c.JSON(http.StatusOK, cachedUser)
		}

		user, err := h.userUC.GetByID(ctx, self.UserID, uID)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

// GetUserByUserName godoc
// @Summary get user by user name
// @Description get string by user name
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "user_name"
// @Success 200 {object} models.User
// @Failure 500 {object} httpErrors.RestError
// @Router /users/username/{user_name} [get]
func (h *UserHandlers) GetUserByUserName() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.GetUserByUserName", nil)
		defer span.End()

		userName := c.Param("user_name")

		self, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		user, err := h.userUC.GetByUserName(ctx, self.UserID, userName)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

// Delete
// @Summary Delete user account
// @Description some description
// @Tags User
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 204 {string} string	"ok"
// @Failure 500 {object} httpErrors.RestError
// @Router /users/{id} [delete]
func (h *UserHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.Delete", nil)
		defer span.End()

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.userUC.Delete(ctx, uID); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusNoContent)
	}
}

// FindByName godoc
// @Summary Find by name
// @Description Find user by name
// @Tags User
// @Accept json
// @Param name query string false "username" Format(username)
// @Produce json
// @Success 200 {object} models.UsersList
// @Failure 500 {object} httpErrors.RestError
// @Router /users/by [get]
func (h *UserHandlers) FindByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.FindByName", nil)
		defer span.End()

		if c.QueryParam("name") == "" {
			utils.LogResponseError(c, h.logger, httpErrors.NewBadRequestError("name is required"))
			return c.JSON(http.StatusBadRequest, httpErrors.NewBadRequestError("name is required"))
		}

		self, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		response, err := h.userUC.FindByName(ctx, self.UserID, c.QueryParam("name"), paginationQuery)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, response)
	}
}

// GetUsers godoc
// @Summary Get users
// @Description Get the list of all users
// @Tags User
// @Accept json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Produce json
// @Success 200 {object} models.UsersList
// @Failure 500 {object} httpErrors.RestError
// @Router /users [get]
func (h *UserHandlers) GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.GetUsers", nil)
		defer span.End()

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		self, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		usersList, err := h.userUC.GetUsers(ctx, self.UserID, paginationQuery)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, usersList)
	}
}

// GetMe godoc
// @Summary Get user by id
// @Description Get current user by id
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 500 {object} httpErrors.RestError
// @Router /users/me [get]
func (h *UserHandlers) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.GetMe", nil)
		defer span.End()

		user, ok := c.Get("user").(*models.User)
		if !ok {
			utils.LogResponseError(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			return utils.ErrResponseWithLog(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		return c.JSON(http.StatusOK, user)
	}
}

// GetCSRFToken godoc
// @Summary Get CSRF token
// @Description Get CSRF token, required user session cookie
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {string} string "Ok"
// @Failure 500 {object} httpErrors.RestError
// @Router /users/token [get]
func (h *UserHandlers) GetCSRFToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.GetCSRFToken", nil)
		defer span.End()

		sid, ok := c.Get("sid").(string)
		if !ok {
			utils.LogResponseError(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			return utils.ErrResponseWithLog(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}
		token := csrf.MakeToken(sid, h.logger)
		c.Response().Header().Set(csrf.CSRFHeader, token)
		c.Response().Header().Set("Access-Control-Expose-Headers", csrf.CSRFHeader)

		return c.NoContent(http.StatusOK)
	}
}

// Update godoc
// @Summary Update user
// @Description update existing user's role
// @Tags User
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 200 {object} models.User
// @Router /users/{id}/role [patch]
func (h *UserHandlers) UpdateRole() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.UpdateRole", nil)
		defer span.End()

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		user := &models.User{}
		user.UserID = uID

		if err = utils.ReadRequest(c, user); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		updatedUser, err := h.userUC.UpdateRole(ctx, user)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}

// UploadAvatar godoc
// @Summary Post avatar
// @Description Post user avatar image
// @Tags User
// @Accept  json
// @Produce  json
// @Param file formData file true "Body with image file"
// @Param bucket query string true "aws s3 bucket" Format(bucket)
// @Param id path int true "user_id"
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpErrors.RestError
// @Router /users/{id}/avatar [post]
func (h *UserHandlers) UploadAvatar() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.UploadAvatar", nil)
		defer span.End()

		user, ok := c.Get("user").(*models.User)
		if !ok {
			utils.LogResponseError(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			return utils.ErrResponseWithLog(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		image, err := utils.ReadImage(c, "file")
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		file, err := image.Open()
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		defer file.Close()

		binaryImage := bytes.NewBuffer(nil)
		if _, err = io.Copy(binaryImage, file); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		contentType, err := utils.CheckImageFileContentType(binaryImage.Bytes())
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		reader := bytes.NewReader(binaryImage.Bytes())

		filepath, err := h.fileUC.PutObject(ctx, models.UploadInput{
			File:        reader,
			Name:        image.Filename,
			Size:        image.Size,
			ContentType: contentType,
		})
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if user.Avatar != nil {
			if err = h.fileUC.RemoveObject(ctx, *user.Avatar); err != nil {
				tracer.AddSpanError(span, err)
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(httpErrors.ErrorResponse(err))
			}
		}

		user.Avatar = filepath

		updatedUser, err := h.userUC.Update(ctx, user)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}

// UploadHeader godoc
// @Summary Post header
// @Description Post user header image
// @Tags User
// @Accept  json
// @Produce  json
// @Param file formData file true "Body with image file"
// @Param bucket query string true "aws s3 bucket" Format(bucket)
// @Param id path int true "user_id"
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpErrors.RestError
// @Router /users/{id}/header [post]
func (h *UserHandlers) UploadHeader() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.UploadHeader", nil)
		defer span.End()

		user, ok := c.Get("user").(*models.User)
		if !ok {
			utils.LogResponseError(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			return utils.ErrResponseWithLog(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		image, err := utils.ReadImage(c, "file")
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		file, err := image.Open()
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		defer file.Close()

		binaryImage := bytes.NewBuffer(nil)
		if _, err = io.Copy(binaryImage, file); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		contentType, err := utils.CheckImageFileContentType(binaryImage.Bytes())
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		reader := bytes.NewReader(binaryImage.Bytes())

		filepath, err := h.fileUC.PutObject(ctx, models.UploadInput{
			File:        reader,
			Name:        image.Filename,
			Size:        image.Size,
			ContentType: contentType,
		})
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if user.Header != nil {
			if err = h.fileUC.RemoveObject(ctx, *user.Header); err != nil {
				tracer.AddSpanError(span, err)
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(httpErrors.ErrorResponse(err))
			}
		}

		user.Header = filepath

		updatedUser, err := h.userUC.Update(ctx, user)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}

func (h *UserHandlers) GetTweetsByUserID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.GetTweetsByUserID", nil)
		defer span.End()

		userID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		tweets, err := h.tweetUC.GetTweetsByUserID(ctx, userID, paginationQuery)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, tweets)
	}
}

// Follow godoc
// @Summary Follow other user
// @Description Follow other user
// @Tags User
// @Accept  json
// @Produce  json
// @Success 204 {string} string	"ok"
// @Failure 500 {object} httpErrors.RestError
// @Router /users/{id} [post]
func (h *UserHandlers) Follow() echo.HandlerFunc {
	type Following struct {
		UserID uuid.UUID `json:"user_id" db:"user_id" validate:"omitempty"`
	}
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.Follow", nil)
		defer span.End()

		followerID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		following := &Following{}
		if err := utils.ReadRequest(c, following); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.followUC.Follow(ctx, followerID, following.UserID); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusCreated)
	}
}

func (h *UserHandlers) GetFollowers() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.GetFollowers", nil)
		defer span.End()

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		self, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		usersList, err := h.followUC.GetFollowers(ctx, self.UserID, uID, paginationQuery)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, usersList)
	}
}

func (h *UserHandlers) GetFollowing() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.GetFollowing", nil)
		defer span.End()

		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		self, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		usersList, err := h.followUC.GetFollowing(ctx, self.UserID, uID, paginationQuery)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, usersList)
	}
}

func (h *UserHandlers) DeleteFollowing() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.DeleteFollowing", nil)
		defer span.End()

		followerID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		followingID, err := uuid.Parse(c.Param("following_id"))
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.followUC.Delete(ctx, followerID, followingID); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (h *UserHandlers) Like() echo.HandlerFunc {
	type Like struct {
		TweetID uint64 `json:"tweet_id" db:"tweet_id" validate:"omitempty"`
	}
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.Like", nil)
		defer span.End()

		userID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		like := &Like{}
		if err := utils.ReadRequest(c, like); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.likeUC.Like(ctx, userID, like.TweetID); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusCreated)
	}
}

func (h *UserHandlers) GetLikedTweets() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.GetLikedTweets", nil)
		defer span.End()

		userID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		tweets, err := h.likeUC.GetLikedTweets(ctx, userID, paginationQuery)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, tweets)
	}
}

func (h *UserHandlers) DeleteLiked() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "UserHandlers.DeleteLiked", nil)
		defer span.End()

		userID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		tweetID, err := strconv.ParseUint(c.Param("tweet_id"), 10, 64)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.likeUC.Delete(ctx, userID, tweetID); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusNoContent)
	}
}
