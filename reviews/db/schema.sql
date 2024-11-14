CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    score INTEGER NOT NULL CHECK (score >= 1 AND score <= 5),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_moderated BOOLEAN DEFAULT FALSE,
    moderation_status INTEGER DEFAULT 0,
    moderator_id VARCHAR(255),
    moderation_note TEXT
);

CREATE INDEX idx_reviews_app_id ON reviews(app_id);
CREATE INDEX idx_reviews_user_id ON reviews(user_id);
