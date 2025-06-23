CREATE TABLE IF NOT EXISTS user_certificates (
    user_cert_id SERIAL PRIMARY KEY,
    user_cert_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_cert_subcategory_id INTEGER NOT NULL REFERENCES subcategories(subcategory_id) ON DELETE CASCADE,

    user_cert_is_up_to_date BOOLEAN NOT NULL DEFAULT true,
    user_cert_slug_url TEXT UNIQUE NOT NULL,

    user_cert_issued_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS certificates (
    certificate_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    certificate_title VARCHAR(255) NOT NULL,
    certificate_description TEXT,
    certificate_template_url TEXT, -- Link ke file template (PDF/JPG)
    certificate_created_by UUID REFERENCES users(id) ON DELETE SET NULL, -- Admin/pengunggah
    certificate_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);