"use client";

import { ReactNode } from "react";
import { Toaster } from "sonner";

interface ClientProviderProps {
  children: ReactNode;
}

export default function ClientProvider({ children }: ClientProviderProps) {
  return (
    <>
      {children}
      <Toaster position="top-center" richColors />
    </>
  );
}
