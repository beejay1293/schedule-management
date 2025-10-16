import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';
import { AppointmentServiceClient } from "../pb/appointment.client";
import {
  CreateAppointmentRequest,
  ListAppointmentsRequest,
  DeleteAppointmentRequest,
  SearchAppointmentsRequest,
  StreamAppointmentsRequest,
  Appointment as PbAppointment,
  AppointmentEvent,
  EventType,
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
  const utcDateTime = new Date(`${pb.date}T${pb.time}:00Z`);
  const localDate = utcDateTime.toLocaleDateString();
  const localTime = utcDateTime.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
  return {
    id: pb.id,
    title: pb.title,
    date: localDate, // "YYYY-MM-DD"
    time: localTime, // "HH:MM"
  };
}

// Create appointment
export async function createAppointment(payload: { title: string; date: string; time: string }): Promise<AppointmentTS> {
  const localDateTime = new Date(`${payload.date}T${payload.time}`);
  const utcDate = localDateTime.toISOString();

  payload.date = utcDate.split("T")[0]
  payload.time = utcDate.split("T")[1].slice(0,5)
  
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

export function streamAppointments(
  onData: (appt: AppointmentTS) => void,
  onDelete?: (id: string) => void,
  filters?: { title?: string; date?: string },
  reconnectDelay = 1
) {
  let stopped = false;
  let abortController: AbortController | null = null;

  const startStream = async () => {
    console.log("stream connected")
    abortController = new AbortController();

    const req = StreamAppointmentsRequest.create({
      titleFilter: filters?.title || "",
      dateFilter: filters?.date || "",
    });

    const call = grpcClient.streamAppointments(req, { abort: abortController.signal });

    try {
      for await (const event of call.responses) {
        if (stopped) break;
        handleStreamEvent(event, onData, onDelete);
      }
    } catch (err) {
      if (!stopped) {
        console.warn("Stream error, reconnecting...", err);
        setTimeout(startStream, reconnectDelay);
      }
    }
  };

  startStream();

  return {
    stop: () => {
      stopped = true;
      abortController?.abort(); // terminate the stream properly
    },
  };
}


function handleStreamEvent(
  event: AppointmentEvent,
  onData: (appt: AppointmentTS) => void,
  onDelete?: (id: string) => void
) {
  switch (event.type) {
    case EventType.CREATED:
      if (event.appointment) {
        onData(pbToAppointment(event.appointment));
      }
      break;

    case EventType.DELETED:
      if (event.appointment?.id && onDelete) {
        onDelete(event.appointment.id);
      }
      break;
  }
}