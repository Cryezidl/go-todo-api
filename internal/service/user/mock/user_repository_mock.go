package mock

import (
	"context"
	"testing"

	"github.com/Cryezidl/go-todo-api/internal/model"
	"github.com/google/uuid"
)

/*
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByUserId(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByUserEmail(ctx context.Context, email string) (*model.User, error)
	FindByUserName(ctx context.Context, name string) (*model.User, error)
	FindByUserEmailOrName(ctx context.Context, login string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
*/

type MockUserRepository struct {
	//Create метод
	CreateError     error
	LastCreateUser  *model.User
	CreateCalled    bool
	CreateCallCount int

	//FindByUserId метод
	FindByUserIdResult    *model.User
	FindByUserIdError     error
	LastFindByUserId      uuid.UUID
	FindByUserIdCalled    bool
	FindByUserIdCallCount int

	//FindByUserEmail метод
	FindByUserEmailResult    *model.User
	FindByUserEmailError     error
	LastFindByUserEmail      string
	FindByUserEmailCalled    bool
	FindByUserEmailCallCount int

	//FindByUserName метод
	FindByUserNameResult    *model.User
	FindByUserNameError     error
	LastFindByUserName      string
	FindByUserNameCalled    bool
	FindByUserNameCallCount int

	//FindByUserEmailOrName метод
	FindByUserEmailOrNameResult    *model.User
	FindByUserEmailOrNameError     error
	LastFindByUserEmailOrName      string
	FindByUserEmailOrNameCalled    bool
	FindByUserEmailOrNameCallCount int

	//Update метод
	UpdateError     error
	LastUpdateUser  *model.User
	UpdateCalled    bool
	UpdateCallCount int

	//Delete метод
	DeleteError     error
	LastDeleteId    uuid.UUID
	DeleteCalled    bool
	DeleteCallCount int
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	m.CreateCalled = true
	m.CreateCallCount++
	m.LastCreateUser = user
	return m.CreateError
}

func (m *MockUserRepository) FindByUserId(ctx context.Context, id uuid.UUID) (*model.User, error) {
	m.FindByUserIdCalled = true
	m.FindByUserIdCallCount++
	m.LastFindByUserId = id

	return m.FindByUserIdResult, m.FindByUserIdError
}
func (m *MockUserRepository) FindByUserEmail(ctx context.Context, email string) (*model.User, error) {
	m.FindByUserEmailCalled = true
	m.FindByUserEmailCallCount++
	m.LastFindByUserEmail = email

	return m.FindByUserEmailResult, m.FindByUserEmailError
}
func (m *MockUserRepository) FindByUserName(ctx context.Context, name string) (*model.User, error) {
	m.FindByUserNameCalled = true
	m.FindByUserNameCallCount++
	m.LastFindByUserName = name

	return m.FindByUserNameResult, m.FindByUserNameError
}
func (m *MockUserRepository) FindByUserEmailOrName(ctx context.Context, login string) (*model.User, error) {
	m.FindByUserEmailOrNameCalled = true
	m.FindByUserEmailOrNameCallCount++
	m.LastFindByUserEmailOrName = login

	return m.FindByUserEmailOrNameResult, m.FindByUserEmailOrNameError
}
func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	m.UpdateCalled = true
	m.UpdateCallCount++
	m.LastUpdateUser = user
	return m.UpdateError
}
func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.DeleteCalled = true
	m.DeleteCallCount++
	m.LastDeleteId = id
	return m.DeleteError
}

func (m *MockUserRepository) AssertCreate(t *testing.T, expectedCallCount int) {
	if m.CreateCallCount != expectedCallCount {
		t.Errorf("Create called %d times, want %d", m.CreateCallCount, expectedCallCount)
	}
}

func (m *MockUserRepository) AssertFindByUserId(t *testing.T, expectedCallCount int) {
	if m.FindByUserIdCallCount != expectedCallCount {
		t.Errorf("FindByUserId called %d times, want %d", m.FindByUserIdCallCount, expectedCallCount)
	}
}

func (m *MockUserRepository) AssertFindByUserEmail(t *testing.T, expectedCallCount int) {
	if m.FindByUserEmailCallCount != expectedCallCount {
		t.Errorf("FindByUserEmail called %d times, want %d", m.FindByUserEmailCallCount, expectedCallCount)
	}
}

func (m *MockUserRepository) AssertFindByUserName(t *testing.T, expectedCallCount int) {
	if m.FindByUserNameCallCount != expectedCallCount {
		t.Errorf("FindByUserName called %d times, want %d", m.FindByUserNameCallCount, expectedCallCount)
	}
}

func (m *MockUserRepository) AssertFindByUserEmailOrName(t *testing.T, expectedCallCount int) {
	if m.FindByUserEmailOrNameCallCount != expectedCallCount {
		t.Errorf("FindByUserEmailOrName called %d times, want %d", m.FindByUserEmailOrNameCallCount, expectedCallCount)
	}
}

func (m *MockUserRepository) AssertUpdate(t *testing.T, expectedCallCount int) {
	if m.UpdateCallCount != expectedCallCount {
		t.Errorf("Update called %d times, want %d", m.UpdateCallCount, expectedCallCount)
	}
}

func (m *MockUserRepository) AssertDelete(t *testing.T, expectedCallCount int) {
	if m.DeleteCallCount != expectedCallCount {
		t.Errorf("Delete called %d times, want %d", m.DeleteCallCount, expectedCallCount)
	}
}

func (m *MockUserRepository) Reset() {
	//Create метод
	m.CreateError = nil
	m.LastCreateUser = nil
	m.CreateCalled = false
	m.CreateCallCount = 0

	//FindByUserId метод
	m.FindByUserIdResult = nil
	m.FindByUserIdError = nil
	m.LastFindByUserId = uuid.Nil
	m.FindByUserIdCalled = false
	m.FindByUserIdCallCount = 0

	//FindByUserEmail метод
	m.FindByUserEmailResult = nil
	m.FindByUserEmailError = nil
	m.LastFindByUserEmail = ""
	m.FindByUserEmailCalled = false
	m.FindByUserEmailCallCount = 0

	//FindByUserName метод
	m.FindByUserNameResult = nil
	m.FindByUserNameError = nil
	m.LastFindByUserName = ""
	m.FindByUserNameCalled = false
	m.FindByUserNameCallCount = 0

	//FindByUserEmailOrName метод
	m.FindByUserEmailOrNameResult = nil
	m.FindByUserEmailOrNameError = nil
	m.LastFindByUserEmailOrName = ""
	m.FindByUserEmailOrNameCalled = false
	m.FindByUserEmailOrNameCallCount = 0

	//Update метод
	m.UpdateError = nil
	m.LastUpdateUser = nil
	m.UpdateCalled = false
	m.UpdateCallCount = 0

	//Delete метод
	m.DeleteError = nil
	m.LastDeleteId = uuid.Nil
	m.DeleteCalled = false
	m.DeleteCallCount = 0
}
