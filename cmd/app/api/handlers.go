package api

import (
	"encoding/json"
	"find-ship/storage"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// @Summary GetShips
// @Description returns ships
// @Produce  json
// @Param type query string false "Channel type"
// @Param name query string false "Channel name"
// @Success 200 {object} ChannelsResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/ship [get]
func (s *server) getShips(w http.ResponseWriter, r *http.Request) {
	if name := r.URL.Query().Get("name"); name != "" {
		s.getShipByName(w, r, name)
		return
	}

	if typeShip := r.URL.Query().Get("type"); typeShip != "" {
		s.getShipsByType(w, r, typeShip)
		return
	}

	ships, err := s.storage.GetVKChannels(r.Context())
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, ChannelsResponse{Data: ships})
}

// @Summary GetShipByID
// @Description returns ship by id
// @Produce  json
// @Param id path string true "Channel ID"
// @Success 200 {object} ChannelResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/ship/{id} [get]
func (s *server) getShipByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("missing id"))
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	ship, err := s.storage.GetVKChannel(r.Context(), idInt)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, ChannelResponse{ship})
}

func (s *server) getShipByName(w http.ResponseWriter, r *http.Request, name string) {
	ship, err := s.storage.GetVKChannelByName(r.Context(), name)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, ChannelsResponse{Data: []storage.VKChannel{ship}})
}

func (s *server) getShipsByType(w http.ResponseWriter, r *http.Request, typeShip string) {
	ships, err := s.storage.GetVKChannelsByType(r.Context(), typeShip)
	if err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, ChannelsResponse{Data: ships})
}

// @Summary CreateShip
// @Description creates new ship
// @Accept  json
// @Produce  json
// @Param vk_channel body CreateVKChannelRequest true "VK Channel"
// @Success 201 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /api/v1/ship [post]
func (s *server) postShip(w http.ResponseWriter, r *http.Request) {
	var shipReq CreateVKChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&shipReq); err != nil {
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}
	defer r.Body.Close()

	if err := s.validate(shipReq); err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	if _, err := s.storage.GetVKChannelByName(r.Context(), shipReq.ChannelName); err == nil {
		s.error(w, r, http.StatusUnprocessableEntity, fmt.Errorf("channel with name %s already exists", shipReq.ChannelName))
		return
	}

	suteURL, err := s.vkClient.GetChannelSiteUrl(r.Context(), shipReq.ChannelURL)
	if err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	ship := storage.VKChannel{
		ChannelName: &shipReq.ChannelName,
		ChannelURL:  &shipReq.ChannelURL,
		ChannelType: &shipReq.ChannelType,
		SiteURL:     &suteURL,
	}

	id, err := s.storage.CreateVKChannel(r.Context(), ship)
	if err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	s.respond(w, r, http.StatusCreated, IDResponse{ID: id})
}

// @Summary DeleteShip
// @Description deletes ship by id
// @Success 200
// @Param id path string true "Channel ID"
// @Failure 404 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/ship/{id} [delete]
func (s *server) deleteShip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("missing id"))
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	if err := s.storage.DeleteVKChannel(r.Context(), idInt); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	s.respond(w, r, http.StatusOK, nil)
}

// @Summary UpdateShip
// @Description updates ship
// @Accept  json
// @Param id path string true "Channel ID"
// @Param vk_channel body PatchVKChannelRequest true "VK Channel"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /api/v1/ship/{id} [patch]
func (s *server) patchShip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("missing id"))
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	var shipReq PatchVKChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&shipReq); err != nil {
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}
	defer r.Body.Close()

	if err := s.validate(shipReq); err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	ship := storage.VKChannel{
		ChannelName: shipReq.ChannelName,
		ChannelURL:  shipReq.ChannelURL,
		ChannelType: shipReq.ChannelType,
		SiteURL:     shipReq.SiteURL,
		Id:          idInt,
	}

	if err := s.storage.UpdateVKChannel(r.Context(), ship); err != nil {
		if strings.Contains(err.Error(), "not found") {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	s.respond(w, r, http.StatusOK, nil)
}

// @Summary HealthCheck
// @Description health check
// @Success 200 {string} string "OK"
// @Produce plain
// @Router /health [get]
func (s *server) healthCheck(w http.ResponseWriter, r *http.Request) {
	s.respond(w, r, http.StatusOK, "OK")
}

// @Summary Metrics
// @Description returns service metrics
// @Produce json
// @Success 200 {object} MetricsResponse
// @Router /metrics [get]
func (s *server) getMetrics(w http.ResponseWriter, r *http.Request) {
	mResp := &MetricsResponse{
		s.metrix.GetAll(),
	}

	s.respond(w, r, http.StatusOK, mResp)
}
