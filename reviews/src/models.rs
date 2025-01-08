use crate::reviews_proto;
use juniper::GraphQLObject;

#[derive(GraphQLObject)]
#[graphql(description = "A review for an application")]
pub struct Review {
    pub id: String,
    pub app_id: String,
    pub user_id: String,
    pub score: i32,
    pub comment: String,
    pub created_at: String,
    pub is_moderated: bool,
    pub moderation_status: i32,
    pub tenant_id: String,
}

impl From<reviews_proto::Review> for Review {
    fn from(review: reviews_proto::Review) -> Self {
        Self {
            id: review.id,
            app_id: review.app_id,
            user_id: review.user_id,
            score: review.score as i32,
            comment: review.comment,
            created_at: review.created_at,
            is_moderated: review.is_moderated,
            moderation_status: review.moderation_status,
            tenant_id: review.tenant_id,
        }
    }
}

#[derive(GraphQLObject)]
#[graphql(description = "Response containing a list of reviews and statistics")]
pub struct ReviewsResponse {
    pub reviews: Vec<Review>,
    pub total_count: i32,
    pub average_score: f64,
}

impl From<reviews_proto::GetReviewsResponse> for ReviewsResponse {
    fn from(response: reviews_proto::GetReviewsResponse) -> Self {
        Self {
            reviews: response.reviews.into_iter().map(Into::into).collect(),
            total_count: response.total_count as i32,
            average_score: response.average_score,
        }
    }
}
