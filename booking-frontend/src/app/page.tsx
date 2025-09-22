"use client";

import { useState, useEffect } from "react";
import { toast } from "sonner";

interface Booking {
  ID: number;
  FirstName: string;
  LastName: string;
  Email: string;
  NumberOfTickets: number;
  Conference: { Name: string };
}

export default function Home() {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [tickets, setTickets] = useState(1);
  const [bookings, setBookings] = useState<Booking[]>([]);

  // Fetch all bookings from backend
  const fetchBookings = async () => {
    try {
      const res = await fetch("http://localhost:8080/api/bookings");
      const data = await res.json();
      setBookings(data);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchBookings();
  }, []);

  const handleBooking = async () => {
    // --- Client-side validation ---
    if (!firstName || !lastName || !email || tickets <= 0) {
      toast.error("Please fill all fields correctly!");
      return;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      toast.error("Please enter a valid email");
      return;
    }

    try {
      const res = await fetch("http://localhost:8080/api/book", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          firstName, // ✅ Must match Go struct
          lastName,
          email,
          tickets,
        }),
      });

      const data = await res.json();

      if (!res.ok) {
        toast.error(data.error || "Booking failed!");
        return;
      }

      toast.success(data.message);

      // Reset form fields
      setFirstName("");
      setLastName("");
      setEmail("");
      setTickets(1);

      fetchBookings();
    } catch (err) {
      console.error(err);
      toast.error("Something went wrong!");
    }
  };

  return (
    <div className="max-w-4xl mx-auto p-8">
      <h1 className="text-3xl font-bold mb-6">🎟 Conference Booking App</h1>

      <div className="grid gap-4 mb-6">
        <input
          type="text"
          placeholder="First Name"
          value={firstName}
          onChange={(e) => setFirstName(e.target.value)}
          className="border px-3 py-2 rounded-md"
        />
        <input
          type="text"
          placeholder="Last Name"
          value={lastName}
          onChange={(e) => setLastName(e.target.value)}
          className="border px-3 py-2 rounded-md"
        />
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="border px-3 py-2 rounded-md"
        />
        <input
          type="number"
          min={1}
          placeholder="Number of Tickets"
          value={tickets}
          onChange={(e) => setTickets(Number(e.target.value))}
          className="border px-3 py-2 rounded-md w-32"
        />
        <button
          onClick={handleBooking}
          className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition"
        >
          Book Tickets
        </button>
      </div>

      <h2 className="text-2xl font-semibold mb-4">All Bookings</h2>
      <div className="overflow-x-auto">
        <table className="w-full border-collapse border border-gray-300">
          <thead className="bg-gray-100">
            <tr>
              <th className="border px-2 py-1">ID</th>
              <th className="border px-2 py-1">First Name</th>
              <th className="border px-2 py-1">Last Name</th>
              <th className="border px-2 py-1">Email</th>
              <th className="border px-2 py-1">Tickets</th>
              <th className="border px-2 py-1">Conference</th>
            </tr>
          </thead>
          <tbody>
            {bookings.map((b) => (
              <tr key={b.ID} className="hover:bg-gray-50">
                <td className="border px-2 py-1">{b.ID}</td>
                <td className="border px-2 py-1">{b.FirstName}</td>
                <td className="border px-2 py-1">{b.LastName}</td>
                <td className="border px-2 py-1">{b.Email}</td>
                <td className="border px-2 py-1">{b.NumberOfTickets}</td>
                <td className="border px-2 py-1">{b.Conference.Name}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
