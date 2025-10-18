CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE appointments (
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  start_time TIMESTAMP WITH TIME ZONE NOT NULL,
  end_time TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  CONSTRAINT no_overlapping_appointments
    EXCLUDE USING GIST (
      tstzrange(start_time, end_time) WITH &&
    )
);