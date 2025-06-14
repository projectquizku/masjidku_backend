CREATE TABLE advices (
  advice_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  advice_description TEXT NOT NULL,
  advice_lecture_id UUID REFERENCES lectures(lecture_id) ON DELETE SET NULL,
  advice_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  advice_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexing (disarankan untuk filter dan relasi)
CREATE INDEX idx_advices_user_id ON advices(advice_user_id);
CREATE INDEX idx_advices_lecture_id ON advices(advice_lecture_id);
CREATE INDEX idx_advices_created_at ON advices(advice_created_at);


CREATE TABLE articles (
  article_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  article_title VARCHAR(255) NOT NULL,
  article_description TEXT NOT NULL,
  article_image_url TEXT,
  article_order_id INT,
  article_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  article_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Untuk urutan tampilan artikel
CREATE INDEX idx_articles_order_id ON articles(article_order_id);

-- Untuk pencarian artikel berdasarkan waktu
CREATE INDEX idx_articles_created_at ON articles(article_created_at);
CREATE INDEX idx_articles_updated_at ON articles(article_updated_at);



CREATE TABLE quotes (
  quote_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  quote_text TEXT NOT NULL,
  is_published BOOLEAN DEFAULT FALSE,
  display_order INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performa filtering dan penampilan
CREATE INDEX idx_quotes_display_order ON quotes(display_order);
CREATE INDEX idx_quotes_created_at ON quotes(created_at);
CREATE INDEX idx_quotes_is_published ON quotes(is_published);