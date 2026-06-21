"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";

export default function DashboardPage() {
  const [user, setUser] = useState<any>(null);

  useEffect(() => {
    async function loadUser() {
      try {
        const response = await api.get("/me");

        setUser(response.data);
      } catch (error) {
        console.error(error);
      }
    }

    loadUser();
  }, []);

  if (!user) {
    return <div>Loading...</div>;
  }

  return (
    <main className="min-h-screen p-10">

      <h1 className="text-4xl font-bold">
        Welcome {user.name}
      </h1>

      <p>{user.email}</p>

    </main>
  );
}