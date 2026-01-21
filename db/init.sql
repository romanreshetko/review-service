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
    with_kids_flag BOOLEAN NOT NULL DEFAULT false,
    with_pets_flag BOOLEAN NOT NULL DEFAULT false,
    pet TEXT,
    business_trip_flag BOOLEAN NOT NULL DEFAULT false,
    physically_challenged_flag NOT NULL DEFAULT false,
    likes_number INTEGER NOT NULL DEFAULT 0,
    trip_type TEXT,
    main_photo TEXT,
    status TEXT NOT NULL CHECK (status IN ('published', 'moderating', 'blocked', 'draft')),
    review_content JSONB NOT NULL
);