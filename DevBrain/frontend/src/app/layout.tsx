import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "DevBrain",
  description: "An AI-powered knowledge graph for developers",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="h-full antialiased">
      <body className="min-h-full flex flex-col">{children}</body>
    </html>
  );
}
