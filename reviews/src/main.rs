use actix_cors::Cors;
use actix_web::{web, App, HttpResponse, HttpServer};
use juniper::http::{graphiql::graphiql_source, GraphQLRequest};
use std::env;
use std::sync::Arc;
use tonic::transport::Server;

mod extension;
mod graphql;
mod models;
mod service;

mod reviews_proto {
    include!("gen/proto/reviews.rs");
}

use graphql::{create_schema, Context, Schema};
use reviews_proto::review_service_server::ReviewServiceServer;
use service::ReviewServiceImpl;

async fn graphql(
    schema: web::Data<Schema>,
    context: web::Data<Context>,
    request: web::Json<GraphQLRequest>,
) -> HttpResponse {
    let response = request.execute(&schema, &context).await;
    HttpResponse::Ok().json(response)
}

async fn graphiql() -> HttpResponse {
    HttpResponse::Ok()
        .content_type("text/html; charset=utf-8")
        .body(graphiql_source("/graphql", None))
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    dotenvy::dotenv().ok();

    let database_url = get_database_url().await;
    let pool = sqlx::PgPool::connect(&database_url)
        .await
        .expect("Failed to create pool");

    sqlx::migrate!("db/migrations")
        .run(&pool)
        .await
        .expect("Failed to migrate database");

    let service = ReviewServiceImpl::new(pool.clone());
    let grpc_service = service.clone();

    // GraphQL setup
    let schema = Arc::new(create_schema());
    let context = Context {
        service: service.clone(),
    };

    // Start servers
    let grpc_addr = "0.0.0.0:50051".parse()?;
    let http_addr = "0.0.0.0:8080";

    println!("gRPC server listening on {}", grpc_addr);
    println!("GraphQL endpoint: http://{}/graphql", http_addr);

    // Run both servers concurrently
    tokio::spawn(async move {
        Server::builder()
            .add_service(ReviewServiceServer::new(grpc_service))
            .serve(grpc_addr)
            .await
            .expect("Failed to start gRPC server");
    });

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(schema.clone()))
            .app_data(web::Data::new(context.clone()))
            .wrap(
                Cors::default()
                    .allow_any_origin()
                    .allow_any_method()
                    .allow_any_header(),
            )
            .service(web::resource("/graphql").route(web::post().to(graphql)))
            .service(web::resource("/graphiql").route(web::get().to(graphiql)))
    })
    .bind(http_addr)?
    .run()
    .await?;

    Ok(())
}

async fn get_database_url() -> String {
    env::var("DATABASE_URL").expect("DATABASE_URL must be set")
}
