import * as grpcWeb from 'grpc-web';

import * as proto_appointment_pb from '../proto/appointment_pb'; // proto import: "proto/appointment.proto"


export class AppointmentServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  createAppointment(
    request: proto_appointment_pb.CreateAppointmentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_appointment_pb.CreateAppointmentResponse) => void
  ): grpcWeb.ClientReadableStream<proto_appointment_pb.CreateAppointmentResponse>;

  listAppointments(
    request: proto_appointment_pb.ListAppointmentsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_appointment_pb.ListAppointmentsResponse) => void
  ): grpcWeb.ClientReadableStream<proto_appointment_pb.ListAppointmentsResponse>;

  deleteAppointment(
    request: proto_appointment_pb.DeleteAppointmentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: proto_appointment_pb.DeleteAppointmentResponse) => void
  ): grpcWeb.ClientReadableStream<proto_appointment_pb.DeleteAppointmentResponse>;

}

export class AppointmentServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  createAppointment(
    request: proto_appointment_pb.CreateAppointmentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_appointment_pb.CreateAppointmentResponse>;

  listAppointments(
    request: proto_appointment_pb.ListAppointmentsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_appointment_pb.ListAppointmentsResponse>;

  deleteAppointment(
    request: proto_appointment_pb.DeleteAppointmentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<proto_appointment_pb.DeleteAppointmentResponse>;

}

