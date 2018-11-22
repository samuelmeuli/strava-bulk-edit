package main

type ActivityListResponse struct {
	Activities []Activity `json:"models"`
	Page       int        `json:"page"`
	PerPage    int        `json:"perPage"`
	Total      int        `json:"total"`
}

type Activity struct {
	Date        string `json:"start_time"`
	Id          int    `json:"id"`
	Title       string `json:"name"`
	Description string `json:"description"`
	Commute     bool   `json:"commute"`
	Sport       string `json:"type"`
	Visibility  string `json:"visibility"`
}
