CREATE TABLE IF NOT EXISTS events (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_title VARCHAR(255) NOT NULL,
    event_description TEXT,
    event_start_time TIMESTAMP NOT NULL,
    event_end_time TIMESTAMP NOT NULL,
    event_location VARCHAR(255),
    event_is_registration_required BOOLEAN DEFAULT FALSE,
    event_capacity INT,
    event_image_url TEXT,
    event_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    event_masjid_id UUID NOT NULL REFERENCES masjids(masjid_id) ON DELETE CASCADE
);

-- Index untuk pencarian berdasarkan masjid
CREATE INDEX idx_event_masjid_id ON events(event_masjid_id);


CREATE TABLE IF NOT EXISTS user_event_registrations (
    user_event_registration_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_event_registration_event_id UUID NOT NULL REFERENCES events(event_id) ON DELETE CASCADE,
    user_event_registration_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_event_registration_status VARCHAR(50) DEFAULT 'registered',
    user_event_registration_registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_event_registration_event_id, user_event_registration_user_id)
);

-- Index untuk efisiensi pencarian
CREATE INDEX idx_user_event_event_id ON user_event_registrations(user_event_registration_event_id);
CREATE INDEX idx_user_event_user_id ON user_event_registrations(user_event_registration_user_id);
