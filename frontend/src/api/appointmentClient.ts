import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import { AppointmentServiceClient } from "../pb/appointment.client";
import {
  CreateAppointmentRequest,
  ListAppointmentsRequest,
  DeleteAppointmentRequest,
  SearchAppointmentsRequest,
  Appointment as PbAppointment,
} from "../pb/appointment";

// gRPC-Web proxy URL
const PROXY_URL = import.meta.env.VITE_GRPC_PROXY_URL ?? "http://localhost:8080";

// Transport and gRPC client
const transport = new GrpcWebFetchTransport({ baseUrl: PROXY_URL });
const grpcClient = new AppointmentServiceClient(transport);

export type AppointmentTS = {
  id: string;
  title: string;
  date: string; // "YYYY-MM-DD"
  time: string; // "HH:MM"
};

// Helper: convert proto Appointment -> plain TS type
function pbToAppointment(pb: PbAppointment): AppointmentTS {
  return {
    id: pb.id,
    title: pb.title,
    date: pb.date,
    time: pb.time,
  };
}

// Create appointment
export async function createAppointment(payload: { title: string; date: string; time: string }): Promise<AppointmentTS> {
  const req = CreateAppointmentRequest.create(payload);
  const resp = await grpcClient.createAppointment(req);
  if (!resp.response?.appointment) throw new Error("No appointment returned");
  return pbToAppointment(resp.response.appointment);
}

// List all appointments
export async function listAppointments(): Promise<AppointmentTS[]> {
  const req = ListAppointmentsRequest.create();
  const resp = await grpcClient.listAppointments(req);
  return (resp.response?.appointments || []).map(pbToAppointment);
}

// Delete appointment
export async function deleteAppointment(id: string): Promise<void> {
  const req = DeleteAppointmentRequest.create({ id });
  await grpcClient.deleteAppointment(req);
}

// Search appointments
export async function searchAppointments(payload: { title?: string; date?: string }): Promise<AppointmentTS[]> {
  const req = SearchAppointmentsRequest.create(payload);
  const resp = await grpcClient.searchAppointments(req);
  return (resp.response?.appointments || []).map(pbToAppointment);
}
