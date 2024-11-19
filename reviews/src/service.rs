use sqlx::PgPool;
use tonic::{Request, Response, Status};
use uuid::Uuid;

use crate::reviews_proto::*;

pub struct ReviewServiceImpl {
    pool: PgPool,
}

impl ReviewServiceImpl {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }
}

#[tonic::async_trait]
impl review_service_server::ReviewService for ReviewServiceImpl {
    async fn add_review(
        &self,
        request: Request<AddReviewRequest>,
    ) -> Result<Response<AddReviewResponse>, Status> {
        let req = request.into_inner();

        let review = sqlx::query_as::<_, Review>(
            r#"
            INSERT INTO reviews (app_id, user_id, score, comment, tenant_id)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id, app_id, user_id, score, comment,
                      created_at, is_moderated, moderation_status, tenant_id
            "#)
            .bind(req.app_id)
            .bind(req.user_id)
        .bind(req.score as i32)
        .bind(req.comment)
        .bind(req.tenant_id)
        .fetch_one(&self.pool)
        .await
        .map_err(|e| Status::internal(format!("Failed to add review: {}", e)))?;

        Ok(Response::new(AddReviewResponse {
            review: Some(review),
            success: true,
            message: "Review added successfully".to_string(),
        }))
    }

    async fn get_reviews(
        &self,
        request: Request<GetReviewsRequest>,
    ) -> Result<Response<GetReviewsResponse>, Status> {
        let req = request.into_inner();
        let offset = (req.page * req.page_size) as i64;

        let reviews = sqlx::query_as::<_, Review>(
            r#"
            SELECT id, app_id, user_id, score, comment,
                   created_at, is_moderated, moderation_status, tenant_id
            FROM reviews
            WHERE app_id = $1 AND tenant_id = $2
            AND ($3 = false OR is_moderated = true)
            ORDER BY created_at DESC
            LIMIT $4 OFFSET $5
            "#,
        )
        .bind(&req.app_id)
        .bind(&req.tenant_id)
        .bind(&req.include_moderated_only)
        .bind(req.page_size as i64)
        .bind(offset)
        .fetch_all(&self.pool)
        .await
        .map_err(|e| Status::internal(format!("Failed to fetch reviews: {}", e)))?;

        let stats = sqlx::query!(
            r#"
            SELECT
                COUNT(*) as total_count,
                COALESCE(AVG(score::float8), 0.0) as average_score
            FROM reviews
            WHERE app_id = $1
            AND tenant_id = $2
            AND ($3 = false OR is_moderated = true)
            "#,
            req.app_id,
            req.tenant_id,
            req.include_moderated_only,
        )
        .fetch_one(&self.pool)
        .await
        .map_err(|e| Status::internal(format!("Failed to fetch review stats: {}", e)))?;

        let total_count = stats
            .total_count
            .ok_or(Status::internal("Failed to get total review count."))?;

        let average_score = stats
            .average_score
            .ok_or(Status::internal("Failed to get average score for reviews."))?;

        Ok(Response::new(GetReviewsResponse {
            reviews,
            total_count: total_count as u32,
            average_score,
        }))
    }

    async fn moderate_comment(
        &self,
        request: Request<ModerateCommentRequest>,
    ) -> Result<Response<ModerateCommentResponse>, Status> {
        let req = request.into_inner();

        let review_id = Uuid::parse_str(&req.review_id)
            .map_err(|e| Status::invalid_argument(format!("Invalid UUID: {}", e)))?;

        let review = sqlx::query_as::<_, Review>(
            r#"
            UPDATE reviews
            SET is_moderated = true,
                moderation_status = $1,
                moderator_id = $2,
                moderation_note = $3
            WHERE id = $4 AND tenant_id = $5
            RETURNING id, app_id, user_id, score, comment,
                      created_at, is_moderated, moderation_status
            "#,
        )
        .bind(req.moderation_status)
        .bind(req.moderator_id)
        .bind(req.moderation_note)
        .bind(review_id)
        .bind(req.tenant_id)
        .fetch_one(&self.pool)
        .await
        .map_err(|e| Status::internal(format!("Failed to moderate review: {}", e)))?;

        Ok(Response::new(ModerateCommentResponse {
            success: true,
            message: "Comment moderated successfully".to_string(),
            updated_review: Some(review),
        }))
    }
}
