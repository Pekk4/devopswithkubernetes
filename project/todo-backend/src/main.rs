mod db;

use axum::{
    extract::{State, Json},
    routing::{get},
    Router,
    response::IntoResponse,
    http::StatusCode
};
use serde::{Deserialize, Serialize};
use std::{sync::{Arc, Mutex}};
use tokio::net::TcpListener;
use tokio::task::spawn_blocking;
use tower_http::cors::{CorsLayer, Any};
use db::TodoStore;

#[derive(Clone)]
struct Config {
    db_creds: String,
    port: u16,
}

#[derive(Clone)]
struct AppState {
    db: Arc<Mutex<TodoStore>>,
}

#[derive(Deserialize)]
struct NewTodo {
    text: String,
}

#[derive(Serialize)]
struct TodosResponse {
    todos: Vec<String>,
}

#[derive(Serialize)]
struct TodoResponse {
    todo: String,
}

#[derive(Serialize)]
struct ErrorResponse {
    error: String,
}

async fn get_todos(State(state): State<AppState>) -> impl IntoResponse {
    let db = state.db.clone();
    let res = spawn_blocking(move || {
        let mut guard = db.lock().unwrap();
        guard.list_todos().map_err(|e| format!("{}", e))
    })
    .await;

    match res {
        Ok(Ok(todos)) => {
            let todos: Vec<String> = todos.into_iter().map(|(_, todo)| todo).collect();
            Json(TodosResponse { todos }).into_response()
        }
        Ok(Err(e)) => (StatusCode::INTERNAL_SERVER_ERROR, format!("db error: {}", e)).into_response(),
        Err(_) => (StatusCode::INTERNAL_SERVER_ERROR, "join error").into_response(),
    }
}

async fn add_todo(State(state): State<AppState>, Json(payload): Json<NewTodo>) -> impl IntoResponse {
    let db = state.db.clone();
    let todo_text = payload.text.clone();

    if todo_text.chars().count() > 140 {
        println!("Todo text exceeds 140 characters: {}", todo_text);
        return (
            StatusCode::BAD_REQUEST,
            Json(ErrorResponse {
                error: "Todo text exceeds 140 characters".to_string(),
            }),
        )
        .into_response();
    }

    let res = spawn_blocking(move || {
        let mut guard = db.lock().unwrap();
        guard.insert_todo(&todo_text).map_err(|e| format!("{}", e))
    })
    .await;

    match res {
        Ok(Ok(inserted)) => (
            StatusCode::CREATED,
            Json(TodoResponse { todo: inserted })
        ).into_response(),
        Ok(Err(e)) => (StatusCode::INTERNAL_SERVER_ERROR, format!("db error: {}", e)).into_response(),
        Err(_) => (StatusCode::INTERNAL_SERVER_ERROR, "join error").into_response(),
    }
}

async fn health_check(State(state): State<AppState>) -> StatusCode {
    let db = state.db.clone();
    let res = spawn_blocking(move || {
        let mut guard = db.lock().unwrap();
        guard.ping()
    })
    .await;

    match res {
        Ok(Ok(_)) => StatusCode::OK,
        _ => StatusCode::INTERNAL_SERVER_ERROR,
    }
}

fn init() -> Config {
    let db_creds = std::env::var("PG_URL")
        .expect("variable PG_URL is not set");
    let port = std::env::var("PORT")
        .unwrap_or_else(|_| "3000".into())
        .parse()
        .expect("PORT must be a number");

    Config { db_creds, port }
}

#[tokio::main]
async fn main() {
    let config = init();

    let db = spawn_blocking(move || {
        let mut store = TodoStore::new(&config.db_creds).expect("db connect");
        store.init().expect("db init");
        Arc::new(Mutex::new(store))
    })
    .await
    .expect("db init join");

    let state = AppState { db };

    let app = Router::new()
        .route("/healthz", get(health_check))
        .route("/todos", get(get_todos).post(add_todo))
        .with_state(state)
        .layer(
            CorsLayer::new()
                .allow_origin(Any)
                .allow_methods(Any)
                .allow_headers(Any)
        );

    let addr = format!("0.0.0.0:{}", config.port);
    let listener = TcpListener::bind(&addr).await.unwrap();

    println!("Server started in port {}", config.port);

    axum::serve(listener, app).await.unwrap();
}
