// docker-compose up dbase

package main

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/waere00/url-shorter-grpc/v2/proto"
)

func init() {
	dbHost = "0.0.0.0"
	dbPort = "5432"
	conninfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
}

var respLink *pb.Link

func TestServerShorten(t *testing.T) {
	var err error
	s := &ShorterServer{}
	req := &pb.Url{Url: "http://example.com"}
	respLink, _ = s.Create(context.Background(), req)
	require.NoError(t, err)
	assert.NotEmpty(t, respLink.Link)
}

func TestServerExpand(t *testing.T) {
	s := &ShorterServer{}
	url := "http://example.com"
	reqGet := &pb.Link{Link: respLink.Link}
	resp, err := s.Get(context.Background(), reqGet)
	require.NoError(t, err)
	assert.Equal(t, url, resp.Url)
}

func TestServerExpand_tokenNotFound(t *testing.T) {
	s := &ShorterServer{}
	reqGet := &pb.Link{Link: "123"}
	resp, _ := s.Get(context.Background(), reqGet)
	assert.Equal(t, "No such link in the database", resp.Url)
}
