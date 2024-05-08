package api

import (
	"bytes"
	"context"
	"encoding/json"
	vkclient "find-ship/cmd/app/vk_client"
	"find-ship/storage"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestGetShips(t *testing.T) {
	t.Parallel()

	now := time.Now()

	var testcases = []struct {
		name       string
		input      string
		response   any
		statusCode int
	}{
		{
			name:       "name query parameter is present and valid",
			input:      "/api/v1/ship?name=channel1",
			statusCode: http.StatusOK,
			response: ChannelsResponse{
				Data: []storage.VKChannel{
					{
						Id:          1,
						ChannelName: pointer("channel1"),
						ChannelURL:  pointer("url1"),
						ChannelType: pointer("type1"),
						SiteURL:     pointer("site1"),
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
			},
		},
		{
			name:       "name query parameter is present, but no ships found",
			input:      "/api/v1/ship?name=test",
			statusCode: http.StatusNotFound,
			response:   ErrorResponse{Error: "channel with name test not found"},
		},
		{
			name:       "type query parameter is present and valid",
			input:      "/api/v1/ship?type=type2",
			statusCode: http.StatusOK,
			response: ChannelsResponse{
				Data: []storage.VKChannel{
					{
						Id:          2,
						ChannelName: pointer("channel2"),
						ChannelURL:  pointer("url2"),
						ChannelType: pointer("type2"),
						SiteURL:     pointer("site2"),
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
			},
		},
		{
			name:       "type query parameter is present, but no ships found",
			input:      "/api/v1/ship?type=test",
			statusCode: http.StatusNotFound,
			response:   ErrorResponse{Error: "no channels found for type test"},
		},
		{
			name:       "neither name nor type query parameter is present",
			input:      "/api/v1/ship",
			statusCode: http.StatusOK,
			response: ChannelsResponse{
				Data: []storage.VKChannel{
					{
						Id:          1,
						ChannelName: pointer("channel1"),
						ChannelURL:  pointer("url1"),
						ChannelType: pointer("type1"),
						SiteURL:     pointer("site1"),
						CreatedAt:   now,
						UpdatedAt:   now,
					},
					{
						Id:          2,
						ChannelName: pointer("channel2"),
						ChannelURL:  pointer("url2"),
						ChannelType: pointer("type2"),
						SiteURL:     pointer("site2"),
						CreatedAt:   now,
						UpdatedAt:   now,
					},
					{
						Id:          3,
						ChannelName: pointer("channel3"),
						ChannelURL:  pointer("url3"),
						ChannelType: pointer("type3"),
						SiteURL:     pointer("site3"),
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
			},
		},
	}

	s := &server{
		storage: NewMockStorage(now),
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r, err := http.NewRequest("GET", tc.input, nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			s.getShips(w, r)
			if w.Code != tc.statusCode {
				t.Errorf("Expected status code %d, got %d", tc.statusCode, w.Code)
			}
			if tc.statusCode >= 200 && tc.statusCode < 300 {
				data := ChannelsResponse{}
				if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
					t.Fatal(err)
				}
				if !cmp.Equal(data, tc.response) {
					t.Errorf("handler returned unexpected body: got %v want %v", data, tc.response)
				}
			} else {
				data := ErrorResponse{}
				if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
					t.Fatal(err)
				}
				if !cmp.Equal(data, tc.response) {
					t.Errorf("handler returned unexpected body: got %v want %v", data, tc.response)
				}
			}
		})
	}
}

func TestGetShipByID(t *testing.T) {
	t.Parallel()

	now := time.Now()

	var testCases = []struct {
		name       string
		input      string
		response   any
		statusCode int
	}{
		{
			name:       "valid id in url",
			input:      "1",
			statusCode: http.StatusOK,
			response: ChannelResponse{
				storage.VKChannel{
					Id:          1,
					ChannelName: pointer("channel1"),
					ChannelURL:  pointer("url1"),
					ChannelType: pointer("type1"),
					SiteURL:     pointer("site1"),
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		{
			name:       "invalid id in url",
			input:      "a",
			statusCode: http.StatusBadRequest,
			response:   ErrorResponse{Error: "invalid id"},
		},
		{
			name:       "non-existent id in url",
			input:      "4",
			statusCode: http.StatusNotFound,
			response:   ErrorResponse{Error: "channel with id 4 not found"},
		},
	}

	s := &server{
		storage: NewMockStorage(now),
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r, err := http.NewRequest("GET", "/api/v1/ship", nil)
			if err != nil {
				t.Fatal(err)
			}
			vars := map[string]string{
				"id": tc.input,
			}
			r = mux.SetURLVars(r, vars)
			w := httptest.NewRecorder()
			s.getShipByID(w, r)
			if w.Code != tc.statusCode {
				t.Errorf("Expected status code %d, got %d", tc.statusCode, w.Code)
			}
			if tc.statusCode >= 200 && tc.statusCode < 300 {
				data := ChannelResponse{}
				if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
					t.Fatal(err)
				}
				if !cmp.Equal(data, tc.response) {
					t.Errorf("handler returned unexpected body: got %v want %v", data, tc.response)
				}
			} else {
				data := ErrorResponse{}
				if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
					t.Fatal(err)
				}
				if !cmp.Equal(data, tc.response) {
					t.Errorf("handler returned unexpected body: got %v want %v", data, tc.response)
				}
			}
		})
	}
}

func TestPostShip(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		name       string
		input      any
		statusCode int
		response   any
	}{
		{
			name: "valid ship data in request",
			input: CreateVKChannelRequest{
				ChannelName: "FanSeries",
				ChannelURL:  "https://vk.com/myserianet",
				ChannelType: "series",
			},
			statusCode: http.StatusCreated,
			response:   IDResponse{ID: 4},
		},
		{
			name: "invalid channel url in request",
			input: CreateVKChannelRequest{
				ChannelName: "channel1",
				ChannelURL:  "url1",
				ChannelType: "type1",
			},
			statusCode: http.StatusUnprocessableEntity,
			response:   ErrorResponse{Error: "invalid channel url"},
		},
		{
			name: "channel from request already exists",
			input: CreateVKChannelRequest{
				ChannelName: "channel1",
				ChannelURL:  "https://vk.com/myserianet",
				ChannelType: "series",
			},
			statusCode: http.StatusUnprocessableEntity,
			response:   ErrorResponse{Error: "channel with name channel1 already exists"},
		},
		{
			name:       "invalid json data in request",
			input:      "1234",
			statusCode: http.StatusBadRequest,
			response:   ErrorResponse{Error: "invalid request body"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := &server{
				storage:  NewMockStorage(time.Now()),
				vkClient: vkclient.MustVKClient(),
			}

			body := new(bytes.Buffer)
			if err := json.NewEncoder(body).Encode(tc.input); err != nil {
				t.Fatal(err)
			}
			r, err := http.NewRequest("POST", "/api/v1/ship", body)
			if err != nil {
				t.Fatal(err)
			}
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			s.postShip(w, r)
			if w.Code != tc.statusCode {
				t.Errorf("Expected status code %d, got %d", tc.statusCode, w.Code)
			}
			if tc.statusCode >= 200 && tc.statusCode < 300 {
				data := IDResponse{}
				if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
					t.Fatal(err)
				}
				if !cmp.Equal(data, tc.response) {
					t.Errorf("handler returned unexpected body: got %v want %v", data, tc.response)
				}
			} else {
				data := ErrorResponse{}
				if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
					t.Fatal(err)
				}
				if !cmp.Equal(data, tc.response) {
					t.Errorf("handler returned unexpected body: got %v want %v", data, tc.response)
				}
			}
		})
	}
}

