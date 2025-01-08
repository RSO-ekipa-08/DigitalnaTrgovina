use actix_web::{get, web, HttpResponse};
use sqlx::PgPool;

#[derive(Debug, serde::Serialize)]
pub struct HealthResponse {
    status: String,
    version: String,
    database: bool,
}

#[get("/health")]
pub async fn health_check(pool: web::Data<PgPool>) -> HttpResponse {
    let db_health = match sqlx::query("SELECT 1").execute(&**pool).await {
        Ok(_) => true,
        Err(e) => {
            log::error!("Database health check failed: {}", e);
            false
        }
    };

    let response = HealthResponse {
        status: "ok".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
        database: db_health,
    };

    if db_health {
        HttpResponse::Ok().json(response)
    } else {
        HttpResponse::ServiceUnavailable().json(response)
    }
}

#[get("/ready")]
pub async fn ready_check() -> HttpResponse {
    HttpResponse::Ok().body("ok")
}
