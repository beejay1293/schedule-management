import React from "react";
import { AppointmentForm } from "../../components/appointments/AppointmentForm";
import { AppointmentList } from "../../components/appointments/AppointmentList";
import "./appointmentsPage.css"

export const AppointmentsPage: React.FC = () => {
  return (
    <div className="appointments-page">
      <AppointmentForm />
      <AppointmentList />
    </div>
  );
};
