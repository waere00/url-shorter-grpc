package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	pb "github.com/waere00/url-shorter-grpc/v2/proto"
	"google.golang.org/grpc"
)

var (
	targetAddr string
	HOST       string = "serv_cont"
	PORT       string = "9080"
	routerPORT string = "1080"
)

func init() {
	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
	}
	if os.Getenv("HOST") != "" {
		HOST = os.Getenv("HOST")
	}
	if os.Getenv("routerPORT") != "" {
		routerPORT = os.Getenv("routerPORT")
	}

	targetAddr = HOST + ":" + PORT
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
