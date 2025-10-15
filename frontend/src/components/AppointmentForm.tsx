import { useState } from "react";
import { createAppointment } from "../api/appointmentApi";
import { useQueryClient } from "@tanstack/react-query";

export function AppointmentForm() {
  const queryClient = useQueryClient();
  const [form, setForm] = useState({ title: "", date: "", time: "" });
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createAppointment(form);
      queryClient.invalidateQueries(["appointments"]);
      setForm({ title: "", date: "", time: "" });
      setError("");
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="p-4 space-y-2">
      <input
        value={form.title}
        onChange={e => setForm({ ...form, title: e.target.value })}
        placeholder="Title"
        required
      />
      <input
        type="date"
        value={form.date}
        onChange={e => setForm({ ...form, date: e.target.value })}
        required
      />
      <input
        type="time"
        value={form.time}
        onChange={e => setForm({ ...form, time: e.target.value })}
        required
      />
      <button type="submit">Create Appointment</button>
      {error && <p style={{ color: "red" }}>{error}</p>}
    </form>
  );
}
