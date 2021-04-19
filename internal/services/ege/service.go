package ege

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"server/config"
	"server/pkg/helpers"
	"server/pkg/middlewares"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Service struct {
	Config *config.Config
}

// NewService creates a new EGE Service
func NewService(cfg *config.Config) *Service {
	return &Service{
		Config: cfg,
	}
}

// Register Service in a provided router
func (serv *Service) Register(r *gin.RouterGroup) {
	// there's a bug that causes wrong route handling but since the question parameter is an integer i can not worry about it
	// but for awareness i'll leave it here
	// https://github.com/gin-gonic/gin/issues/2682
	r.POST("/:question/solve", middlewares.EnsureParamIsInt("question"), serv.handleQuestionSolve)
	r.GET("/available", serv.handleAvailable)
	r.GET("/:question/types", middlewares.EnsureParamIsInt("question"), serv.handleQuestionTypes)
}

func (serv *Service) handleQuestionSolve(c *gin.Context) {
	mReader := multipart.NewReader(c.Request.Body, helpers.GetBoundary(c.GetHeader("Content-Type")))
	metadataPart, err := mReader.NextPart()
	if err != nil {
		code := http.StatusBadRequest
		c.JSON(code, gin.H{
			"code":    code,
			"message": "no metadata part provided",
		})
		return
	}
	var reqData question24Request
	err = helpers.ParseBodyPartToJSON(metadataPart, &reqData)
	if err != nil || reqData.Type == 0 {
		code := http.StatusBadRequest
		c.JSON(code, gin.H{
			"code":    code,
			"message": "wrong metadata body",
		})
		return
	}
	textPart, err := mReader.NextPart()
	if err != nil {
		code := http.StatusBadRequest
		c.JSON(code, gin.H{
			"code":    code,
			"message": "no text part provided",
		})
		return
	}
	text, err := io.ReadAll(textPart)
	if err != nil {
		code := http.StatusBadRequest
		c.JSON(code, gin.H{
			"code":    code,
			"message": err.Error(),
		})
		return
	}
	fName := helpers.SaveToFile(serv.Config.TempDir, text)
	fPath := filepath.Join(serv.Config.TempDir, fName)
	defer os.Remove(fPath)
	questionNumber, _ := strconv.Atoi(c.Param("question")) // can ignore the error since middleware validates that param is a number
	result, err := processQuestion(serv.Config.PythonScriptPath, questionNumber, fPath, &reqData)
	if err != nil {
		log.Println(result, err)
		code := http.StatusInternalServerError
		c.JSON(code, gin.H{
			"code":    code,
			"message": "internal server error",
		})
		return
	}
	result = strings.TrimRight(result, "\r\n") // since python prints everything with an endline character we need to trim it
	resultInt, _ := strconv.Atoi(result)
	code := http.StatusOK
	c.JSON(code, gin.H{
		"code":   code,
		"result": resultInt,
	})
}

func (serv *Service) handleAvailable(c *gin.Context) {
	result, err := executeScript(serv.Config.PythonScriptPath, "available")
	if err != nil {
		log.Println(result, err)
		code := http.StatusInternalServerError
		c.JSON(code, gin.H{
			"code":    code,
			"message": "internal server error",
		})
		return
	}
	result = strings.TrimRight(result, "\r\n") // since python prints everything with an endline character we need to trim it
	code := http.StatusOK
	c.JSON(code, gin.H{
		"code":                code,
		"questions_available": result,
	})
}

func (serv *Service) handleQuestionTypes(c *gin.Context) {
	questionNumber := c.Param("question")
	result, err := executeScript(serv.Config.PythonScriptPath, "types", questionNumber)
	if err != nil {
		log.Println(result, err)
		code := http.StatusInternalServerError
		c.JSON(code, gin.H{
			"code":    code,
			"message": "internal server error",
		})
		return
	}
	result = strings.TrimRight(result, "\r\n") // since python prints everything with an endline character we need to trim it
	code := http.StatusOK
	c.JSON(code, gin.H{
		"code":            code,
		"types_available": result,
	})
}
