use crate::reviews_proto::Review;
use sqlx::types::time::OffsetDateTime;
use sqlx::Row;
use uuid::Uuid;

impl<'r> sqlx::FromRow<'r, sqlx::postgres::PgRow> for Review {
    fn from_row(row: &'r sqlx::postgres::PgRow) -> Result<Self, sqlx::Error> {
        Ok(Review {
            id: row.try_get::<Uuid, _>("id")?.to_string(),
            app_id: row.try_get("app_id")?,
            tenant_id: row.try_get("tenant_id")?,
            user_id: row.try_get("user_id")?,
            score: row.try_get::<i32, _>("score")? as u32,
            comment: row
                .try_get::<Option<String>, _>("comment")?
                .unwrap_or_default(),
            created_at: row
                .try_get::<Option<OffsetDateTime>, _>("created_at")?
                .unwrap_or_else(|| OffsetDateTime::now_utc())
                .to_string(),
            is_moderated: row
                .try_get::<Option<bool>, _>("is_moderated")?
                .unwrap_or_default(),
            moderation_status: row
                .try_get::<Option<i32>, _>("moderation_status")?
                .unwrap_or(0),
        })
    }
}
