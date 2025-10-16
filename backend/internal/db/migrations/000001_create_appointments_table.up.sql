CREATE TABLE appointments (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Add index on time for efficient lookups
CREATE INDEX idx_appointments_time ON appointments(date);