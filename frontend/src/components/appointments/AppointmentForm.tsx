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

    // Date & time validation
    if (!form.date) {
      errs.date = "Date is required";
    }

    if (!form.time) {
      errs.time = "Time is required";
    }

    if (form.date && form.time) {
      const [hours, minutes] = form.time.split(":").map(Number);
      const [year, month, day] = form.date.split("-").map(Number);

      // Create local date/time
      const selectedDateTime = new Date(year, month - 1, day, hours, minutes, 0);

      const now = new Date();
  
      if (selectedDateTime < now) {
        errs.date = "Date and time must be in the future";
      }
   }

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
