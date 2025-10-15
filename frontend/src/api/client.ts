import { AppointmentServiceClient } from "../pb/appointment_grpc_web_pb";

const client = new AppointmentServiceClient("http://localhost:8080", null, null);

// function pbToAppointment(pb: PbAppointment): Appointment {
//   return {
//     id: pb.getId(),
//     title: pb.getTitle(),
//     date: pb.getDate(),
//     time: pb.getTime(),
//   };
// }

export default client;
