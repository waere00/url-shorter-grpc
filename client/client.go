package main

import (
	"context"
	"fmt"
	pb "github.com/waere00/url-shorter-grpc/proto"
	"google.golang.org/grpc"
	"log"
	"os"
)

var (
	targetAddr string
	host       string = "localhost"
	port       string = "9080"
)

func choice(cli pb.ShorterClient) {

	const (
		Create = "1"
		Get    = "2"
		Exit   = "3"
	)
	var start string

	fmt.Println("Выберите действие:\n" +
		"1. Генерировать сокращенную ссылку\n" +
		"2. Использовать короткую ссылку для получения первичной ссылки\n" +
		"3. Выйти")
	fmt.Scan(&start)

	var link string
	switch start {

	case Create:
		fmt.Println("Вставьте ссылку:")
		fmt.Scan(&link)
		result, err := cli.Create(context.Background(), &pb.UrlRequest{Url: link})
		if err != nil {
			log.Fatalln("Ошибка grpcserver, метод Create: ", err)
		}
		fmt.Printf("Короткая ссылка: %s\n", result.Link)
		choice(cli)

	case Get:
		fmt.Println("Вставьте короткую ссылку:")
		fmt.Scan(&link)
		result, err := cli.Get(context.Background(), &pb.LinkRequest{Link: link})
		if err != nil {
			log.Fatalln("Ошибка grpcserver, метод Get: ", err)
		}
		if result.Url != "empty" {
			fmt.Printf("Оригинальная ссылка: %s\n", result.Url)
		} else {
			fmt.Println("Оригинальная ссылка не нашлась")
		}
		choice(cli)

	case Exit:
		fmt.Println("Завершение работы клиента")
		os.Exit(0)

	default:
		fmt.Println(link, "Такой команды нет")
		choice(cli)
	}

}

func init() {
	if os.Getenv("PORT") != "" {
		port = os.Getenv("port")
	}
	if os.Getenv("HOST") != "" {
		host = os.Getenv("host")
	}
	targetAddr = host + ":" + port
}

func main() {
	fmt.Println("Инициализация клиента")
	conn, err := grpc.Dial(targetAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Ошибка при соединении с grpc сервером: ", err)
	} else {
		log.Println("fine!")
	}

	cli := pb.NewShorterClient(conn)

	choice(cli)
}