func TestDeleteShip(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		name       string
		input      string
		statusCode int
		response   any
	}{
		{
			name:       "valid id in url",
			input:      "1",
			statusCode: http.StatusOK,
			response:   nil,
		},
		{
			name:       "invalid id in url",
			input:      "a",
			statusCode: http.StatusBadRequest,
			response:   ErrorResponse{Error: "invalid id"},
		},
		{
			name:       "non-existent id in url",
			input:      "4",
			statusCode: http.StatusNotFound,
			response:   ErrorResponse{Error: "channel with id 4 not found"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := &server{
				storage: NewMockStorage(time.Now()),
			}

			r, err := http.NewRequest("DELETE", "/api/v1/ship", nil)
			if err != nil {
				t.Fatal(err)
			}
			vars := map[string]string{
				"id": tc.input,
			}
			r = mux.SetURLVars(r, vars)
			w := httptest.NewRecorder()
			s.deleteShip(w, r)
			if w.Code != tc.statusCode {
				t.Errorf("Expected status code %d, got %d", tc.statusCode, w.Code)
			}
			if tc.statusCode >= 200 && tc.statusCode < 300 {
				r.Method = "GET"
				w = httptest.NewRecorder()
				s.getShipByID(w, r)
				if w.Code != http.StatusNotFound {
					t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
				}
			} else {
				data := ErrorResponse{}
				if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
					t.Fatal(err)
				}
				if !cmp.Equal(data, tc.response) {
					t.Errorf("handler returned unexpected body: got %v want %v", data, tc.response)
				}
			}
		})
	}
}

