import { 
  CreateAppointmentRequest, 
  ListAppointmentsRequest, 
  DeleteAppointmentRequest, 
  Appointment as PbAppointment } from '../pb/appointment_pb';
import { AppointmentServiceClient } from "../pb/appointment_grpc_web_pb";

// gRPC-Web proxy URL
const PROXY_URL = import.meta.env.VITE_GRPC_PROXY_URL ?? "http://localhost:8080";

// Create gRPC client
const client = new AppointmentServiceClient(PROXY_URL, null, null);

// TypeScript-friendly Appointment type
export type Appointment = {
  id: string;
  title: string;
  date: string; // "YYYY-MM-DD"
  time: string; // "HH:MM"
};

// Helper: convert proto Appointment -> plain TS type
function pbToAppointment(pb: PbAppointment): Appointment {
  return {
    id: pb.getId(),
    title: pb.getTitle(),
    date: pb.getDate(),
    time: pb.getTime(),
  };
}

// Create appointment
export async function createAppointment(payload: {
  title: string;
  date: string;
  time: string;
}): Promise<Appointment> {
  const req = new CreateAppointmentRequest();
  req.setTitle(payload.title);
  req.setDate(payload.date);
  req.setTime(payload.time);

  return new Promise((resolve, reject) => {
    client.createAppointment(req, {}, (err, resp) => {
      if (err) return reject(err);
      if (!resp || !resp.getAppointment()) return reject(new Error("No response"));
      resolve(pbToAppointment(resp.getAppointment()!));
    });
  });
}

// List all appointments
export async function listAppointments(): Promise<Appointment[]> {
  const req = new ListAppointmentsRequest();

  return new Promise((resolve, reject) => {
    client.listAppointments(req, {}, (err, resp) => {
      if (err) return reject(err);
      if (!resp) return resolve([]);
      const items = resp.getAppointmentsList().map(pbToAppointment);

      // Sort chronologically by date + time
      items.sort((a: { date: any; time: any; }, b: { date: any; time: any; }) => `${a.date}T${a.time}`.localeCompare(`${b.date}T${b.time}`));
      resolve(items);
    });
  });
}

// Delete appointment
export async function deleteAppointment(id: string): Promise<void> {
  const req = new DeleteAppointmentRequest();
  req.setId(id);

  return new Promise((resolve, reject) => {
    client.deleteAppointment(req, {}, (err) => {
      if (err) reject(err);
      else resolve();
    });
  });
}
