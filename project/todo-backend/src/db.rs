use postgres::{Client, NoTls};
use std::error::Error;

pub struct TodoStore {
    client: Client,
}

impl TodoStore {
    pub fn new(conn: &str) -> Result<Self, Box<dyn Error>> {
        let client = Client::connect(conn, NoTls)?;
        Ok(Self { client })
    }

    pub fn init(&mut self) -> Result<(), Box<dyn Error>> {
        self.client.batch_execute(
            "
              CREATE TABLE IF NOT EXISTS todos (
                  id SERIAL PRIMARY KEY,
                  todo TEXT NOT NULL,
                  done BOOLEAN NOT NULL DEFAULT FALSE
              );
            ",
        )?;
        Ok(())
    }

    pub fn insert_todo(&mut self, todo: &str) -> Result<String, Box<dyn Error>> {
        let row = self.client.query_one(
            "INSERT INTO todos (todo) VALUES ($1) RETURNING todo",
            &[&todo],
        )?;
        let inserted: String = row.get(0);
        println!("Inserted todo: {}", inserted);
        Ok(inserted)
    }

    pub fn list_todos(&mut self) -> Result<Vec<(i32, String, bool)>, Box<dyn Error>> {
        let rows = self
            .client
            .query("SELECT id, todo, done FROM todos ORDER BY id", &[])?;
        let mut todos = Vec::with_capacity(rows.len());
        for row in rows {
            let id: i32 = row.get(0);
            let todo: String = row.get(1);
            let done: bool = row.get(2);
            todos.push((id, todo, done));
        }
        Ok(todos)
    }

    pub fn update_todo_status(&mut self, id: i32, done: bool) -> Result<String, Box<dyn Error>> {
        let query = "UPDATE todos SET done = $1 WHERE id = $2 RETURNING todo";
        let row = self.client.query_one(query, &[&done, &id])?;
        let updated: String = row.get(0);
        println!("Updated todo: '{}' to be done", updated);
        Ok(updated)
    }

    pub fn ping(&mut self) -> Result<(), String> {
        self.client.query("SELECT 1", &[])
            .map(|_| ())
            .map_err(|e| e.to_string())
    }
}
