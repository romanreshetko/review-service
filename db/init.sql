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
    review_mark NUMERIC,
    with_kids_flag BOOLEAN NOT NULL DEFAULT false,
    with_pets_flag BOOLEAN NOT NULL DEFAULT false,
    pet TEXT NOT NULL DEFAULT '',
    physically_challenged_flag BOOLEAN NOT NULL DEFAULT false,
    limited_mobility_flag BOOLEAN NOT NULL DEFAULT false,
    elderly_people_flag BOOLEAN NOT NULL DEFAULT false,
    special_diet_flag BOOLEAN NOT NULL DEFAULT false,
    likes_number INTEGER NOT NULL DEFAULT 0,
    trip_type TEXT NOT NULL DEFAULT '',
    main_photo TEXT NOT NULL DEFAULT 'default',
    status TEXT NOT NULL CHECK (status IN ('published', 'moderating', 'blocked', 'draft', 'reported', 'blocked_reported', 'undefined', 'moderation_error')),
    review_content JSONB NOT NULL,
    review_tsv tsvector
);

CREATE OR REPLACE FUNCTION reviews_tsvector_update()
RETURNS trigger AS $$
BEGIN
    NEW.review_tsv :=
        to_tsvector(
        'russian',
        (
        SELECT string_agg(sec->>'text', ' ')
        FROM jsonb_array_elements(NEW.review_content) AS sec
        )
                   );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_reviews_tsv
BEFORE INSERT OR UPDATE OF review_content
ON reviews
FOR EACH ROW
EXECUTE FUNCTION reviews_tsvector_update();


CREATE TABLE IF NOT EXISTS review_likes (
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

CREATE INDEX idx_reviews_likes_number
ON reviews (likes_number DESC)
WHERE status = 'published';


CREATE TEMPORARY TABLE cities_temp (
    address TEXT,
    postal_code TEXT,
    country TEXT,
    federal_district TEXT,
    region_type TEXT,
    region TEXT,
    area_type TEXT,
    area TEXT,
    city_type TEXT,
    city TEXT,
    settlement_type TEXT,
    settlement TEXT,
    kladr_id TEXT,
    fias_id TEXT,
    fias_level INTEGER,
    capital_marker INTEGER,
    okato TEXT,
    oktmo TEXT,
    tax_office TEXT,
    timezone TEXT,
    geo_lat NUMERIC,
    geo_lon NUMERIC,
    population INTEGER,
    foundation_year INTEGER
);

COPY cities_temp FROM '/docker-entrypoint-initdb.d/city.csv'
WITH (FORMAT csv, HEADER true, DELIMITER ',');

INSERT INTO cities (city, region, longitude, latitude)
SELECT ct.city, ct.region || ' ' || ct.region_type, geo_lon, geo_lat
FROM cities_temp ct;
