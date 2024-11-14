use std::env;
use tonic::transport::Server;

mod service;
mod reviews_proto {
    tonic::include_proto!("reviews");
}

use reviews_proto::review_service_server::ReviewServiceServer;
use service::ReviewServiceImpl;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Load environment variables
    dotenvy::dotenv().ok();

    // Database connection
    let database_url = get_database_url().await;

    let pool = sqlx::PgPool::connect(&database_url)
        .await
        .expect("Failed to create pool");

    // Run migrations
    sqlx::migrate!("db/migrations")
        .run(&pool)
        .await
        .expect("Failed to migrate database");

    // Create service
    let addr = "0.0.0.0:50051".parse().unwrap();
    let service = ReviewServiceImpl::new(pool);

    println!("Review service listening on {}", addr);

    // Start server
    Server::builder()
        .add_service(ReviewServiceServer::new(service))
        .serve(addr)
        .await?;

    Ok(())
}

async fn get_database_url() -> String {
    // In Kubernetes, we'll use environment variables
    let user = env::var("POSTGRES_USER").expect("POSTGRES_USER must be set");
    let password = env::var("POSTGRES_PASSWORD").expect("POSTGRES_PASSWORD must be set");
    let host = env::var("POSTGRES_HOST").unwrap_or_else(|_| "localhost".to_string());
    let db = env::var("POSTGRES_DB").unwrap_or_else(|_| "reviews_db".to_string());

    format!("postgres://{}:{}@{}/{}", user, password, host, db)
}
