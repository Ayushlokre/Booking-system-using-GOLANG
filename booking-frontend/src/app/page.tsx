"use client";

import { useState, useEffect } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Table, TableHeader, TableRow, TableCell, TableBody } from "@/components/ui/table";
import { toast } from "sonner";
import { bookTickets, fetchBookings as fetchBookingsAction } from "./actions/bookTickets";

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

  // Fetch all bookings from Go API
  const fetchBookings = async () => {
    try {
      const data = await fetchBookingsAction();
      setBookings(data);
    } catch (err: any) {
      console.error(err);
      toast.error(err.message || "Failed to fetch bookings");
    }
  };

  useEffect(() => {
    fetchBookings();
  }, []);

  // Handle booking submission
  const handleBooking = async () => {
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
      const message = await bookTickets({ firstName, lastName, email, tickets });
      toast.success(message);

      // Reset form fields
      setFirstName("");
      setLastName("");
      setEmail("");
      setTickets(1);

      fetchBookings();
    } catch (err: any) {
      toast.error(err.message || "Booking failed!");
    }
  };

  return (
    <div className="max-w-4xl mx-auto p-8">
      <h1 className="text-3xl font-bold mb-6">🎟 Conference Booking App</h1>

      {/* Booking Form */}
      <div className="grid gap-4 mb-6">
        <Input placeholder="First Name" value={firstName} onChange={(e) => setFirstName(e.target.value)} />
        <Input placeholder="Last Name" value={lastName} onChange={(e) => setLastName(e.target.value)} />
        <Input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
        <Input
          type="number"
          min={1}
          placeholder="Number of Tickets"
          value={tickets}
          onChange={(e) => setTickets(Number(e.target.value))}
        />
        <Button onClick={handleBooking}>Book Tickets</Button>
      </div>

      {/* Bookings Table */}
      <h2 className="text-2xl font-semibold mb-4">All Bookings</h2>
      <div className="overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>First Name</TableCell>
              <TableCell>Last Name</TableCell>
              <TableCell>Email</TableCell>
              <TableCell>Tickets</TableCell>
              <TableCell>Conference</TableCell>
            </TableRow>
          </TableHeader>
          <TableBody>
            {bookings.map((b) => (
              <TableRow key={b.ID}>
                <TableCell>{b.ID}</TableCell>
                <TableCell>{b.FirstName}</TableCell>
                <TableCell>{b.LastName}</TableCell>
                <TableCell>{b.Email}</TableCell>
                <TableCell>{b.NumberOfTickets}</TableCell>
                <TableCell>{b.Conference.Name}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