func TestPatchShip(t *testing.T) {
	t.Parallel()

	now := time.Now()

	var testCases = []struct {
		name       string
		input      any
		id         string
		statusCode int
		response   any
	}{
		{
			name:       "invalid json data in request",
			input:      "1234",
			id:         "1",
			statusCode: http.StatusBadRequest,
			response:   ErrorResponse{Error: "invalid request body"},
		},
		{
			name: "invalid id in url",
			input: PatchVKChannelRequest{
				ChannelName: pointer("name"),
				ChannelType: pointer("type"),
				ChannelURL:  pointer("https://vk.com/name"),
				SiteURL:     pointer("https://site.com"),
			},
			id:         "aaa",
			statusCode: http.StatusBadRequest,
			response:   ErrorResponse{Error: "invalid id"},
		},
		{
			name: "non-existent id in url",
			input: PatchVKChannelRequest{
				ChannelName: pointer("name"),
				ChannelType: pointer("type"),
				ChannelURL:  pointer("https://vk.com/name"),
				SiteURL:     pointer("https://site.com"),
			},
			id:         "4",
			statusCode: http.StatusNotFound,
			response:   ErrorResponse{Error: "channel with id 4 not found"},
		},
		{
			name: "all fields in request filled",
			input: PatchVKChannelRequest{
				ChannelName: pointer("name"),
				ChannelType: pointer("type"),
				ChannelURL:  pointer("https://vk.com/name"),
				SiteURL:     pointer("https://site.com"),
			},
			id:         "1",
			statusCode: http.StatusOK,
			response: ChannelResponse{
				storage.VKChannel{
					Id:          1,
					ChannelName: pointer("name"),
					ChannelURL:  pointer("https://vk.com/name"),
					ChannelType: pointer("type"),
					SiteURL:     pointer("https://site.com"),
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		{
			name:       "filled only channel name",
			input:      PatchVKChannelRequest{ChannelName: pointer("name")},
			id:         "1",
			statusCode: http.StatusOK,
			response: ChannelResponse{
				storage.VKChannel{
					Id:          1,
					ChannelName: pointer("name"),
					ChannelURL:  pointer("url1"),
					ChannelType: pointer("type1"),
					SiteURL:     pointer("site1"),
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		{
			name:       "filled only channel type",
			input:      PatchVKChannelRequest{ChannelType: pointer("type")},
			id:         "1",
			statusCode: http.StatusOK,
			response: ChannelResponse{
				storage.VKChannel{
					Id:          1,
					ChannelName: pointer("channel1"),
					ChannelURL:  pointer("url1"),
					ChannelType: pointer("type"),
					SiteURL:     pointer("site1"),
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		{
			name:       "filled only channel url",
			input:      PatchVKChannelRequest{ChannelURL: pointer("https://vk.com/name")},
			id:         "1",
			statusCode: http.StatusOK,
			response: ChannelResponse{
				storage.VKChannel{
					Id:          1,
					ChannelName: pointer("channel1"),
					ChannelURL:  pointer("https://vk.com/name"),
					ChannelType: pointer("type1"),
					SiteURL:     pointer("site1"),
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		{
			name:       "invalid channel url in request",
			input:      PatchVKChannelRequest{ChannelURL: pointer("url")},
			id:         "1",
			statusCode: http.StatusUnprocessableEntity,
			response:   ErrorResponse{Error: "invalid channel url"},
		},
		{
			name:       "filled only site url",
			input:      PatchVKChannelRequest{SiteURL: pointer("https://site.com")},
			id:         "1",
			statusCode: http.StatusOK,
			response: ChannelResponse{
				storage.VKChannel{
					Id:          1,
					ChannelName: pointer("channel1"),
					ChannelURL:  pointer("url1"),
					ChannelType: pointer("type1"),
					SiteURL:     pointer("https://site.com"),
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		{
			name:       "invalid site url in request",
			input:      PatchVKChannelRequest{SiteURL: pointer("site")},
			id:         "1",
			statusCode: http.StatusUnprocessableEntity,
			response:   ErrorResponse{Error: "invalid site url"},
		},
		{
			name:       "empty json in request",
			input:      PatchVKChannelRequest{},
			id:         "1",
			statusCode: http.StatusUnprocessableEntity,
			response:   ErrorResponse{Error: "nothing to update"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := &server{
				storage: NewMockStorage(now),
			}

			body := new(bytes.Buffer)
			if err := json.NewEncoder(body).Encode(tc.input); err != nil {
				t.Fatal(err)
			}
			r, err := http.NewRequest("PATCH", "/api/v1/ship", body)
			if err != nil {
				t.Fatal(err)
			}
			vars := map[string]string{
				"id": tc.id,
			}
			r = mux.SetURLVars(r, vars)
			w := httptest.NewRecorder()
			s.patchShip(w, r)
			if w.Code != tc.statusCode {
				t.Errorf("Expected status code %d, got %d", tc.statusCode, w.Code)
			}
			if tc.statusCode >= 200 && tc.statusCode < 300 {
				r.Method = "GET"
				w = httptest.NewRecorder()
				s.getShipByID(w, r)
				data := ChannelResponse{}
				if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
					t.Fatal(err)
				}
				if !cmp.Equal(data, tc.response, cmpopts.IgnoreFields(storage.VKChannel{}, "UpdatedAt")) && data.UpdatedAt == now {
					t.Errorf("handler returned unexpected body: got %v want %v", data, tc.response)
				}
			} else {
				data := ErrorResponse{}
				if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
					t.Fatal(err)
				}
				if !cmp.Equal(data, tc.response) {
					t.Errorf("handler returned unexpected body: got %v want %v", data, tc.response)
				}
			}
		})
	}
}

// MockStorage is a mock implementation of the storage interface
type MockStorage struct {
	VKChannels []storage.VKChannel
	indexId    int64
}

func NewMockStorage(time time.Time) *MockStorage {
	vkChannels := []storage.VKChannel{
		{Id: 1, ChannelName: pointer("channel1"), ChannelType: pointer("type1"), ChannelURL: pointer("url1"), SiteURL: pointer("site1"), CreatedAt: time, UpdatedAt: time},
		{Id: 2, ChannelName: pointer("channel2"), ChannelType: pointer("type2"), ChannelURL: pointer("url2"), SiteURL: pointer("site2"), CreatedAt: time, UpdatedAt: time},
		{Id: 3, ChannelName: pointer("channel3"), ChannelType: pointer("type3"), ChannelURL: pointer("url3"), SiteURL: pointer("site3"), CreatedAt: time, UpdatedAt: time},
	}

	return &MockStorage{
		VKChannels: vkChannels,
		indexId:    3,
	}
}

func pointer(s string) *string {
	return &s
}

func (m *MockStorage) GetVKChannels(ctx context.Context) ([]storage.VKChannel, error) {
	return m.VKChannels, nil
}

func (m *MockStorage) GetVKChannel(ctx context.Context, id int64) (storage.VKChannel, error) {
	for _, channel := range m.VKChannels {
		if channel.Id == id {
			return channel, nil
		}
	}

	return storage.VKChannel{}, fmt.Errorf("channel with id %d not found", id)
}

func (m *MockStorage) GetVKChannelByName(ctx context.Context, name string) (storage.VKChannel, error) {
	for _, channel := range m.VKChannels {
		if *channel.ChannelName == name {
			return channel, nil
		}
	}

	return storage.VKChannel{}, fmt.Errorf("channel with name %v not found", name)
}

func (m *MockStorage) GetVKChannelsByType(ctx context.Context, channelType string) ([]storage.VKChannel, error) {
	channels := make([]storage.VKChannel, 0)

	for _, channel := range m.VKChannels {
		if *channel.ChannelType == channelType {
			channels = append(channels, channel)
		}
	}

	if len(channels) == 0 {
		return []storage.VKChannel{}, fmt.Errorf("no channels found for type %v", channelType)
	}

	return channels, nil
}

func (m *MockStorage) UpdateVKChannel(ctx context.Context, channel storage.VKChannel) error {
	for _, c := range m.VKChannels {
		if c.Id == channel.Id {
			if channel.ChannelName != nil {
				*c.ChannelName = *channel.ChannelName
			}
			if channel.ChannelType != nil {
				*c.ChannelType = *channel.ChannelType
			}
			if channel.ChannelURL != nil {
				*c.ChannelURL = *channel.ChannelURL
			}
			if channel.SiteURL != nil {
				*c.SiteURL = *channel.SiteURL
			}
			c.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("channel with id %d not found", channel.Id)
}

func (m *MockStorage) DeleteVKChannel(ctx context.Context, id int64) error {
	for i, channel := range m.VKChannels {
		if channel.Id == id {
			m.VKChannels = append(m.VKChannels[:i], m.VKChannels[i+1:]...)
			m.indexId--
			return nil
		}
	}

	return fmt.Errorf("channel with id %d not found", id)
}

func (m *MockStorage) CreateVKChannel(ctx context.Context, channel storage.VKChannel) (int64, error) {
	m.indexId++
	m.VKChannels = append(m.VKChannels, storage.VKChannel{
		Id:          m.indexId,
		ChannelName: channel.ChannelName,
		ChannelType: channel.ChannelType,
		ChannelURL:  channel.ChannelURL,
		SiteURL:     channel.SiteURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	return m.indexId, nil
}

func (m *MockStorage) Disconnect() {}

func (m *MockStorage) Stat() *pgxpool.Stat {
	return nil
}

func (m *MockStorage) Ping() error {
	return nil
}

func (m *MockStorage) Reconnect(ctx context.Context) error {
	return nil
}
