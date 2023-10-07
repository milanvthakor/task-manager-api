# Task Manager API

The RESTful API that allow users to create, read, update, and delete tasks. Users can register, log in, and access their own tasks only when authenticated.

## Table of Contents
1. [Project Design](#project-design)
    1. [Code Structure](#code-structure)
    2. [Database](#database)
    3. [Authentication](#authentication)
2. [Set up and Run API locally](#set)
3. [API Documentation](#api-documentation)
    1. [Base URL](#base-url)
    2. [Endpoints](#endpoints)
        1. [Health Check](#health-check)
        2. [User Registration](#user-registration)
        3. [User Login](#user-login)
        4. [Get Tasks](#get-tasks)
        5. [Create Task](#create-task)
        6. [Get Task By ID](#get-task-by-id)
        7. [Delete Task By ID](#delete-task-by-id)
        8. [Update Task](#update-task)
        9. [Mark Tasks as Done](#mark-tasks-as-done)

## Project Design

### Code Structure

The project's code structure has been designed to ensure that the codebase remains organized, maintainable, and adaptable as the project evolves.

1. **Separation of Concerns**: The structure separates different concerns of the application into distinct directories, making it easy to identify and work on specific components such as authentication, database handling, models, and API endpoints.
2. **Modularity**: Each major component of the application (e.g., authentication, tasks) is encapsulated within its own directory, allowing for easy module-level testing and reusability.
3. **Scalability**: The structure is scalable, making it relatively straightforward to add new features or components without disrupting the existing codebase.
4. **Dependency Injection**: The "config" packages suggest the use of dependency injection, making components more testable.

### Database 

PostgreSQL is chosen as the database for this API due to its suitability for handling structured and relational data. The main reason for this choice is the need to establish relationships between entities, such as users and tasks, which are naturally represented in a relational structure. PostgreSQL provides robust support for relational data modeling and offers features like foreign keys and transactions, making it well-suited for ensuring data integrity and consistency.

### Authentication

JWT (JSON Web Tokens) is adopted for authentication in this API as it allows for the secure inclusion of user-related data, such as the userID and email, within the token itself. This design choice simplifies the authentication process by eliminating the need to query the database for user information during each request, enhancing performance and reducing database load. Additionally, JWT provides a stateless authentication mechanism, ensuring scalability and compatibility with modern RESTful API architectures.

## Set up and Run API locally

1. Open your command line terminal.
2. Navigate to the directory where you want to clone the project.
3. Clone the GitHub repository using its URL
4. Navigate to the project directory that you just cloned.
5. Install project dependencies using the below command:

    ```
    go mod tidy
    ```
6. Inside the project folder, create a new file named `.env`. This file should be located at the same level as the `go.mod` file. In the `.env` file, make sure to set all the required environment variables according to the project's configuration. You can refer to the `env.example` file as a reference for the environment variables that need to be defined in your `.env` file.
7. Start the API server:

    ```
    go run cmd/api/main.go
    ```
    The terminal should display the message `Server is running on <PORT>`

You can now access the API endpoints using a tool like Postman or via the `curl` command.

## API Documentation 

An overview of the endpoints, required parameters, response formats, and example requests/responses of the API.

### Base URL

```
http://localhost:<PORT>
```

### Endpoints

#### Health Check
- **URL**: `/health`
- **Method**: `GET`
- **Description**: This API endpoint allows users to check the health of the server to ensure it is running properly.
- **Example Request**:
    ```
    GET /health
    ```
- **Example Response**:
    ```
    Status Code: 200

    {
        "message": "Server is up and running"
    }
    ```

#### User Registration
- **URL**: `/api/register`
- **Method**: `POST`
- **Description**: This API endpoint allows users to register by providing their email and password.
- **Request Body**: The request body must be in JSON format and include the following fields:
    - `email` (string, required): The email address of the user.
    - `password` (string, required): The password for the user account.
- **Example Request**:
    ```
    POST /api/register
    Content-Type: application/json

    {
        "email": "alice@example.com",
        "password": "Password@123"
    }
    ```
- **Example Response**:
    ```
    Status Code: 201

    {
        "message": "User registered successfully"
    }
    ```

#### User Login
- **URL**: `/api/login`
- **Method**: `POST`
- **Description**: This API endpoint allows users to log in by providing their email and password.
- **Request Body**: The request body must be in JSON format and include the following fields:
    - `email` (string, required): The email address of the user.
    - `password` (string, required): The password for the user account.
- **Example Request**:
    ```
    POST /api/login
    Content-Type: application/json

    {
        "email": "alice@example.com",
        "password": "Password@123"
    }
    ```
- **Example Response**:
    ```
    Status Code: 200

    {
        "message": "Login successful",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFsaWNlQGV4YW1wbGUuY29tIiwiaWQiOiI2NGY5Y2U1N2I4ZDg4ZGEyZDdkNjAwN2UiLCJpYXQiOjE2OTQwOTI4OTQsImV4cCI6MTY5NDA5NjQ5NH0.dXUcXMMHJsWspD89oiG43LyrbUjNBwKMKPK0_kcs-ug"
    }    
    ```

#### Get Tasks
- **URL**: `/api/tasks`
- **Method**: `GET`
- **Description**: This API endpoint allows users to retrieve a list of tasks. 
- **Example Request**:
    ```
    GET /api/tasks
    ```
- **Example Response**:
    ```
    Status Code: 200

    [
        {
            "id": 3,
            "title": "Task #1",
            "description": "Description of the Task #1",
            "status": "done"
        },
        {
            "id": 2,
            "title": "Task #1",
            "description": "Description of the Task #1",
            "status": "done"
        },
        {
            "id": 1,
            "title": "Task #1",
            "description": "Description of the Task #1",
            "status": "done"
        }
    ]
    ```
#### Create Task
- **URL**: `/api/tasks`
- **Method**: `POST`
- **Description**: This API endpoint allows users to create a task.
- **Headers**:
    - `Authorization` (string, required): The `Authorization` header must be set with a valid authentication token obtained from the `/login` endpoint. Use the format `Authorization: Bearer <token>`.
- **Request Body**: The request body must be in JSON format and include the following fields:
    - `title` (string, required): The title of the task.
    - `description` (string, optional): The description of the task.
    - `status` (string, required): The status of the task. It can have one of the following values: "todo", "in progress", or "done".
- **Example Request**:
    ```
    POST /api/tasks
    Content-Type: application/json
    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFsaWNlQGV4YW1wbGUuY29tIiwiaWQiOiI2NGY5Y2U1N2I4ZDg4ZGEyZDdkNjAwN2UiLCJpYXQiOjE2OTQwOTI4OTQsImV4cCI6MTY5NDA5NjQ5NH0.dXUcXMMHJsWspD89oiG43LyrbUjNBwKMKPK0_kcs-ug

    {
        "title": "Task #1",
        "description": "Description of the Task #1",
        "status": "done"
    }
    ```
- **Example Response**:
    ```
    Status Code: 201

    {
        "message": "Task created successfully",
        "task": {
            "id": 6,
            "title": "Task #1",
            "description": "Description of the Task #1",
            "status": "done"
        }
    }
    ```

#### Get Task by ID
- **URL**: `/api/tasks/{id}`
- **Method**: `GET`
- **Description**: This API endpoint allows users to retrieve details of a task by providing its unique ID. The user is allowed to retrieve details of only his/her task.
- **Headers**:
    - `Authorization` (string, required): The `Authorization` header must be set with a valid authentication token obtained from the `/login` endpoint. Use the format `Authorization: Bearer <token>`.
- **Path Parameters**:
    - `id` (string, required): The unique ID of the task to retrieve.
- **Example Request**:
    ```
    GET /api/tasks/1
    ```
- **Example Response**:
    ```
    Status Code: 200

    {
        "id": 1,
        "title": "Task #1",
        "description": "Description of the Task #1",
        "status": "done"
    }
    ```
#### Delete Task by ID
- **URL**: `/api/tasks/{id}`
- **Method**: `DELETE`
- **Description**: This API endpoint allows users to delete a task by providing its unique ID. The user is allowed to delete only his/her task.
- **Headers**:
    - `Authorization` (string, required): The `Authorization` header must be set with a valid authentication token obtained from the `/login` endpoint. Use the format `Authorization: Bearer <token>`.
- **Path Parameters**:
    - `id` (string, required): The unique ID of the task to delete.
- **Example Request**:
    ```
    DELETE /api/tasks/1
    ```
- **Example Response**:
    ```
    Status Code: 200

    {
        "message": "Task deleted successfully"
    }
    ```

#### Update Task
- **URL**: `/api/tasks/{id}`
- **Method**: `PUT`
- **Description**: This API endpoint allows users to update task details by providing its unique ID. The user is allowed to update details of only his/her task.
- **Headers**:
    - `Authorization` (string, required): The `Authorization` header must be set with a valid authentication token obtained from the `/login` endpoint. Use the format `Authorization: Bearer <token>`.
- **Request Body**: The request body must be in JSON format and can include any of the following fields:
    - `title` (string, optional): The title of the task.
    - `description` (string, optional): The description of the task.
    - `status` (string, optional): The status of the task. It can have one of the following values: "todo", "in progress", or "done".
- **Example Request**:
    ```
    PUT /api/tasks/1
    Content-Type: application/json
    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFsaWNlQGV4YW1wbGUuY29tIiwiaWQiOiI2NGY5Y2U1N2I4ZDg4ZGEyZDdkNjAwN2UiLCJpYXQiOjE2OTQwOTI4OTQsImV4cCI6MTY5NDA5NjQ5NH0.dXUcXMMHJsWspD89oiG43LyrbUjNBwKMKPK0_kcs-ug

    {
        "title": "Task #1",
        "description": "Description of the Task #1",
        "status": "done"
    }
    ```
- **Example Response**:
    ```
    Status Code: 200

    {
        "message": "Task updated successfully",
        "task": {
            "id": 1,
            "title": "Task #1",
            "description": "Description of the Task #1",
            "status": "done"
        }
    }
    ```

#### Mark Tasks as Done
- **URL**: `/api/tasks/mark-done`
- **Method**: `PATCH`
- **Description**: This API endpoint allows users to mark the status of multiple tasks as "done" by providing their unique IDs. The user is allowed to update details of only his/her tasks.
- **Headers**:
    - `Authorization` (string, required): The `Authorization` header must be set with a valid authentication token obtained from the `/login` endpoint. Use the format `Authorization: Bearer <token>`.
- **Request Body**: The request body must be in JSON format and include the list of task ID(s).
- **Example Request**:
    ```
    PATCH /api/tasks/mark-done
    Content-Type: application/json
    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFsaWNlQGV4YW1wbGUuY29tIiwiaWQiOiI2NGY5Y2U1N2I4ZDg4ZGEyZDdkNjAwN2UiLCJpYXQiOjE2OTQwOTI4OTQsImV4cCI6MTY5NDA5NjQ5NH0.dXUcXMMHJsWspD89oiG43LyrbUjNBwKMKPK0_kcs-ug

    [1,2,3]
    ```
- **Example Response**:
    ```
    Status Code: 200

    [
        {
            "id": 3,
            "message": "Task marked as done successfully"
        },
        {
            "id": 1,
            "message": "Task marked as done successfully"
        },
        {
            "id": 2,
            "message": "Task marked as done successfully"
        }
    ]
    ```
