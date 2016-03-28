package omniscient

import (
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

const (
	fieldNoteID        = "id"
	fieldNoteContent   = "content"
	fieldNoteCreatedAt = "created_at"
	fieldNoteUpdatedAt = "updated_at"

	catalogKey = "catalog"
)

var (
	defaultRedisClient RedisClient
	defaultIDGenFn     = func() string {
		return uuid.NewV4().String()
	}
)

// Note is note.
type Note struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func noteFromMap(m map[string]string) *Note {
	n := &Note{
		ID:      m[fieldNoteID],
		Content: m[fieldNoteContent],
	}

	if t, err := time.Parse(time.RFC3339, m[fieldNoteCreatedAt]); err == nil {
		n.CreatedAt = t
	}

	if t, err := time.Parse(time.RFC3339, m[fieldNoteUpdatedAt]); err == nil {
		n.UpdatedAt = t
	}

	return n
}

// NoteRepository is a repository for managing notes.
type NoteRepository interface {
	Create(content string) (*Note, error)
	Retrieve(id string) (*Note, error)
	Update(id, content string) (*Note, error)
	Delete(id string) error
	List() ([]Note, error)
}

type idGenFn func() string

// RedisNoteRepository is an implementation of NoteRepository backed by Redis.
type RedisNoteRepository struct {
	base        string
	redisClient RedisClient
	idGen       idGenFn
}

var _ NoteRepository = (*RedisNoteRepository)(nil)

// RedisNoteRepositoryOption is an option for configuring RedisNoteRepository.
type RedisNoteRepositoryOption func(*RedisNoteRepository) error

// NewRedisNoteRepository creates an instance of NoteRepository.
func NewRedisNoteRepository(opts ...RedisNoteRepositoryOption) (NoteRepository, error) {
	rnr := &RedisNoteRepository{
		base:        "notes",
		redisClient: defaultRedisClient,
		idGen:       defaultIDGenFn,
	}

	for _, opt := range opts {
		err := opt(rnr)
		if err != nil {
			return nil, err
		}
	}

	return rnr, nil
}

// RedisClientOption sets the Redis client for RedisNoteRepository.
func RedisClientOption(rc RedisClient) RedisNoteRepositoryOption {
	return func(rnr *RedisNoteRepository) error {
		rnr.redisClient = rc
		return nil
	}
}

// NoteIDGenFn sets the id generator for RedisNoteRepository.
func NoteIDGenFn(idGen func() string) RedisNoteRepositoryOption {
	return func(rnr *RedisNoteRepository) error {
		rnr.idGen = idGen
		return nil
	}
}

// Create creates a new note.
func (nr *RedisNoteRepository) Create(content string) (*Note, error) {
	now := time.Now()

	note := Note{
		ID:        nr.idGen(),
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := nr.save(&note)
	if err != nil {
		return nil, err
	}

	_, err = nr.redisClient.LPush(nr.keyForID(catalogKey),
		note.ID)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

// Retrieve retrieves an existing note.
func (nr *RedisNoteRepository) Retrieve(id string) (*Note, error) {
	return nr.load(id)
}

// Update updates an existing note.
func (nr *RedisNoteRepository) Update(id, content string) (*Note, error) {
	m, err := nr.redisClient.HGetAllMap(nr.keyForID(id))
	if err != nil {
		return nil, err
	}

	n := noteFromMap(m)
	n.Content = content
	n.UpdatedAt = time.Now()

	err = nr.save(n)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// Delete deletes an existing note.
func (nr *RedisNoteRepository) Delete(id string) error {
	_, err := nr.redisClient.LRem(nr.keyForID(catalogKey), 0, id)
	if err != nil {
		return err
	}

	_, err = nr.redisClient.Delete(nr.keyForID(id))

	return err
}

// List retrieves all the notes.
func (nr *RedisNoteRepository) List() ([]Note, error) {
	notes := []Note{}
	ids, err := nr.redisClient.LRange(nr.keyForID(catalogKey), 0, -1)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		note, err := nr.load(id)
		if err != nil {
			return nil, err
		}

		notes = append(notes, *note)
	}

	return notes, nil
}

func (nr *RedisNoteRepository) keyForID(id string) string {
	return strings.Join([]string{
		nr.base, id,
	}, ":")
}

func (nr *RedisNoteRepository) save(note *Note) error {
	key := nr.keyForID(note.ID)

	_, err := nr.redisClient.HMSet(key,
		fieldNoteID, note.ID,
		fieldNoteContent, note.Content,
		fieldNoteCreatedAt, note.CreatedAt.Format(time.RFC3339),
		fieldNoteUpdatedAt, note.UpdatedAt.Format(time.RFC3339))

	return err
}

func (nr *RedisNoteRepository) load(id string) (*Note, error) {
	m, err := nr.redisClient.HGetAllMap(nr.keyForID(id))
	if err != nil {
		return nil, err
	}

	return noteFromMap(m), nil
}
