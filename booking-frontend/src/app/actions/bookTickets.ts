// src/app/actions/bookTickets.ts

export interface BookingInput {
  firstName: string;
  lastName: string;
  email: string;
  tickets: number;
}

export interface Booking {
  ID: number;
  FirstName: string;
  LastName: string;
  Email: string;
  NumberOfTickets: number;
  Conference: { Name: string };
}

const API_BASE = "http://localhost:8080/api"; // Go API base URL

// Server action: Book tickets
export const bookTickets = async (booking: BookingInput): Promise<string> => {
  try {
    const payload = {
      firstName: booking.firstName,
      lastName: booking.lastName,
      email: booking.email,
      tickets: booking.tickets,
    };

    const res = await fetch(`${API_BASE}/book`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    const data = await res.json();

    console.log("Booking response:", res.status, data); // Debug backend response

    if (!res.ok) {
      // Use backend message or fallback
      throw new Error(data.error || data.message || `Booking failed with status ${res.status}`);
    }

    return data.message || "Booking successful";
  } catch (err: any) {
    console.error("Error in bookTickets:", err);
    throw new Error(err.message || "Something went wrong while booking");
  }
};

// Server action: Fetch all bookings
export const fetchBookings = async (): Promise<Booking[]> => {
  try {
    const res = await fetch(`${API_BASE}/bookings`);

    const data = await res.json();

    console.log("Fetch bookings response:", res.status, data); // Debug backend response

    if (!res.ok) {
      throw new Error(data.error || data.message || `Failed to fetch bookings (status ${res.status})`);
    }

    return data;
  } catch (err: any) {
    console.error("Error in fetchBookings:", err);
    throw new Error(err.message || "Something went wrong while fetching bookings");
  }
};
