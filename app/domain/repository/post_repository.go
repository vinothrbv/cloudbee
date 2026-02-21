package repository

import (
	"github.com/vinothrbv/cloudbee/app/domain/entity"
	"errors"
	"sync"
)

type PostRepository interface {
	Create(post *entity.Post) error
	Get(id int64) (*entity.Post, error)
	Update(post *entity.Post) error
	Delete(id int64) error
}

type InMemoryPostRepository struct {
	mu     sync.Mutex
	store  map[int64]*entity.Post
	nextID int64
}

func NewInMemoryPostRepository() *InMemoryPostRepository {
	return &InMemoryPostRepository{
		store:  make(map[int64]*entity.Post),
		nextID: 1,
	}
}

func (r *InMemoryPostRepository) Create(post *entity.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if post.ID == 0 {
		post.ID = r.nextID
		r.nextID++
	} else {
		if _, ok := r.store[post.ID]; ok {
			return errors.New("post already exists")
		}
		if post.ID >= r.nextID {
			r.nextID = post.ID + 1
		}
	}
	r.store[post.ID] = post
	return nil
}

func (r *InMemoryPostRepository) Get(id int64) (*entity.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	post, ok := r.store[id]
	if !ok {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (r *InMemoryPostRepository) Update(post *entity.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.store[post.ID]
	if !ok {
		return errors.New("post not found")
	}
	r.store[post.ID] = post
	return nil

}

func (r *InMemoryPostRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.store[id]
	if !ok {
		// idempotent operation
		return nil
	}
	delete(r.store, id)
	return nil
}
