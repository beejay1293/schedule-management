import React, { useState } from "react";
import { useAppointments } from "../hooks/useAppointments";

export const AppointmentForm: React.FC = () => {
  const { createMutation } = useAppointments();
  const [form, setForm] = useState({ title: "", date: "", time: "" });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createMutation.mutate(form);
    setForm({ title: "", date: "", time: "" });
  };

  return (
    <form onSubmit={handleSubmit} className="p-4 border rounded mb-4">
      <h2 className="text-lg font-bold mb-2">Create Appointment</h2>
      <input
        type="text"
        placeholder="Title"
        value={form.title}
        onChange={(e) => setForm({ ...form, title: e.target.value })}
        className="border p-2 w-full mb-2 rounded"
        required
      />
      <input
        type="date"
        value={form.date}
        onChange={(e) => setForm({ ...form, date: e.target.value })}
        className="border p-2 w-full mb-2 rounded"
        required
      />
      <input
        type="time"
        value={form.time}
        onChange={(e) => setForm({ ...form, time: e.target.value })}
        className="border p-2 w-full mb-4 rounded"
        required
      />
      <button
        type="submit"
        className="bg-blue-500 text-white px-4 py-2 rounded disabled:opacity-50"
        disabled={createMutation.isPending}
      >
        {createMutation.isPending ? "Creating..." : "Create"}
      </button>
      {createMutation.isError && (
        <p className="text-red-500 mt-2">
          {(createMutation.error as any).message || "Error creating appointment"}
        </p>
      )}
    </form>
  );
};
