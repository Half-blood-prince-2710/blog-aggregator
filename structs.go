package main

import "github.com/half-blood-prince-2710/blog-aggregator/internal/config"


type state struct {
	config *config.Config
}

type command struct {
	name string
	arguments []string
} 

func hanslerLogin