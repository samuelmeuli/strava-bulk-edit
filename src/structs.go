package main

type ActivityList struct {
	Activities []Activity `json:"models"`
}

type Activity struct {
	Id          int    `json:"id"`
	Title       string `json:"name"`
	Description string `json:"description"`
	Commute     bool   `json:"commute"`
	Sport       string `json:"type"`
	Visibility  string `json:"visibility"`
}
