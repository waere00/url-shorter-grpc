package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	pb "github.com/waere00/url-shorter-grpc/proto"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

var (
	selfAddr   string
	host       string = "localhost"
	port       string = "9080"
	dbHost     string = "localhost"
	dbPort     string = "5433"
	dbName     string = "links_db"
	dbUser     string = "postgres"
	dbPassword string = "postgres"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

var db *sql.DB

func init() {
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	if os.Getenv("HOST") != "" {
		host = os.Getenv("HOST")
	}
	if os.Getenv("DB_HOST") != "" {
		dbHost = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") != "" {
		dbPort = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_NAME") != "" {
		dbName = os.Getenv("DB_NAME")
	}
	if os.Getenv("DB_USER") != "" {
		dbUser = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PASSWORD") != "" {
		dbPassword = os.Getenv("DB_PASSWORD")
	}
	selfAddr = host + ":" + port
}

var conninfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

func main() {
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		log.Fatalf("Can't connect to postgres: %s", err)
	}
	defer db.Close()

	db = GetDB()

	// initializing Server
	log.Println("Initializing Server")
	srv := grpc.NewServer()
	instance := new(ShorterServer)
	pb.RegisterShorterServer(srv, instance)
	log.Println("Successfully initialized")

	listener, err := net.Listen("tcp", selfAddr)
	if err != nil {
		log.Fatalf("Unable to create grpc listener: %s", err)
	} else {
		log.Println("Listening on", selfAddr)
	}
	// starting the Server
	if err = srv.Serve(listener); err != nil {
		log.Fatalf("Unable to start Server: %s", err)
	}
}

func GetDB() *sql.DB {
	if db != nil {
		return db
	} else {
		db, err := sql.Open("postgres", conninfo)
		if err != nil {
			log.Fatalf("Can't connect to postgres: %s", err)
		}
		return db
	}
}

type ShorterServer struct {
}

func (s *ShorterServer) Create(ctx context.Context, req *pb.UrlRequest) (*pb.LinkResponse, error) {
	response := new(pb.LinkResponse)
	genLink := genShortLink()
	db = GetDB()
	//var checkLink string
	//	db.QueryRow("SELECT link FROM links WHERE url = $1;", req.Url).Scan(&checkLink)
	//	if checkLink != "" {
	//		log.Println("Ссылка уже есть в базе")
	//		response.Link = "localshorter.local/" + checkLink
	//		return response, nil
	//	}
	log.Println(req.Url, genLink)
	_, err := db.Exec("INSERT INTO links VALUES ($1, $2);", req.Url, genLink)
	if err != nil {
		log.Fatalf("Unable to insert into table: %s", err)
	}
	response.Link = "localshorter.local/" + genLink
	return response, nil
}

func (s *ShorterServer) Get(ctx context.Context, req *pb.LinkRequest) (*pb.UrlResponse, error) {
	response := new(pb.UrlResponse)

	if strings.Contains(req.Link, "localshorter.local/") {
		req.Link = strings.Replace(req.Link, "localshorter.local/", "", 1)
	}
	db = GetDB()
	db.QueryRow("SELECT url FROM links WHERE link = $1;", req.Link).Scan(&response.Url)
	if response.Url != "" {
		return response, nil
	} else {
		response.Url = "Wrong link!"
		return response, nil
	}
}

func genShortLink() string {
	unique := false
	link := make([]byte, 10)
	for !unique {
		for i := range link {
			link[i] = chars[seededRand.Intn(len(chars))]
		}
		unique = isUnique(string(link))
	}
	return string(link)
}

func isUnique(link string) bool {

	db = GetDB()
	var check string
	db.QueryRow("SELECT EXISTS (SELECT link FROM links WHERE link = $1 LIMIT 1);", link).Scan(&check)
	if check != "false" {
		return false
	}
	return true
}
