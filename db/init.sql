CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS cities (
    id UUID PRIMARY KEY DEFAULT uuid_generate(),
    city TEXT NOT NULL,
    region TEXT NOT NULL,
    longitude NUMERIC NOT NULL,
    latitude NUMERIC NOT NULL,
    reviews_number NUMERIC DEFAULT 0,
    mark NUMERIC DEFAULT 0
);

CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author_id INTEGER,
    creation_date TIMESTAMP,
    city_id UUID NOT NULL REFERENCES cities(name),
    season TEXT NOT NULL CHECK (season IN ('winter', 'spring', 'summer', 'autumn')),
    budget INTEGER,
    tags TEXT,
    transport_mark INTEGER,
    cleanliness_mark INTEGER,
    preservation_mark INTEGER,
    safety_mark INTEGER,
    hospitality_mark INTEGER,
    price_quality_ratio INTEGER,
    review_mark INTEGER
    with_kids_flag BOOLEAN NOT NULL DEFAULT false,
    with_pets_flag BOOLEAN NOT NULL DEFAULT false,
    pet TEXT,
    business_trip_flag BOOLEAN NOT NULL DEFAULT false,
    physically_challenged_flag BOOLEAN NOT NULL DEFAULT false,
    limited_mobility_flag BOOLEAN NOT NULL DEFAULT false,
    elderly_people_flag BOOLEAN NOT NULL DEFAULT false,
    special_diet_flag BOOLEAN NOT NULL DEFAULT false,
    likes_number INTEGER NOT NULL DEFAULT 0,
    trip_type TEXT,
    main_photo TEXT,
    status TEXT NOT NULL CHECK (status IN ('published', 'moderating', 'blocked', 'draft')),
    review_content JSONB NOT NULL,
    review_tsv tsvector
    GENERATED ALWAYS AS (
        to_tsvector(
        'russian',
        (
        SELECT string_agg(sec->>'text', ' ')
        FROM jsonb_array_elements(review_content->'sections') AS sec
        )
                   )
                        ) STORED
);

CREATE INDEX idx_reviews_city_published
ON reviews (city_id)
WHERE status = 'published';

CREATE INDEX idx_reviews_city_rating
ON reviews (city_id, review_mark DESC)
WHERE status = 'published';

CREATE INDEX idx_reviews_tsv
ON reviews USING GIN(review_tsv)
WHERE status = 'published';