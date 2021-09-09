package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	pb "github.com/waere00/url-shorter-grpc/v2/proto"
	"google.golang.org/grpc"
)

var (
	targetAddr string
	host       string = "serv_cont"
	port       string = "9080"
	routerPORT string = "1080"
)

func init() {
	if os.Getenv("port") != "" {
		port = os.Getenv("port")
	}
	if os.Getenv("host") != "" {
		host = os.Getenv("host")
	}
	if os.Getenv("routerPORT") != "" {
		routerPORT = os.Getenv("routerPORT")
	}

	targetAddr = host + ":" + port
}

func main() {
	conn, err := grpc.Dial(targetAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %s", err)
	}
	client := pb.NewShorterClient(conn)
	router := gin.Default()

	router.Static("/css", "html/css")
	router.LoadHTMLGlob("html/*.html")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "TestShorter",
		})
	})

	router.POST("/", func(ctx *gin.Context) {
		req := &pb.Url{Url: ctx.PostForm("url")}
		if !strings.HasPrefix(req.Url, "http://") && !strings.HasPrefix(req.Url, "https://") {
			req.Url = "http://" + req.Url
		}
		if !isValidUrl(req.Url) {
			// err := errors.New("bad URL")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad URL"})
			return
		}
		if response, err := client.Create(ctx, req); err == nil {
			ctx.HTML(http.StatusOK, "create.html", gin.H{
				"title": "TestShorter",
				"url":   req.Url,
				"link":  response.Link,
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	router.GET("/get", func(ctx *gin.Context) {
		req := &pb.Link{Link: ctx.Query("link")}
		if response, err := client.Get(ctx, req); err == nil {
			ctx.HTML(http.StatusOK, "get.html", gin.H{
				"title": "TestShorter",
				"url":   response.Url,
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	if err := router.Run("0.0.0.0:" + routerPORT); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// Check if given string is valid URL
func isValidUrl(urlToCheck string) bool {
	_, err := url.ParseRequestURI(urlToCheck)
	if err != nil {
		log.Printf("Bad URL: %s", urlToCheck)
		return false
	}
	return true
}
