-- Content is the data object accessed by the Content Management function of
-- the Epileptic application component.
CREATE TABLE IF NOT EXISTS content (
    id         UUID PRIMARY KEY,
    title      TEXT NOT NULL,
    body       TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
