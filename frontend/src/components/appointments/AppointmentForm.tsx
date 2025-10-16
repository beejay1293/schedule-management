import React, { useState } from "react";
import { useAppointments } from "../../hooks/useAppointments";
import "./Appointments.css";

interface FormErrors {
  title?: string;
  date?: string;
  time?: string;
  conflict?: string;
}

export const AppointmentForm: React.FC = () => {
  const { createMutation } = useAppointments();
  const [form, setForm] = useState({ title: "", date: "", time: "" });
  const [errors, setErrors] = useState<FormErrors>({});

  const validate = () => {
  const errs: FormErrors = {};

  // Title validation
  if (!form.title.trim()) errs.title = "Title is required";

  // Date validation
  if (!form.date) {
    errs.date = "Date is required";
  } else {
    const selectedDate = new Date(form.date);
    const today = new Date();
    today.setHours(0, 0, 0, 0); // start of today
    if (selectedDate < today) {
      errs.date = "Date cannot be in the past";
    }
  }

  // Time validation
  if (!form.time) errs.time = "Time is required";

  return errs;
};

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setErrors({}); // reset

    const validationErrors = validate();
    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      return;
    }

    createMutation.mutate(form, {
      onError: (err: any) => {
        if (err?.message?.includes("conflict")) {
          setErrors({ conflict: "Another appointment exists at that time" });
        } else {
          setErrors({ conflict: err.message || "Failed to schedule appointment" });
        }
      },
      onSuccess: () => {
        setForm({ title: "", date: "", time: "" });
        setErrors({});
      },
    });
  };

  return (
    <form onSubmit={handleSubmit} className="appointment-form">
      <h2>Schedule Appointment</h2>

      {errors.conflict && <p className="error conflict">{errors.conflict}</p>}

      <input
        type="text"
        placeholder="Title"
        value={form.title}
        onChange={(e) => setForm({ ...form, title: e.target.value })}
      />
      {errors.title && <p className="error">{errors.title}</p>}

      <input
        type="date"
        value={form.date}
        onChange={(e) => setForm({ ...form, date: e.target.value })}
      />
      {errors.date && <p className="error">{errors.date}</p>}

      <input
        type="time"
        value={form.time}
        onChange={(e) => setForm({ ...form, time: e.target.value })}
      />
      {errors.time && <p className="error">{errors.time}</p>}

      <button type="submit" disabled={createMutation.isPending}>
        {createMutation.isPending ? "Creating..." : "Schedule"}
      </button>
    </form>
  );
};
