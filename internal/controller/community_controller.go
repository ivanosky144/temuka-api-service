package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/repository"
	httputil "github.com/temuka-api-service/pkg/http"
)

type CommunityController interface {
	CreateCommunity(w http.ResponseWriter, r *http.Request)
	JoinCommunity(w http.ResponseWriter, r *http.Request)
}

type CommunityControllerImpl struct {
	CommunityRepository repository.CommunityRepository
}

func NewCommunityController(repo repository.CommunityRepository) CommunityController {
	return &CommunityControllerImpl{
		CommunityRepository: repo,
	}
}

func (c *CommunityControllerImpl) CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"desc"`
		LogoPicture string `json:"logopicture"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	newCommunity := model.Community{
		Name:        requestBody.Name,
		Description: requestBody.Description,
		LogoPicture: requestBody.LogoPicture,
	}

	if err := c.CommunityRepository.CreateCommunity(context.Background(), &newCommunity); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating community"})
		return
	}

	response := struct {
		Message string          `json:"message"`
		Data    model.Community `json:"data"`
	}{
		Message: "Community has been created",
		Data:    newCommunity,
	}
	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *CommunityControllerImpl) JoinCommunity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityIDstr := vars["community_id"]

	communityID, err := strconv.Atoi(communityIDstr)
	if err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid community id"})
		return
	}

	var requestBody struct {
		UserID int `json:"user_id"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	community, err := c.CommunityRepository.GetCommunityDetailByID(context.Background(), communityID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving community"})
		return
	}
	if community == nil {
		httputil.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "Community not found"})
		return
	}

	existingMember, err := c.CommunityRepository.CheckMembership(context.Background(), communityID, requestBody.UserID)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error checking community membership"})
		return
	}
	if existingMember != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "User already a member of the community"})
		return
	}

	newMember := model.CommunityMember{
		UserID:      requestBody.UserID,
		CommunityID: communityID,
	}

	if err := c.CommunityRepository.AddCommunityMember(context.Background(), &newMember); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error adding community member"})
		return
	}

	community.MembersCount++
	if err := c.CommunityRepository.UpdateCommunity(context.Background(), communityID, community); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error updating community"})
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Successfully joined the community",
	}
	httputil.WriteResponse(w, http.StatusOK, response)
}
