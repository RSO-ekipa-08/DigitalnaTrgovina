use crate::models::{Review, ReviewsResponse};
use crate::reviews_proto::*;
use crate::service::ReviewServiceImpl;
use juniper::{graphql_object, EmptySubscription, FieldError, FieldResult, RootNode};
use review_service_server::ReviewService;
use tonic::Request;

#[derive(Clone)]
pub struct Context {
    pub service: ReviewServiceImpl,
}

impl juniper::Context for Context {}

pub struct Query;

#[graphql_object(Context = Context)]
impl Query {
    async fn review(
        context: &Context,
        review_id: String,
        tenant_id: String,
    ) -> FieldResult<Review> {
        let request = Request::new(GetReviewRequest {
            review_id,
            tenant_id,
        });

        let response = context.service.get_review(request).await?;
        Ok(response
            .into_inner()
            .review
            .ok_or_else(|| <&str as Into<FieldError>>::into("Review not found"))
            .map(Into::into)?)
    }

    async fn reviews(
        context: &Context,
        app_id: String,
        page: i32,
        page_size: i32,
        include_moderated_only: bool,
        tenant_id: String,
    ) -> FieldResult<ReviewsResponse> {
        let request = Request::new(GetReviewsRequest {
            app_id,
            page: page as u32,
            page_size: page_size as u32,
            include_moderated_only,
            tenant_id,
        });

        let response = context.service.get_reviews(request).await?;
        Ok(response.into_inner().into())
    }
}

pub struct Mutation;

#[graphql_object(Context = Context)]
impl Mutation {
    async fn create_review(
        context: &Context,
        app_id: String,
        user_id: String,
        score: i32,
        comment: String,
        tenant_id: String,
    ) -> FieldResult<Review> {
        let request = Request::new(AddReviewRequest {
            app_id,
            user_id,
            score: score as u32,
            comment,
            tenant_id,
        });

        let response = context.service.add_review(request).await?;
        Ok(response
            .into_inner()
            .review
            .ok_or_else(|| <&str as Into<FieldError>>::into("Failed to create review"))
            .map(Into::into)?)
    }

    async fn update_review(
        context: &Context,
        review_id: String,
        score: i32,
        comment: String,
        tenant_id: String,
    ) -> FieldResult<Review> {
        let request = Request::new(UpdateReviewRequest {
            review_id,
            score: score as u32,
            comment,
            tenant_id,
        });

        let response = context.service.update_review(request).await?;
        Ok(response
            .into_inner()
            .review
            .ok_or_else(|| <&str as Into<FieldError>>::into("Failed to update review"))
            .map(Into::into)?)
    }

    async fn delete_review(
        context: &Context,
        review_id: String,
        tenant_id: String,
    ) -> FieldResult<bool> {
        let request = Request::new(DeleteReviewRequest {
            review_id,
            tenant_id,
        });

        let response = context.service.delete_review(request).await?;
        Ok(response.into_inner().success)
    }

    async fn moderate_review(
        context: &Context,
        review_id: String,
        moderation_status: i32,
        moderator_id: String,
        moderation_note: String,
        tenant_id: String,
    ) -> FieldResult<Review> {
        let request = Request::new(ModerateCommentRequest {
            review_id,
            moderation_status,
            moderator_id,
            moderation_note,
            tenant_id,
        });

        let response = context.service.moderate_comment(request).await?;
        Ok(response
            .into_inner()
            .updated_review
            .ok_or_else(|| <&str as Into<FieldError>>::into("Failed to moderate review"))
            .map(Into::into)?)
    }
}

pub type Schema = RootNode<'static, Query, Mutation, EmptySubscription<Context>>;

pub fn create_schema() -> Schema {
    Schema::new(Query, Mutation, EmptySubscription::new())
}
