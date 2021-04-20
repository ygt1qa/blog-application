package main

import (
	"context"
	"testing"

	"github.com/ygt1qa/blog-application/global"
	proto "github.com/ygt1qa/blog-application/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Test_authServer_Login(t *testing.T) {
	global.ConnectToTestDB()
	pw, _ := bcrypt.GenerateFromPassword([]byte("example"), bcrypt.DefaultCost)
	global.DB.Collection("user").InsertOne(context.Background(), global.User{ID: primitive.NewObjectID(), Email: "test@gmail.com", Username: "Yasuyuki", Password: string(pw)})
	server := authServer{}
	_, err := server.Login(context.Background(), &proto.LoginRequest{Login: "test@gmail.com", Password: "example"})
	if err != nil {
		t.Error("1: An error was returned: ", err.Error())
	}
	_, err = server.Login(context.Background(), &proto.LoginRequest{Login: "something", Password: "something"})
	if err == nil {
		t.Error("2: Error was nil")
	}

	_, err = server.Login(context.Background(), &proto.LoginRequest{Login: "Yasuyuki", Password: "example"})
	if err != nil {
		t.Error("3: Error was nil")
	}
}

func Test_authServer_UsernameUsed(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Username: "Yasuyuki"})
	server := authServer{}
	res, err := server.UsernameUsed(context.Background(), &proto.UsernameUsedRequest{Username: "Yasu"})
	if err != nil {
		t.Error("1: An error was returned: ", err.Error())
	}
	if res.GetUsed() {
		t.Error("1: Wrong result")
	}
	res, err = server.UsernameUsed(context.Background(), &proto.UsernameUsedRequest{Username: "Yasuyuki"})
	if err != nil {
		t.Error("2: An error was returned: ", err.Error())
	}
	if !res.GetUsed() {
		t.Error("2: Wrong result")
	}
}

func Test_authServer_EmailUsed(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Email: "yasu@gmail.com"})
	server := authServer{}
	res, err := server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: "y@gmail.com"})
	if err != nil {
		t.Error("1: An error was returned: ", err.Error())
	}
	if res.GetUsed() {
		t.Error("1: Wrong result")
	}
	res, err = server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: "yasu@gmail.com"})
	if err != nil {
		t.Error("2: An error was returned: ", err.Error())
	}
	if !res.GetUsed() {
		t.Error("2: Wrong result")
	}
}

func Test_authServer_Signup(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Username: "Yasuyuki", Email: "yasu@gmail.com"})
	server := authServer{}
	_, err := server.Signup(context.Background(), &proto.SignupRequest{Username: "Yasuyuki", Email: "example@gmail.com", Password: "examplestring"})
	if err.Error() != "Username is used" {
		t.Error("1: No or the wrong Error was returned")
	}
	_, err = server.Signup(context.Background(), &proto.SignupRequest{Username: "example", Email: "yasu@gmail.com", Password: "examplestring"})
	if err.Error() != "Email is used" {
		t.Error("2: No or the wrong Error was returned")
	}
	_, err = server.Signup(context.Background(), &proto.SignupRequest{Username: "example", Email: "example@gmail.com", Password: "examplestring"})
	if err != nil {
		t.Error("3: No or the wrong Error was returned")
	}
	_, err = server.Signup(context.Background(), &proto.SignupRequest{Username: "example", Email: "example@gmail.com", Password: "examp"})
	if err.Error() != "Validation Fail" {
		t.Error("4: No or the wrong Error was returned")
	}
}

func Test_authServer_AuthUser(t *testing.T) {
	server := authServer{}
	res, err := server.AuthUser(context.Background(), &proto.AuthUserRequest{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoie1wiSURcIjpcIjYwN2FkYzc2ZTYxMTJlYTU5MTE3ZTVkOVwiLFwiVXNlcm5hbWVcIjpcIllhc3V5dWtpXCIsXCJFbWFpbFwiOlwidGVzdEBnbWFpbC5jb21cIixcIlBhc3N3b3JkXCI6XCIkMmEkMTAkaklWRUoveWtqQzhyaWtGa3F2UXZmLmlNdk5zZHhRQUR2aFcvY3lBdlB3TTlDVW1SS2piazJcIn0ifQ.CrcQCA5nVlZP0MvZtI4NE23s40xjASYWq8aq5ne-I-g"})
	if err != nil {
		t.Error("an error was returned")
	}
	if res.GetID() != "607adc76e6112ea59117e5d9" || res.GetUsername() != "Yasuyuki" || res.GetEmail() != "test@gmail.com" {
		t.Error("wrong result returned: ", res)
	}
}
