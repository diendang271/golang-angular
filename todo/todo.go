package todo

import (
	"errors"
	"github.com/rs/xid"
	"sync"
)

var (
	list []Todo
	mtx  sync.Mutex
	once sync.Once
	lock sync.RWMutex
)

func Init() {
	once.Do(InitializeList)
}

func InitializeList() {
	list = []Todo{}
}

type Todo struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	Complete bool   `json:"complete"`
}

// Get retrieves all elements from the todo list
func Get() []Todo {
	return list
}

// Add will add a new todo based on a message
func Add(message string) string {
	t := newTodo(message)
	mtx.Lock()
	list = append(list, t)
	mtx.Unlock()
	return t.ID
}

// Delete will remove a Todo from the Todo list
func Delete(id string) error {
	location, err := findTodoLocation(id)
	if err != nil {
		return err
	}

	removeElementByLocation(location)
	return nil
}

func Complete(id string) error {
	location, err := findTodoLocation(id)
	if err != nil {
		return err
	}
	setTodoCompleteByLocation(location)
	return nil
}

func newTodo(msg string) Todo {
	return Todo{
		ID:       xid.New().String(),
		Message:  msg,
		Complete: false,
	}
}

func findTodoLocation(id string) (int, error) {
	lock.RLock()
	defer lock.RUnlock()

	for i, t := range list {
		if isMatchingID(t.ID, id) {
			return i, nil
		}
	}

	return 0, errors.New("could not find todo based on id")
}

func removeElementByLocation(i int) {
	mtx.Lock()
	list = append(list[:i], list[i+1:]...)
	mtx.Unlock()
}

func setTodoCompleteByLocation(location int) {
	mtx.Lock()
	list[location].Complete = true
	mtx.Unlock()
}

func isMatchingID(a string, b string) bool {
	return a == b
}
