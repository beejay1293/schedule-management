export async function listAppointments() {
    const res = await fetch("http://localhost:8080/appointments");
    return res.json();
  }
  
  export async function createAppointment(data: { title: string; date: string; time: string; }) {
    const res = await fetch("http://localhost:8080/appointments", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
  }
  
  export async function deleteAppointment(id: string) {
    const res = await fetch(`http://localhost:8080/appointments/${id}`, { method: "DELETE" });
    if (!res.ok) throw new Error(await res.text());
  }
  