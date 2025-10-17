import React, { useState } from "react";
import { AppointmentForm } from "../../components/appointments/AppointmentForm";
import { AppointmentList } from "../../components/appointments/AppointmentList";
import "./appointmentsPage.css";

export const AppointmentsPage: React.FC = () => {
  const [showFormModal, setShowFormModal] = useState(false);

  return (
    <div className="appointments-page">
      <div className="appointments-header">
        <h2>My Appointments</h2>
        <button className="create-btn" onClick={() => setShowFormModal(true)}>
          + Create Appointment
        </button>
      </div>

      <AppointmentList />
      {showFormModal && (
        <div className="modal-overlay">
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <button
              className="modal-close-modal-btn"
              onClick={() => setShowFormModal(false)}
            >
              ✕
            </button>
            <AppointmentForm onSuccess={() => setShowFormModal(false)} />
          </div>
        </div>
      )}
    </div>
  );
};
