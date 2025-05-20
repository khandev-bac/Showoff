package main

import (
	"exceapp/pkg/jwt"
	"fmt"

	"github.com/google/uuid"
)

func main() {
	uid := uuid.New()
	fmt.Println("uid::  ", uid)
	jwtToken, err := jwt.GenerateJWTToken(uid)
	if err != nil {
		fmt.Println("error while generating jwt token")
	}
	fmt.Println("token::  ", jwtToken.AccessToken)
	jwtParsed, err := jwt.ValidateToken(jwtToken.AccessToken)
	if err != nil {
		fmt.Println("error while validating jwt token")
	}
	fmt.Println(jwtParsed)
	fmt.Println("only id from valifdatetokrn: ", jwtParsed["userID"])
}
