import * as jspb from 'google-protobuf'



export class Appointment extends jspb.Message {
  getId(): string;
  setId(value: string): Appointment;

  getTitle(): string;
  setTitle(value: string): Appointment;

  getDate(): string;
  setDate(value: string): Appointment;

  getTime(): string;
  setTime(value: string): Appointment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Appointment.AsObject;
  static toObject(includeInstance: boolean, msg: Appointment): Appointment.AsObject;
  static serializeBinaryToWriter(message: Appointment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Appointment;
  static deserializeBinaryFromReader(message: Appointment, reader: jspb.BinaryReader): Appointment;
}

export namespace Appointment {
  export type AsObject = {
    id: string;
    title: string;
    date: string;
    time: string;
  };
}

export class CreateAppointmentRequest extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): CreateAppointmentRequest;

  getDate(): string;
  setDate(value: string): CreateAppointmentRequest;

  getTime(): string;
  setTime(value: string): CreateAppointmentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAppointmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAppointmentRequest): CreateAppointmentRequest.AsObject;
  static serializeBinaryToWriter(message: CreateAppointmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAppointmentRequest;
  static deserializeBinaryFromReader(message: CreateAppointmentRequest, reader: jspb.BinaryReader): CreateAppointmentRequest;
}

export namespace CreateAppointmentRequest {
  export type AsObject = {
    title: string;
    date: string;
    time: string;
  };
}

export class CreateAppointmentResponse extends jspb.Message {
  getAppointment(): Appointment | undefined;
  setAppointment(value?: Appointment): CreateAppointmentResponse;
  hasAppointment(): boolean;
  clearAppointment(): CreateAppointmentResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAppointmentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAppointmentResponse): CreateAppointmentResponse.AsObject;
  static serializeBinaryToWriter(message: CreateAppointmentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAppointmentResponse;
  static deserializeBinaryFromReader(message: CreateAppointmentResponse, reader: jspb.BinaryReader): CreateAppointmentResponse;
}

export namespace CreateAppointmentResponse {
  export type AsObject = {
    appointment?: Appointment.AsObject;
  };
}

export class ListAppointmentsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAppointmentsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAppointmentsRequest): ListAppointmentsRequest.AsObject;
  static serializeBinaryToWriter(message: ListAppointmentsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAppointmentsRequest;
  static deserializeBinaryFromReader(message: ListAppointmentsRequest, reader: jspb.BinaryReader): ListAppointmentsRequest;
}

export namespace ListAppointmentsRequest {
  export type AsObject = {
  };
}

export class ListAppointmentsResponse extends jspb.Message {
  getAppointmentsList(): Array<Appointment>;
  setAppointmentsList(value: Array<Appointment>): ListAppointmentsResponse;
  clearAppointmentsList(): ListAppointmentsResponse;
  addAppointments(value?: Appointment, index?: number): Appointment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAppointmentsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAppointmentsResponse): ListAppointmentsResponse.AsObject;
  static serializeBinaryToWriter(message: ListAppointmentsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAppointmentsResponse;
  static deserializeBinaryFromReader(message: ListAppointmentsResponse, reader: jspb.BinaryReader): ListAppointmentsResponse;
}

export namespace ListAppointmentsResponse {
  export type AsObject = {
    appointmentsList: Array<Appointment.AsObject>;
  };
}

export class DeleteAppointmentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeleteAppointmentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAppointmentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAppointmentRequest): DeleteAppointmentRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteAppointmentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAppointmentRequest;
  static deserializeBinaryFromReader(message: DeleteAppointmentRequest, reader: jspb.BinaryReader): DeleteAppointmentRequest;
}

export namespace DeleteAppointmentRequest {
  export type AsObject = {
    id: string;
  };
}

export class DeleteAppointmentResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAppointmentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAppointmentResponse): DeleteAppointmentResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteAppointmentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAppointmentResponse;
  static deserializeBinaryFromReader(message: DeleteAppointmentResponse, reader: jspb.BinaryReader): DeleteAppointmentResponse;
}

export namespace DeleteAppointmentResponse {
  export type AsObject = {
  };
}

