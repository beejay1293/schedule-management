import React, { useState, useEffect, useCallback } from "react";
import { useAppointments } from "../../hooks/useAppointments";
import { ConfirmationModal } from "../modals/Confirm";
import {
  searchAppointments,
  streamAppointments,
} from "../../api/appointmentClient";
import "./appointments.css";
import calendar from "../../assets/calendar.png";

export const AppointmentList: React.FC = () => {
  const { appointmentsQuery, deleteMutation } = useAppointments();
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  const [searchTitle, setSearchTitle] = useState("");
  const [searchDate, setSearchDate] = useState("");
  const [filteredAppointments, setFilteredAppointments] = useState(
    appointmentsQuery.data || []
  );


  useEffect(() => {
    const stream = streamAppointments(
      // Handle new appointment
      (appt) => {
        setFilteredAppointments((prev) => {
          const exists = prev.find((a) => a.id === appt.id);
          if (exists) return prev;
          return [...prev, appt];
        });
      },
      // Handle deleted appointment
      (deletedId) => {
        setFilteredAppointments((prev) =>
          prev.filter((a) => a.id !== deletedId)
        );
      }
    );

    // Stop streaming when unmounted
    return () => {
      stream.stop();
    };
  }, []);

  const debounce = (fn: Function, delay = 500) => {
    let timeout: NodeJS.Timeout;
    return (...args: any[]) => {
      clearTimeout(timeout);
      timeout = setTimeout(() => fn(...args), delay);
    };
  };

  const handleSearch = useCallback(
    debounce(async () => {
      if (!searchTitle && !searchDate) {
        setFilteredAppointments(appointmentsQuery.data || []);
        return;
      }
      try {
        const results = await searchAppointments({
          title: searchTitle || undefined,
          date: searchDate || undefined,
        });
        setFilteredAppointments(results);
      } catch (err) {
        console.error("Search error:", err);
      }
    }, 500),
    [searchTitle, searchDate, appointmentsQuery.data]
  );

  useEffect(() => {
    handleSearch();
  }, [searchTitle, searchDate, appointmentsQuery.data, handleSearch]);

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
      <div className="search-bar">
        <input
          type="text"
          placeholder="Search by title..."
          value={searchTitle}
          onChange={(e) => setSearchTitle(e.target.value)}
          className="search-input"
        />
        <input
          type="date"
          value={searchDate}
          onChange={(e) => setSearchDate(e.target.value)}
          className="search-input"
        />
      </div>

      <ul className="appointment-list">
        {appointmentsQuery.isLoading ? (
          Array.from({ length: 5 }).map((_, idx) => (
            <li key={idx} className="appointment-item skeleton">
              <span className="skeleton-text">&nbsp;</span>
              <button
                className="appointment-delete-btn skeleton-btn"
                disabled
              />
            </li>
          ))
        ) : filteredAppointments.length == 0 ? (
          <div className="no-appointments-card">
            <img
              src={calendar}
              alt="No appointments"
              className="no-appointments-img"
            />
            <p className="no-appointments-text">No appointments scheduled</p>
          </div>
        ) : (
          filteredAppointments.map((a) => (
            <li key={a.id} className="appointment-card">
              <div className="card-content">
                <h3>{a.title}</h3>
                <p>
                  {a.date} at {a.time}
                </p>
              </div>
              <button
                className="appointment-delete-btn"
                onClick={() => handleDeleteClick(a.id)}
                disabled={deleteMutation.isPending}
              >
                {deleteMutation.isPending ? "Deleting..." : "Delete"}
              </button>
            </li>
          ))
        )}
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
