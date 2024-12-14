-- Create categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create applications table
CREATE TABLE applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    developer_id UUID NOT NULL,
    category TEXT NOT NULL REFERENCES categories(name),
    price NUMERIC(10,2) NOT NULL DEFAULT 0,
    size BIGINT NOT NULL,
    min_android_version TEXT NOT NULL,
    current_version TEXT NOT NULL,
    tags TEXT[] NOT NULL DEFAULT '{}',
    screenshots TEXT[] NOT NULL DEFAULT '{}',
    storage_url TEXT NOT NULL,
    rating NUMERIC(3,2) NOT NULL DEFAULT 0,
    downloads INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    search_vector tsvector
);

-- Create downloads table
CREATE TABLE downloads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    application_id UUID NOT NULL REFERENCES applications(id),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    ip_address TEXT NOT NULL,
    success BOOLEAN NOT NULL DEFAULT true,
    error_message TEXT
);

-- Create indexes
CREATE INDEX applications_developer_id_idx ON applications(developer_id);
CREATE INDEX applications_category_idx ON applications(category);
CREATE INDEX applications_price_idx ON applications(price);
CREATE INDEX applications_rating_idx ON applications(rating);
CREATE INDEX applications_downloads_idx ON applications(downloads);
CREATE INDEX downloads_user_id_idx ON downloads(user_id);
CREATE INDEX downloads_application_id_idx ON downloads(application_id);
CREATE INDEX downloads_timestamp_idx ON downloads(timestamp);
CREATE INDEX applications_search_idx ON applications USING GIN(search_vector);

-- Create trigger for updating updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for updating search_vector
CREATE OR REPLACE FUNCTION applications_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.description, '')), 'B') ||
        setweight(to_tsvector('english', array_to_string(NEW.tags, ' ')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_applications_search_vector
    BEFORE INSERT OR UPDATE ON applications
    FOR EACH ROW
    EXECUTE FUNCTION applications_search_vector_update();

CREATE TRIGGER update_applications_updated_at
    BEFORE UPDATE ON applications
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default categories
INSERT INTO categories (name, description) VALUES
    ('Games', 'Mobile games for entertainment'),
    ('Productivity', 'Apps to help you get things done'),
    ('Education', 'Learning and educational apps'),
    ('Social', 'Social networking and communication apps'),
    ('Entertainment', 'Entertainment and media apps'),
    ('Health & Fitness', 'Health tracking and fitness apps'),
    ('Business', 'Business and professional apps'),
    ('Tools', 'Utility and tool apps'),
    ('Travel', 'Travel and navigation apps'),
    ('Shopping', 'Shopping and e-commerce apps')
ON CONFLICT (name) DO NOTHING; 