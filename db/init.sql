CREATE TABLE IF NOT EXISTS cities (
    id BIGSERIAL PRIMARY KEY,
    city TEXT NOT NULL,
    region TEXT NOT NULL,
    longitude NUMERIC NOT NULL,
    latitude NUMERIC NOT NULL,
    reviews_number NUMERIC NOT NULL DEFAULT 0,
    mark NUMERIC NOT NULL DEFAULT 0
);

CREATE INDEX idx_cities_location
ON cities USING gist (point(latitude, longitude));

CREATE TABLE IF NOT EXISTS reviews (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    author_id BIGINT NOT NULL,
    creation_date TIMESTAMP NOT NULL,
    city_id BIGINT NOT NULL REFERENCES cities(id),
    season TEXT NOT NULL CHECK (season IN ('Зима', 'Весна', 'Лето', 'Осень')),
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
    pet TEXT,
    physically_challenged_flag BOOLEAN NOT NULL DEFAULT false,
    limited_mobility_flag BOOLEAN NOT NULL DEFAULT false,
    elderly_people_flag BOOLEAN NOT NULL DEFAULT false,
    special_diet_flag BOOLEAN NOT NULL DEFAULT false,
    likes_number INTEGER NOT NULL DEFAULT 0,
    trip_type TEXT NOT NULL DEFAULT '',
    main_photo TEXT,
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
SELECT
    ct.city,
    TRIM(
        CASE
            WHEN ct.region_type = 'обл' THEN ct.region || ' область'
            WHEN ct.region_type = 'АО' THEN ct.region || ' автономный округ'
            WHEN ct.region_type = 'Аобл' THEN ct.region || ' автономная область'
            WHEN ct.region_type = 'Респ' THEN 'Республика ' || ct.region
            WHEN ct.region_type = 'край' THEN ct.region || ' край'
            ELSE ct.region
        END
        ),
    geo_lon,
    geo_lat
FROM cities_temp ct;

INSERT INTO reviews (author_id, creation_date, city_id, season, budget, tags, transport_mark, cleanliness_mark, preservation_mark, safety_mark,
                     hospitality_mark, price_quality_ratio, review_mark,
                     with_kids_flag, with_pets_flag, pet,
                     physically_challenged_flag, limited_mobility_flag,
                     elderly_people_flag, special_diet_flag,
                     trip_type, main_photo, status, review_content)
VALUES (1, NOW(), 536, 'Лето', 20000,
        '["Природа", "Спорт", "Поездка с животными"]'::jsonb, 5,
        5, 5, 5, 5, 5,
        5, false, false, NULL, true,
        true, true, true, 'Активная',
        'test/colomna1.jpg', 'published', '[{"text": "Коломна летом — это прекрасное место для семейного отдыха, сочетающее в себе историческую атмосферу и современные развлечения. Старинные улочки и архитектура создают неповторимый колорит, а разнообразие мероприятий и активностей позволяет найти занятие по душе как детям, так и взрослым.", "title": "Общее", "photos": [], "places": [{"name": "Исторический центр Коломны", "latitude": 55.102, "longitude": 38.783}, {"name": "Парк Коломенское", "latitude": 55.105, "longitude": 38.79}]}]'::jsonb);


INSERT INTO reviews (author_id, creation_date, city_id, season, budget, tags, transport_mark, cleanliness_mark, preservation_mark, safety_mark, hospitality_mark, price_quality_ratio, review_mark, with_kids_flag, with_pets_flag, pet, physically_challenged_flag, limited_mobility_flag, elderly_people_flag, special_diet_flag, trip_type, main_photo, status, review_content)
VALUES (3, NOW(), 929, 'Зима', 45000, '["Музеи", "Еда", "Поездка с детьми"]'::jsonb, 5, 5, 4, 5, 5, 4, 4.7, true, false, NULL, false, false, true, false, 'Семейная', 'kazan_winter.jpg', 'published',
        '[
          {"title": "Общее", "text": "Зимняя Казань встретила нас сказочной иллюминацией. Город очень чистый, а набережная идеальна для прогулок даже в мороз.", "photos": ["kzn1.jpg"], "places": []},
          {"title": "Еда", "text": "Попробовали настоящий эчпочмак и чак-чак. Рекомендуем семейные кафе на улице Баумана.", "photos": ["kzn_food.jpg"], "places": [{"name": "Дом Чая", "latitude": 55.792, "longitude": 49.11}]},
          {"title": "Проживание", "text": "Снимали апартаменты в центре, очень тепло и уютно.", "photos": [], "places": []},
          {"title": "Достопримечательности", "text": "Кремль обязателен к посещению. Мечеть Кул-Шариф на фоне снега выглядит невероятно.", "photos": ["kzn_kremlin.jpg"], "places": [{"name": "Кремль", "latitude": 55.799, "longitude": 49.10}]},
          {"title": "Особенности", "text": "Берите с собой очень теплую обувь, у Волги сильный ветер.", "photos": [], "places": []}
        ]'::jsonb);

INSERT INTO reviews (author_id, creation_date, city_id, season, budget, tags, transport_mark, cleanliness_mark, preservation_mark, safety_mark, hospitality_mark, price_quality_ratio, review_mark, with_kids_flag, with_pets_flag, pet, physically_challenged_flag, limited_mobility_flag, elderly_people_flag, special_diet_flag, trip_type, main_photo, status, review_content)
VALUES (1, NOW(), 394, 'Осень', 60000, '["Природа", "Здоровье и СПА", "Романтическая"]'::jsonb, 4, 4, 5, 5, 4, 3, 4.2, false, false, NULL, false, false, false, true, 'Романтическая', 'sochi_autumn.jpg', 'published',
        '[
          {"title": "Общее", "text": "Осень в Сочи — это лучшее время. Людей меньше, а море еще теплое.", "photos": ["sochi1.jpg"], "places": []},
          {"title": "Еда", "text": "Много ресторанов с морепродуктами, но цены кусаются. Искали места с веганским меню — их достаточно.", "photos": ["fish.jpg"], "places": []},
          {"title": "Проживание", "text": "Отель в Красной поляне с видом на горы. Спа-комплекс выше всяких похвал.", "photos": [], "places": []},
          {"title": "Достопримечательности", "text": "Ездили на Роза Хутор. Виды на облака сверху — это нечто.", "photos": ["mountains.jpg"], "places": [{"name": "Роза Пик", "latitude": 43.62, "longitude": 40.31}]},
          {"title": "Особенности", "text": "Осенью в горах погода меняется за 5 минут, берите дождевики.", "photos": [], "places": []}
        ]'::jsonb);

INSERT INTO reviews (author_id, creation_date, city_id, season, budget, tags, transport_mark, cleanliness_mark, preservation_mark, safety_mark, hospitality_mark, price_quality_ratio, review_mark, with_kids_flag, with_pets_flag, pet, physically_challenged_flag, limited_mobility_flag, elderly_people_flag, special_diet_flag, trip_type, main_photo, status, review_content)
VALUES (3, NOW(), 788, 'Лето', 35000, '["Музеи", "Достопримечательности", "Шоппинг"]'::jsonb, 5, 4, 5, 4, 4, 4, 4.5, false, true, 'Корги', false, false, false, false, 'Культурная', 'spb_summer.jpg', 'published',
        '[
          {"title": "Общее", "text": "Белые ночи — это то, что должен увидеть каждый. Город живет круглосуточно.", "photos": ["spb_night.jpg"], "places": []},
          {"title": "Еда", "text": "Корюшка и пышки — классика. Нашли дог-френдли кофейни почти на каждом углу.", "photos": ["pishki.jpg"], "places": [{"name": "Пышечная на Конюшенной", "latitude": 59.93, "longitude": 30.32}]},
          {"title": "Проживание", "text": "Жили в историческом здании с высокими потолками. Очень атмосферно.", "photos": [], "places": []},
          {"title": "Достопримечательности", "text": "Эрмитаж и Петергоф. В Петергофе пускают с маленькими собаками в сумках.", "photos": ["peterhof.jpg"], "places": [{"name": "Нижний парк", "latitude": 59.88, "longitude": 29.90}]},
          {"title": "Особенности", "text": "Развод мостов — красиво, но следите за временем, чтобы не остаться на другом берегу.", "photos": [], "places": []}
        ]'::jsonb);

INSERT INTO reviews (author_id, creation_date, city_id, season, budget, tags, transport_mark, cleanliness_mark, preservation_mark, safety_mark, hospitality_mark, price_quality_ratio, review_mark, with_kids_flag, with_pets_flag, pet, physically_challenged_flag, limited_mobility_flag, elderly_people_flag, special_diet_flag, trip_type, main_photo, status, review_content)
VALUES (3, NOW(), 120, 'Весна', 12000, '["Бюджетный отдых", "Паломничество", "История"]'::jsonb, 3, 5, 5, 5, 5, 5, 4.8, false, false, NULL, true, true, false, false, 'Паломническая', 'suzdal_spring.jpg', 'published',
        '[
          {"title": "Общее", "text": "Суздаль — город-музей. Очень тихо, спокойно, душа отдыхает.", "photos": ["suz1.jpg"], "places": []},
          {"title": "Еда", "text": "Медовуха и блины. Покупали у местных, очень вкусно и дешево.", "photos": [], "places": []},
          {"title": "Проживание", "text": "Гостевой дом на окраине. Чисто, по-домашнему.", "photos": [], "places": []},
          {"title": "Достопримечательности", "text": "Много старинных храмов. Город адаптирован для людей на колясках, почти везде есть пандусы.", "photos": ["suz_church.jpg"], "places": [{"name": "Кремль", "latitude": 56.41, "longitude": 40.44}]},
          {"title": "Особенности", "text": "Весной бывает много грязи на неасфальтированных тропинках.", "photos": [], "places": []}
        ]'::jsonb);

INSERT INTO reviews (author_id, creation_date, city_id, season, budget, tags, transport_mark, cleanliness_mark, preservation_mark, safety_mark, hospitality_mark, price_quality_ratio, review_mark, with_kids_flag, with_pets_flag, pet, physically_challenged_flag, limited_mobility_flag, elderly_people_flag, special_diet_flag, trip_type, main_photo, status, review_content)
VALUES (3, NOW(), 834, 'Лето', 25000, '["Спорт", "Кемпинг", "Ночная жизнь"]'::jsonb, 4, 4, 4, 4, 5, 5, 4.4, false, false, NULL, false, false, false, false, 'Молодежная', 'ekb_summer.jpg', 'published',
        '[
          {"title": "Общее", "text": "Столица Урала удивила своей энергией. Очень много молодежи и крутых баров.", "photos": ["ekb_view.jpg"], "places": []},
          {"title": "Еда", "text": "Уральские пельмени — это любовь. Ходили на гастрофестиваль.", "photos": ["pelmeni.jpg"], "places": []},
          {"title": "Проживание", "text": "Остановились в хостеле, а потом уехали с палатками за город.", "photos": [], "places": []},
          {"title": "Достопримечательности", "text": "Ельцин-центр — современно и технологично. Сходили на смотровую Высоцкого.", "photos": ["vysotsky.jpg"], "places": [{"name": "БЦ Высоцкий", "latitude": 56.83, "longitude": 60.61}]},
          {"title": "Особенности", "text": "В выходные в центре очень шумно, если хотите тишины — выбирайте другие районы.", "photos": [], "places": []}
        ]'::jsonb);

INSERT INTO reviews (author_id, creation_date, city_id, season, budget, tags, transport_mark, cleanliness_mark, preservation_mark, safety_mark, hospitality_mark, price_quality_ratio, review_mark, with_kids_flag, with_pets_flag, pet, physically_challenged_flag, limited_mobility_flag, elderly_people_flag, special_diet_flag, trip_type, main_photo, status, review_content)
VALUES (1, NOW(), 120, 'Зима', 35000, '["Поездка с детьми", "Достопримечательности", "Здоровье и СПА"]'::jsonb, 4, 5, 5, 5, 5, 4, 4.8, true, false, NULL, false, false, false, false, 'Семейная', 'test/suzdal_win.jpg', 'published',
        '[
          {"title": "Общее", "text": "Зимний Суздаль — это настоящая русская сказка. Снег здесь белый-белый, а воздух чистейший. Катались на санях, запряженных лошадьми.", "photos": ["suzdal_sani.jpg"], "places": []},
          {"title": "Еда", "text": "Обедали в ресторане при ГТК. Очень вкусные щи в хлебе и детское меню достойное.", "photos": ["suzdal_soup.jpg"], "places": []},
          {"title": "Проживание", "text": "Жили в деревянном срубе. Запах дерева и настоящая печь создали невероятный уют.", "photos": [], "places": []},
          {"title": "Достопримечательности", "text": "Музей деревянного зодчества зимой выглядит очень аутентично. Дети были в восторге от мастер-класса по кузнечному делу.", "photos": ["zodchestvo.jpg"], "places": [{"name": "Музей деревянного зодчества", "latitude": 56.411, "longitude": 40.440}]},
          {"title": "Особенности", "text": "В мороз очень быстро садятся телефоны, берите повербанки.", "photos": [], "places": []}
        ]'::jsonb);

INSERT INTO reviews (author_id, creation_date, city_id, season, budget, tags, transport_mark, cleanliness_mark, preservation_mark, safety_mark, hospitality_mark, price_quality_ratio, review_mark, with_kids_flag, with_pets_flag, pet, physically_challenged_flag, limited_mobility_flag, elderly_people_flag, special_diet_flag, trip_type, main_photo, status, review_content)
VALUES (1, NOW(), 120, 'Лето', 15000, '["Кемпинг", "Эко туризм", "Природа", "Спорт"]'::jsonb, 3, 5, 5, 5, 4, 5, 4.5, false, true, 'Джек-рассел', false, false, false, false, 'Молодежная', 'test/suzdal_eco.jpg', 'published',
        '[
          {"title": "Общее", "text": "Приехали с палатками и собакой. Нашли отличное место на берегу реки Каменка. Город очень dog-friendly!", "photos": ["suzdal_tent.jpg"], "places": []},
          {"title": "Еда", "text": "Закупались медовухой и овощами на местном рынке. Сами готовили на горелке, но заходили за кофе в центр — цены московские.", "photos": [], "places": [{"name": "Торговые ряды", "latitude": 56.420, "longitude": 40.448}]},
          {"title": "Проживание", "text": "Кемпинг-стоянка. Удобно, есть доступ к воде и электричеству.", "photos": [], "places": []},
          {"title": "Достопримечательности", "text": "Катались на сапах по реке Каменка. Вид на церкви с воды — это лучший ракурс, который можно придумать.", "photos": ["suzdal_sup.jpg"], "places": []},
          {"title": "Особенности", "text": "Летом очень много комаров у воды, берите мощные репелленты.", "photos": [], "places": []}
        ]'::jsonb);

INSERT INTO reviews (author_id, creation_date, city_id, season, budget, tags, transport_mark, cleanliness_mark, preservation_mark, safety_mark, hospitality_mark, price_quality_ratio, review_mark, with_kids_flag, with_pets_flag, pet, physically_challenged_flag, limited_mobility_flag, elderly_people_flag, special_diet_flag, trip_type, main_photo, status, review_content)
VALUES (3, NOW(), 120, 'Осень', 28000, '["Еда", "Музеи", "История", "Фестивали"]'::jsonb, 4, 4, 5, 5, 5, 3, 4.6, false, false, NULL, true, true, true, false, 'Культурная', 'test/suzdal_gast.jpg', 'published',
        '[
          {"title": "Общее", "text": "Осенний Суздаль в золоте — идеальное место для неспешных прогулок. Попали на праздник урожая, атмосфера непередаваемая.", "photos": ["suzdal_autumn.jpg"], "places": []},
          {"title": "Еда", "text": "Главное открытие — огуречное варенье! Также очень советуем попробовать запеченную утку в местных ресторанах.", "photos": ["suzdal_jam.jpg"], "places": []},
          {"title": "Проживание", "text": "Жили в гостевом доме в пяти минутах от Кремля. Идеально для пожилых родителей, всё рядом.", "photos": [], "places": []},
          {"title": "Достопримечательности", "text": "Суздальский кремль — мощь и история. Экспозиция внутри очень современная и интересная.", "photos": ["suzdal_kremlin.jpg"], "places": [{"name": "Суздальский Кремль", "latitude": 56.416, "longitude": 40.443}]},
          {"title": "Особенности", "text": "Город почти полностью адаптирован под маломобильных людей, но высокие пороги в старых зданиях всё еще встречаются.", "photos": [], "places": []}
        ]'::jsonb);