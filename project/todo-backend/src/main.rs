use axum::{
    extract::{State, Json},
    routing::{get, post},
    Router,
    response::IntoResponse,
};
use serde::{Deserialize, Serialize};
use std::{sync::{Arc, Mutex}};
use tokio::net::TcpListener;
use tower_http::cors::{CorsLayer, Any};

type TodoList = Arc<Mutex<Vec<String>>>;

#[derive(Deserialize)]
struct NewTodo {
    text: String,
}

#[derive(Serialize)]
struct TodosResponse {
    todos: Vec<String>,
}

async fn get_todos(State(todos): State<TodoList>) -> impl IntoResponse {
    let todos = todos.lock().unwrap();
    axum::Json(TodosResponse {
        todos: todos.clone(),
    })
}

async fn add_todo(
    State(todos): State<TodoList>,
    Json(payload): Json<NewTodo>,
) -> impl IntoResponse {
    let mut todos = todos.lock().unwrap();
    todos.push(payload.text);
    "Todo added"
}

#[tokio::main]
async fn main() {
    let todos: TodoList = Arc::new(Mutex::new(Vec::new()));

    let app = Router::new()
        .route("/todos", get(get_todos).post(add_todo))
        .with_state(todos)
        .layer(
            CorsLayer::new()
                .allow_origin(Any)
                .allow_methods(Any)
                .allow_headers(Any)
        );

    let port = std::env::var("PORT").expect("Environment variable PORT is required");
    let addr = format!("0.0.0.0:{}", port);
    let listener = TcpListener::bind(&addr).await.unwrap();

    println!("Server started in port {}", port);

    axum::serve(listener, app).await.unwrap();
}
