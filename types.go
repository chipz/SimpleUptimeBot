package main

var Websites []Website

type Website struct {
	Id       int    `json:"id"`
	Url      string `json:"url"`
	Interval int    `json:"interval"` // in seconds
	ChatId   int64  `json:"chat-id"`
}