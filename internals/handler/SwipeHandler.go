package handler

import "exceapp/internals/service"

type SwipeHandler struct {
	service *service.SwipeService
}

func NewSwipeHandler(service *service.SwipeService) *SwipeHandler {
	return &SwipeHandler{
		service: service,
	}
}
