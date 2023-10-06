package models

import "database/sql"

// Task represents a task in the application.
type Task struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	UserID      uint       `json:"-"`
}

// TaskStatus represents the status of a task.
type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in progress"
	TaskStatusDone       TaskStatus = "done"
)

// TaskRepository provides an interface for task-related database operations.
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository creates a new instance of TaskRepository.
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// CreateTasks inserts a new task into the database.
func (r *TaskRepository) CreateTask(task *Task) (*Task, error) {
	row := r.db.QueryRow("INSERT INTO tasks (title, description, status, userID) VALUES ($1, $2, $3, $4) RETURNING *", task.Title, task.Description, task.Status, task.UserID)

	var newTask Task
	err := row.Scan(&newTask.ID, &newTask.Title, &newTask.Description, &newTask.Status, &newTask.UserID)
	if err != nil {
		return nil, err
	}

	return &newTask, nil
}

// GetTaskByID retrieves a task by its ID from the database.
func (r *TaskRepository) GetTaskByID(taskID uint) (*Task, error) {
	row := r.db.QueryRow("SELECT * FROM tasks WHERE id = $1", taskID)

	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.UserID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// UpdateTask updates a task in the database.
func (r *TaskRepository) UpdateTask(task *Task) (*Task, error) {
	row := r.db.QueryRow("UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4 RETURNING *", task.Title, task.Description, task.Status, task.ID)

	var updatedTask Task
	err := row.Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Description, &updatedTask.Status, &updatedTask.UserID)
	if err != nil {
		return nil, err
	}

	return &updatedTask, nil
}

// DeleteTask deletes a task from the database.
func (r *TaskRepository) DeleteTask(taskID uint) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	return err
}

// ListTasksByUserID retrieves a list of tasks belonging to a user in the database.
func (r *TaskRepository) ListTasksByUserID(userID uint) ([]Task, error) {
	rows, err := r.db.Query("SELECT * FROM tasks WHERE userID = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.UserID)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
