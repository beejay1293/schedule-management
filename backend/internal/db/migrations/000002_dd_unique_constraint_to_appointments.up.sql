ALTER TABLE appointments
ADD CONSTRAINT unique_appointment_time UNIQUE (date);
