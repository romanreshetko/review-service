CREATE TABLE IF NOT EXISTS cities (
    id BIGSERIAL PRIMARY KEY,
    city TEXT NOT NULL,
    region TEXT NOT NULL,
    longitude NUMERIC NOT NULL,
    latitude NUMERIC NOT NULL,
    reviews_number NUMERIC NOT NULL DEFAULT 0,
    mark NUMERIC NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS reviews (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    author_id BIGINT NOT NULL,
    creation_date TIMESTAMP NOT NULL,
    city_id BIGINT NOT NULL REFERENCES cities(id),
    season TEXT NOT NULL CHECK (season IN ('winter', 'spring', 'summer', 'autumn')),
    budget INTEGER NOT NULL,
    tags JSONB NOT NULL DEFAULT '[]',
    transport_mark INTEGER,
    cleanliness_mark INTEGER,
    preservation_mark INTEGER,
    safety_mark INTEGER,
    hospitality_mark INTEGER,
    price_quality_ratio INTEGER,
    review_mark INTEGER
    with_kids_flag BOOLEAN NOT NULL DEFAULT false,
    with_pets_flag BOOLEAN NOT NULL DEFAULT false,
    pet TEXT NOT NULL DEFAULT '',
    business_trip_flag BOOLEAN NOT NULL DEFAULT false,
    physically_challenged_flag BOOLEAN NOT NULL DEFAULT false,
    limited_mobility_flag BOOLEAN NOT NULL DEFAULT false,
    elderly_people_flag BOOLEAN NOT NULL DEFAULT false,
    special_diet_flag BOOLEAN NOT NULL DEFAULT false,
    likes_number INTEGER NOT NULL DEFAULT 0,
    trip_type TEXT NOT NULL DEFAULT '',
    main_photo TEXT NOT NULL DEFAULT 'default',
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

CREATE TABLE IF NOT EXISTS reviews_likes (
    user_id BIGINT NOT NULL,
    review_id BIGINT NOT NULL REFERENCES reviews(id),
    PRIMARY KEY (user_id, review_id)
);

CREATE INDEX idx_reviews_tsv
ON reviews USING GIN(review_tsv)
WHERE status = 'published';

CREATE INDEX idx_reviews_city_published
ON reviews (city_id)
WHERE status = 'published';

CREATE INDEX idx_reviews_city_rating
ON reviews (city_id, review_mark DESC)
WHERE status = 'published';

CREATE INDEX idx_reviews_tags
ON reviews USING GIN(tags)
WHERE status = 'published';

CREATE INDEX idx_review_likes_user
ON review_likes(user_id);
