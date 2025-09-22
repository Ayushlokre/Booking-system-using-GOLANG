export interface BookingRequest {
  firstName: string;
  lastName: string;
  email: string;
  tickets: number;
}

export interface BookingResponse {
  message: string;
}

export interface BookingData {
  ID: number;
  FirstName: string;
  LastName: string;
  Email: string;
  NumberOfTickets: number;
  Conference: {
    Name: string;
  };
}

// Book tickets
export async function bookTicket(req: BookingRequest): Promise<BookingResponse> {
  const res = await fetch("http://localhost:8080/api/book", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(errorText || "Booking failed");
  }
  return res.json();
}

// Fetch all bookings
export async function fetchBookings(): Promise<BookingData[]> {
  const res = await fetch("http://localhost:8080/api/bookings");
  if (!res.ok) throw new Error("Failed to fetch bookings");
  return res.json();
}
