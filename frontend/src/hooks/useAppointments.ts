import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  listAppointments,
  createAppointment,
  deleteAppointment,
  type Appointment,
} from "../api/appointmentClient";

export function useAppointments() {
  const queryClient = useQueryClient();

  const appointmentsQuery = useQuery<Appointment[]>({
    queryKey: ["appointments"],
    queryFn: listAppointments,
  });

  const createMutation = useMutation({
    mutationFn: createAppointment,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["appointments"] }),
  });

  const deleteMutation = useMutation({
    mutationFn: deleteAppointment,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["appointments"] }),
  });

  return { appointmentsQuery, createMutation, deleteMutation };
}
