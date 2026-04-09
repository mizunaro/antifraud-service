CREATE TABLE IF NOT EXISTS url_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url TEXT NOT NULL UNIQUE,
    status SMALLINT NOT NULL DEFAULT 0, -- 0: Pending, 1: Safe, 2: Malicious
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_url_checks_status ON url_checks (status);
