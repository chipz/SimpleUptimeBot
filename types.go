package main

import "github.com/satori/go.uuid"

var Websites []Website
var TempWebsites []Website

type Website struct {
	Id       uuid.UUID `json:"id"`
	Url      string `json:"url"`
	Interval int    `json:"interval"` // in seconds
	ChatId   int64  `json:"chat-id"`
}

type ActorStatus struct {
	Website    Website
	Active	   bool
}