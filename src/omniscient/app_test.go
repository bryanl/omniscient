package omniscient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

type appTestFn func(u *url.URL, mnr *MockNoteRepository, h *Health)

func withApp(t *testing.T, healthOpts []HealthOption, fn appTestFn) {
	for _, c := range metricsCollectors {
		prometheus.Unregister(c)
	}

	health, err := NewHealth(healthOpts...)
	assert.NoError(t, err)

	mnr := &MockNoteRepository{}

	app, err := NewApp(
		AppNoteRepository(mnr),
		AppHealth(health))
	assert.NoError(t, err)

	ts := httptest.NewServer(app.Mux)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)

	fn(u, mnr, app.health)

	app.health.mu.Lock()
	err = app.health.Stop()
	app.health.mu.Unlock()

	assert.NoError(t, err)
}

func TestAppCreate(t *testing.T) {
	withApp(t, nil, func(u *url.URL, mnr *MockNoteRepository, h *Health) {
		newNote := &Note{
			ID: "1",
		}
		mnr.On("Create", "new note").Return(newNote, nil)

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(&createNoteReq{Content: "new note"})
		assert.NoError(t, err)

		u.Path = "/notes"

		res, err := http.Post(u.String(), "application/json", &buf)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)

		defer res.Body.Close()
		var note Note
		err = json.NewDecoder(res.Body).Decode(&note)
		assert.NoError(t, err)

		assert.Equal(t, "1", note.ID)
	})

}

func TestAppRetrieveSingleNote(t *testing.T) {
	withApp(t, nil, func(u *url.URL, mnr *MockNoteRepository, h *Health) {
		existingNote := &Note{
			ID: "1",
		}
		mnr.On("Retrieve", "1").Return(existingNote, nil)

		u.Path = "/notes/1"

		res, err := http.Get(u.String())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		defer res.Body.Close()
		var note Note
		err = json.NewDecoder(res.Body).Decode(&note)
		assert.NoError(t, err)

		assert.Equal(t, "1", note.ID)
	})
}

func TestAppRetrieveAllNotes(t *testing.T) {
	withApp(t, nil, func(u *url.URL, mnr *MockNoteRepository, h *Health) {
		existingNote := Note{
			ID: "1",
		}
		mnr.On("List").Return([]Note{existingNote}, nil)

		u.Path = "/notes"

		res, err := http.Get(u.String())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		defer res.Body.Close()
		var notes []Note
		err = json.NewDecoder(res.Body).Decode(&notes)
		assert.NoError(t, err)

		assert.Len(t, notes, 1)
	})
}

func TestAppUpdateNote(t *testing.T) {
	withApp(t, nil, func(u *url.URL, mnr *MockNoteRepository, h *Health) {
		existingNote := &Note{
			ID:      "1",
			Content: "new content",
		}
		mnr.On("Update", "1", "new content").Return(existingNote, nil)

		u.Path = "/notes/1"

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(&updateNoteReq{Content: "new content"})
		assert.NoError(t, err)

		req, err := http.NewRequest("PUT", u.String(), &buf)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		defer res.Body.Close()
		var note Note
		err = json.NewDecoder(res.Body).Decode(&note)
		assert.NoError(t, err)

		assert.Equal(t, "1", note.ID)
		assert.Equal(t, "new content", note.Content)
	})
}

func TestAppDeleteNote(t *testing.T) {
	withApp(t, nil, func(u *url.URL, mnr *MockNoteRepository, h *Health) {
		mnr.On("Delete", "1").Return(nil)

		u.Path = "/notes/1"

		req, err := http.NewRequest("DELETE", u.String(), nil)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
	})
}

func TestAppHealthCheck(t *testing.T) {
	withApp(t, nil, func(u *url.URL, mnr *MockNoteRepository, h *Health) {
		u.Path = "/healthz"

		res, err := http.Get(u.String())
		assert.NoError(t, err)

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "OK", string(body))
	})
}

func TestAppInfo(t *testing.T) {
	withApp(t, nil, func(u *url.URL, mnr *MockNoteRepository, h *Health) {
		u.Path = "/app/info"

		res, err := http.Get(u.String())
		assert.NoError(t, err)

		defer res.Body.Close()

		var ai appInfo
		err = json.NewDecoder(res.Body).Decode(&ai)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "dev", ai.Revision)
	})
}
