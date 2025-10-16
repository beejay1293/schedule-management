import React from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { AppointmentsPage } from "./pages/appointments/AppointmentsPage";

const queryClient = new QueryClient();

const App: React.FC = () => (
  <QueryClientProvider client={queryClient}>
    <AppointmentsPage />
  </QueryClientProvider>
);

export default App;
