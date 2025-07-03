-- 🎨 THEMES TABLE
CREATE TABLE IF NOT EXISTS themes (
    theme_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    theme_name VARCHAR NOT NULL,
    theme_type INT NOT NULL DEFAULT 1,           -- 1=color, 2=wallpaper, 3=mix
    theme_colors JSONB NOT NULL,                 -- JSON object of UI colors
    wallpapers JSONB DEFAULT '[]'::jsonb,        -- Array of wallpapers
    required_level INT NOT NULL,                 -- Minimum level to unlock
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 🔎 Indexing for themes
CREATE INDEX IF NOT EXISTS idx_theme_type ON themes(theme_type);



-- 👤 USER THEMES TABLE
CREATE TABLE IF NOT EXISTS user_themes (
    user_theme_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    theme_id UUID NOT NULL REFERENCES themes(theme_id) ON DELETE CASCADE,
    is_selected BOOLEAN DEFAULT FALSE,
    selected_wallpaper_tag VARCHAR,
    unlocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 🔎 Indexing for user_themes
CREATE INDEX IF NOT EXISTS idx_user_themes_user_id ON user_themes(user_id);
CREATE INDEX IF NOT EXISTS idx_user_themes_user_selected ON user_themes(user_id, is_selected);
CREATE INDEX IF NOT EXISTS idx_user_themes_theme_id ON user_themes(theme_id);
