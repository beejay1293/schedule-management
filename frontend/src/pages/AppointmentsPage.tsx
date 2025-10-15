import React from "react";
import { AppointmentForm } from "../components/AppointmentForm";
import { AppointmentList } from "../components/AppointmentList";

export const AppointmentsPage: React.FC = () => {
  return (
    <div className="max-w-lg mx-auto mt-8">
      <AppointmentForm />
      <AppointmentList />
    </div>
  );
};
