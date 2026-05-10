package Env

import (
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	TargetUrl string
	JWTSecret string
}

func NewEnv() *Environment {
	_ = godotenv.Load()

	return &Environment{
		TargetUrl: os.Getenv("Target_Url"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
