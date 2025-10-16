import React, { useState } from "react";
import { useAppointments } from "../../hooks/useAppointments";
import { ConfirmationModal } from "../modals/Confirm";
import "./Appointments.css";

export const AppointmentList: React.FC = () => {
  const { appointmentsQuery, deleteMutation } = useAppointments();
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  const handleDeleteClick = (id: string) => {
    setSelectedId(id);
    setModalOpen(true);
  };

  const confirmDelete = () => {
    if (selectedId) deleteMutation.mutate(selectedId);
    setModalOpen(false);
    setSelectedId(null);
  };

  return (
    <div className="appointment-list-container">
      <ul className="appointment-list">
        {appointmentsQuery.isLoading
          ? Array.from({ length: 5 }).map((_, idx) => (
              <li key={idx} className="appointment-item skeleton">
                <span className="skeleton-text">&nbsp;</span>
                <button className="appointment-delete-btn skeleton-btn" disabled />
              </li>
            ))
          : appointmentsQuery.data?.map((a) => (
              <li key={a.id} className="appointment-item">
                <span>
                  <strong>{a.title}</strong> — {a.date} at {a.time}
                </span>
                <button
                  className="appointment-delete-btn"
                  onClick={() => handleDeleteClick(a.id)}
                  disabled={deleteMutation.isPending}
                >
                  {deleteMutation.isPending ? "Deleting..." : "Delete"}
                </button>
              </li>
            ))}
      </ul>

      <ConfirmationModal
        isOpen={modalOpen}
        title="Confirm Delete"
        message="Are you sure you want to delete this appointment?"
        onConfirm={confirmDelete}
        onCancel={() => setModalOpen(false)}
      />
    </div>
  );
};
