package http

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/JamesHsu333/go-twitter/config"
	"github.com/JamesHsu333/go-twitter/internal/file"
	"github.com/JamesHsu333/go-twitter/internal/like"
	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/internal/tweet"
	"github.com/JamesHsu333/go-twitter/pkg/httpErrors"
	"github.com/JamesHsu333/go-twitter/pkg/logger"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Tweet handlers
type TweetHandlers struct {
	cfg     *config.Config
	tweetUC tweet.UseCase
	fileUC  file.FileRepository
	likeUC  like.Repository
	logger  logger.Logger
}

// NewTweetHandlers User handlers constructor
func NewTweetHandlers(cfg *config.Config, tweetUC tweet.UseCase, fileUC file.FileRepository, likeUC like.Repository, logger logger.Logger) tweet.Handlers {
	return &TweetHandlers{cfg: cfg, tweetUC: tweetUC, fileUC: fileUC, likeUC: likeUC, logger: logger}
}

// Create godoc
// @Summary Create new tweet
// @Description create new tweet, returns tweet
// @Tags Tweet
// @Accept file formData file true "Body with image file"
// @Produce json
// @Success 201 {object} models.Tweet
// @Router /tweets [post]
func (h *TweetHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "TweetHandlers.Create", nil)
		defer span.End()

		tweet := &models.Tweet{}
		if err := utils.ReadRequest(c, tweet); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		image, err := utils.ReadImage(c, "image")
		if err != nil {
			if !strings.Contains(err.Error(), "no such file") {
				tracer.AddSpanError(span, err)
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(httpErrors.ErrorResponse(err))
			}
		}

		if image != nil {
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

			imageURL, err := h.fileUC.PutObject(ctx, models.UploadInput{
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

			tweet.Image = imageURL

		}

		createdTweet, err := h.tweetUC.Create(ctx, tweet)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdTweet)
	}
}

// CreateReply godoc
// @Summary Create new reply tweet
// @Description create new reply tweet, returns tweet
// @Tags Tweet
// @Accept file formData file true "Body with image file"
// @Produce json
// @Success 201 {object} models.Tweet
// @Router /tweets/{id}/reply [post]
func (h *TweetHandlers) CreateReply() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "TweetHandlers.CreateReply", nil)
		defer span.End()

		tweetID, err := strconv.ParseUint(c.Param("tweet_id"), 10, 64)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		tweet := &models.Tweet{}
		if err := utils.ReadRequest(c, tweet); err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		image, err := utils.ReadImage(c, "image")
		if err != nil {
			if !strings.Contains(err.Error(), "no such file") {
				tracer.AddSpanError(span, err)
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(httpErrors.ErrorResponse(err))
			}
		}

		if image != nil {
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

			imageURL, err := h.fileUC.PutObject(ctx, models.UploadInput{
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

			tweet.Image = imageURL

		}

		createdTweet, err := h.tweetUC.CreateReply(ctx, tweetID, tweet)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdTweet)
	}
}

// GetTweetByID godoc
// @Summary get tweet by id
// @Description get tweet by ID
// @Tags Tweet
// @Accept  json
// @Produce  json
// @Param id path int true "tweet_id"
// @Success 200 {object} models.Tweet
// @Failure 500 {object} httpErrors.RestError
// @Router /tweets/{id} [get]
func (h *TweetHandlers) GetTweetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "TweetHandlers.GetTweetByID", nil)
		defer span.End()

		tweetID, err := strconv.ParseUint(c.Param("tweet_id"), 10, 64)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		tweet, err := h.tweetUC.GetTweetByID(ctx, tweetID)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, tweet)
	}
}

// GetTweets godoc
// @Summary Get tweets
// @Description Get the list of all tweets
// @Tags Tweet
// @Accept json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Produce json
// @Success 200 {object} models.TweetsList
// @Failure 500 {object} httpErrors.RestError
// @Router /tweets [get]
func (h *TweetHandlers) GetTweets() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "TweetHandlers.GetTweets", nil)
		defer span.End()

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		var tweetsList *models.TweetsList

		if userID := c.QueryParam("userID"); userID != "" {
			uID, err := uuid.Parse(userID)
			if err != nil {
				tracer.AddSpanError(span, err)
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(httpErrors.ErrorResponse(err))
			}
			tweetsList, err = h.tweetUC.GetTweetsByUserID(ctx, uID, paginationQuery)
			if err != nil {
				tracer.AddSpanError(span, err)
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(httpErrors.ErrorResponse(err))
			}
		} else {
			tweetsList, err = h.tweetUC.GetTweets(ctx, paginationQuery)
			if err != nil {
				tracer.AddSpanError(span, err)
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(httpErrors.ErrorResponse(err))
			}
		}

		return c.JSON(http.StatusOK, tweetsList)
	}
}

// GetReplyTweets godoc
// @Summary Get reply tweets by tweet id
// @Description Get the list of reply tweets by tweet id
// @Tags Tweet
// @Accept json
// @Param id path int true "tweet_id"
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Produce json
// @Success 200 {object} models.TweetsList
// @Failure 500 {object} httpErrors.RestError
// @Router /tweets/{id}/replys [get]
func (h *TweetHandlers) GetReplyTweets() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "TweetHandlers.GetReplyTweets", nil)
		defer span.End()

		tweetID, err := strconv.ParseUint(c.Param("tweet_id"), 10, 64)
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

		tweetsList, err := h.tweetUC.GetReplyTweets(ctx, tweetID, paginationQuery)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, tweetsList)
	}
}

// Delete
// @Summary Delete tweet by id
// @Description delete tweet by id
// @Tags Tweet
// @Accept json
// @Param id path int true "tweet_id"
// @Produce json
// @Success 204 {string} string	"ok"
// @Failure 500 {object} httpErrors.RestError
// @Router /tweets/{id} [delete]
func (h *TweetHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "TweetHandlers.Delete", nil)
		defer span.End()

		tweetID, err := strconv.ParseUint(c.Param("tweet_id"), 10, 64)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		tweet, err := h.tweetUC.GetTweetByID(ctx, tweetID)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if tweet.Image != nil {
			if err = h.fileUC.RemoveObject(ctx, *tweet.Image); err != nil {
				tracer.AddSpanError(span, err)
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(httpErrors.ErrorResponse(err))
			}
		}

		err = h.tweetUC.Delete(ctx, tweetID)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (h *TweetHandlers) GetLikedUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(utils.GetRequestCtx(c), "TweetHandlers.GetLikedUsers", nil)
		defer span.End()

		tweetID, err := strconv.ParseUint(c.Param("tweet_id"), 10, 64)
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

		users, err := h.likeUC.GetLikedUsers(ctx, self.UserID, tweetID, paginationQuery)
		if err != nil {
			tracer.AddSpanError(span, err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, users)
	}
}
