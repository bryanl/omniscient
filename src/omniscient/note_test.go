package omniscient

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRedisNoteRepoCreate(t *testing.T) {
	mrc := &MockRedisClient{}
	mrc.On("HMSet",
		"notes:1",
		"id", "1",
		mock.AnythingOfType("[]string")).Return("", nil)

	mrc.On("LPush", "notes:catalog", []string{"1"}).Return(int64(0), nil)

	id := 0
	igf := func() string {
		id++
		return fmt.Sprintf("%d", id)
	}

	rnr, err := NewRedisNoteRepository(
		RedisClientOption(mrc),
		NoteIDGenFn(igf))
	assert.NoError(t, err)

	n, err := rnr.Create("test")
	assert.Equal(t, "1", n.ID)
	assert.Equal(t, "test", n.Content)
	assert.NotEmpty(t, n.UpdatedAt)
	assert.NotEmpty(t, n.CreatedAt)
}

func TestRedisNoteRepoRetrieve(t *testing.T) {
	mrc := &MockRedisClient{}

	now := time.Now()

	m := map[string]string{
		fieldNoteID:        "1",
		fieldNoteContent:   "test",
		fieldNoteCreatedAt: now.Add(-30 * time.Minute).Format(time.RFC3339),
		fieldNoteUpdatedAt: now.Format(time.RFC3339),
	}
	mrc.On("HGetAllMap", "notes:1").Return(m, nil)

	rnr, err := NewRedisNoteRepository(
		RedisClientOption(mrc),
	)
	assert.NoError(t, err)

	n, err := rnr.Retrieve("1")
	assert.Equal(t, "1", n.ID)
	assert.Equal(t, "test", n.Content)
	assert.NotEmpty(t, n.UpdatedAt)
	assert.NotEmpty(t, n.CreatedAt)
}

func TestRedisNoteRepoUpdate(t *testing.T) {
	mrc := &MockRedisClient{}

	now := time.Now()
	updatedAt := now.Add(-15 * time.Minute)

	m := map[string]string{
		fieldNoteID:        "1",
		fieldNoteContent:   "test",
		fieldNoteCreatedAt: now.Add(-30 * time.Minute).Format(time.RFC3339),
		fieldNoteUpdatedAt: updatedAt.Format(time.RFC3339),
	}
	mrc.On("HGetAllMap", "notes:1").Return(m, nil)
	mrc.On("HMSet",
		"notes:1",
		"id", "1",
		mock.AnythingOfType("[]string")).Return("", nil)

	rnr, err := NewRedisNoteRepository(
		RedisClientOption(mrc),
	)
	assert.NoError(t, err)

	n, err := rnr.Update("1", "updated contents")
	assert.Equal(t, "1", n.ID)
	assert.Equal(t, "updated contents", n.Content)
	assert.True(t, updatedAt.Before(n.UpdatedAt))
	assert.NotEmpty(t, n.CreatedAt)
}

func TestRedisNoteRepoDelete(t *testing.T) {
	mrc := &MockRedisClient{}

	mrc.On("LRem", "notes:catalog", int64(0), "1").Return(int64(1), nil)
	mrc.On("Delete", []string{"notes:1"}).Return(int64(0), nil)

	rnr, err := NewRedisNoteRepository(
		RedisClientOption(mrc),
	)
	assert.NoError(t, err)

	err = rnr.Delete("1")
	assert.NoError(t, err)
}

func TestRedisNoteRepoList(t *testing.T) {
	mrc := &MockRedisClient{}

	ids := []string{"1"}
	mrc.On("LRange", "notes:catalog", int64(0), int64(-1)).Return(ids, nil)

	now := time.Now()

	m := map[string]string{
		fieldNoteID:        "1",
		fieldNoteContent:   "test",
		fieldNoteCreatedAt: now.Add(-30 * time.Minute).Format(time.RFC3339),
		fieldNoteUpdatedAt: now.Format(time.RFC3339),
	}
	mrc.On("HGetAllMap", "notes:1").Return(m, nil)

	rnr, err := NewRedisNoteRepository(
		RedisClientOption(mrc),
	)
	assert.NoError(t, err)

	notes, err := rnr.List()
	assert.Len(t, notes, 1)
	assert.NoError(t, err)
}
