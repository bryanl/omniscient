package omniscient

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

var (
	defaultNoteRepository NoteRepository
	defaultHealth         *Health

	revision string
)

func init() {
	dnr, err := NewRedisNoteRepository()
	if err != nil {
		panic(fmt.Sprintf("unable to create default note repository: %v", err))
	}

	defaultNoteRepository = dnr

	health, err := NewHealth()
	if err != nil {
		panic(fmt.Sprintf("unable to create default health check service: %v", err))
	}

	defaultHealth = health
}

// App is the application.
type App struct {
	Mux      *echo.Echo
	noteRepo NoteRepository
	health   *Health
}

// AppOption is an option for configuring App.
type AppOption func(*App) error

// NewApp creates an instance of App.
func NewApp(opts ...AppOption) (*App, error) {
	e := echo.New()
	a := &App{
		Mux:      e,
		noteRepo: defaultNoteRepository,
		health:   defaultHealth,
	}

	for _, opt := range opts {
		if err := opt(a); err != nil {
			return nil, err
		}
	}

	// routes
	e.Post("/notes", a.createNote())
	e.Get("/notes", a.retrieveNotes())
	e.Get("/notes/:id", a.retrieveNote())
	e.Put("/notes/:id", a.updateNote())
	e.Delete("/notes/:id", a.deleteNote())

	e.Get("/healthz", a.healthz())
	e.Get("/app/info", a.appInfo())

	if a.health == nil {
		return nil, errors.New("no health checker")
	}

	err := a.health.Start()
	if err != nil {
		return nil, fmt.Errorf("unable to start health checker: %v", err)
	}

	return a, nil
}

// AppNoteRepository sets the note repository for App.
func AppNoteRepository(nr NoteRepository) AppOption {
	return func(a *App) error {
		a.noteRepo = nr
		return nil
	}
}

// AppHealth sets the app health checker option.
func AppHealth(h *Health) AppOption {
	return func(a *App) error {
		a.health = h
		return nil
	}
}

type createNoteReq struct {
	Content string `json:"content"`
}

type updateNoteReq struct {
	Content string `json:"content"`
}

func (a *App) createNote() echo.HandlerFunc {
	return func(c *echo.Context) error {
		cnr := &createNoteReq{}
		if err := c.Bind(cnr); err != nil {
			return err
		}

		note, err := a.noteRepo.Create(cnr.Content)
		if err != nil {
			msg := map[string]interface{}{
				"error": "unable to create note",
			}
			return c.JSON(http.StatusInternalServerError, msg)
		}

		return c.JSON(http.StatusCreated, note)
	}
}

func (a *App) retrieveNote() echo.HandlerFunc {
	return func(c *echo.Context) error {
		id := c.Param("id")
		note, err := a.noteRepo.Retrieve(id)
		if err != nil {
			msg := map[string]interface{}{
				"error": "note not found",
			}
			return c.JSON(http.StatusNotFound, msg)
		}

		return c.JSON(http.StatusOK, note)
	}
}

func (a *App) retrieveNotes() echo.HandlerFunc {
	return func(c *echo.Context) error {
		notes, err := a.noteRepo.List()
		if err != nil {
			msg := map[string]interface{}{
				"error": "unable to retrieve notes",
			}
			return c.JSON(http.StatusInternalServerError, msg)
		}

		return c.JSON(http.StatusOK, notes)
	}
}

func (a *App) updateNote() echo.HandlerFunc {
	return func(c *echo.Context) error {
		id := c.Param("id")

		cnr := &createNoteReq{}
		if err := c.Bind(cnr); err != nil {
			return err
		}

		note, err := a.noteRepo.Update(id, cnr.Content)
		if err != nil {
			msg := map[string]interface{}{
				"error": "unable to update note",
			}
			return c.JSON(http.StatusNotFound, msg)
		}

		return c.JSON(http.StatusOK, note)
	}
}

func (a *App) deleteNote() echo.HandlerFunc {
	return func(c *echo.Context) error {
		id := c.Param("id")

		err := a.noteRepo.Delete(id)
		if err != nil {
			msg := map[string]interface{}{
				"error": "unable to delete note",
			}
			return c.JSON(http.StatusBadRequest, msg)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (a *App) healthz() echo.HandlerFunc {
	return func(c *echo.Context) error {
		if a.health.IsOK() {
			return c.String(http.StatusOK, "OK")
		}

		return c.NoContent(http.StatusInternalServerError)
	}
}

type appInfo struct {
	Revision string `json:"revision"`
}

func (a *App) appInfo() echo.HandlerFunc {
	return func(c *echo.Context) error {
		ai := appInfo{
			Revision: revision,
		}
		if ai.Revision == "" {
			ai.Revision = "dev"
		}

		return c.JSON(http.StatusOK, ai)
	}
}
