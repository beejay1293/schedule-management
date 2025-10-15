import React from "react";
import { useAppointments } from "../hooks/useAppointments";

export const AppointmentList: React.FC = () => {
  const { appointmentsQuery, deleteMutation } = useAppointments();

  if (appointmentsQuery.isLoading) return <p>Loading...</p>;
  if (appointmentsQuery.isError)
    return <p>Error loading appointments.</p>;

  return (
    <div className="p-4">
      <h2 className="text-xl font-bold mb-4">Appointments</h2>
      <ul className="space-y-2">
        {appointmentsQuery.data?.map((a) => (
          <li
            key={a.id}
            className="flex justify-between items-center border p-2 rounded"
          >
            <span>
              <strong>{a.title}</strong> — {a.date} at {a.time}
            </span>
            <button
              className="text-red-500"
              onClick={() => deleteMutation.mutate(a.id)}
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending ? "Deleting..." : "Delete"}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
};
