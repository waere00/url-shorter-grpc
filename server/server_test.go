// docker-compose up dbase

package main

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	pb "github.com/waere00/url-shorter-grpc/v2/proto"
)

func init() {
	dbHost = "0.0.0.0"
	dbPort = "5432"
	conninfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
}

var respLink *pb.Link

func TestCreate(t *testing.T) {
	var err error
	s := &ShorterServer{}
	req := &pb.Url{Url: "http://example.com"}
	respLink, _ = s.Create(context.Background(), req)
	if err != nil {
		t.Error("TestCreate error", err)
	}
	if respLink.Link == "" {
		t.Error("TestCreate error, exptected a Link, got nothing: ", err)
	}
}

func TestGet(t *testing.T) {
	s := &ShorterServer{}
	url := "http://example.com"
	reqGet := &pb.Link{Link: respLink.Link}
	resp, err := s.Get(context.Background(), reqGet)
	if err != nil {
		t.Error("TestGet error", err)
	}
	if url != resp.Url {
		t.Errorf("TestGet error, exptected %s, but got %s: %s", url, resp.Url, err)
	}
}

func TestGet_LinkNotFound(t *testing.T) {
	s := &ShorterServer{}
	reqGet := &pb.Link{Link: "qwert"}
	resp, err := s.Get(context.Background(), reqGet)
	if err != nil {
		t.Error("TestGet_LinkNotFound error", err)
	}
	if resp.Url != "No such link in the database" {
		t.Errorf("TestGet_LinkNotFound error, got %s: %s", resp.Url, err)
	}
}
