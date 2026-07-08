"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import {
  LayoutDashboard,
  Library,
  Network,
  Target,
  Search,
  MessageSquare,
  BarChart3,
  LogOut,
} from "lucide-react";

const NAV_ITEMS = [
  { href: "/dashboard", label: "dashboard", icon: LayoutDashboard },
  { href: "/resources", label: "resources", icon: Library },
  { href: "/graph", label: "graph", icon: Network },
  { href: "/gaps", label: "gaps", icon: Target },
  { href: "/search", label: "search", icon: Search },
  { href: "/mentor", label: "mentor", icon: MessageSquare },
  { href: "/analytics", label: "analytics", icon: BarChart3 },
];

export default function AppShell({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();

  function handleLogout() {
    localStorage.removeItem("token");
    router.push("/login");
  }

  return (
    <div className="min-h-screen flex bg-paper text-ink">
      <aside className="w-56 shrink-0 border-r border-line flex flex-col justify-between py-6 px-4 h-screen sticky top-0">
        <div>
          <Link href="/dashboard" className="block mb-8 px-2">
            <span className="font-label text-sm text-ink">
              &gt; devbrain
            </span>
          </Link>

          <nav className="flex flex-col gap-1">
            {NAV_ITEMS.map((item) => {
              const Icon = item.icon;
              const active = pathname === item.href;

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`flex items-center gap-2.5 px-2.5 py-2 rounded-md text-sm font-label transition-colors ${
                    active
                      ? "bg-accent-soft text-accent"
                      : "text-ink-soft hover:bg-line-soft hover:text-ink"
                  }`}
                >
                  <Icon size={15} strokeWidth={2} />
                  {item.label}
                </Link>
              );
            })}
          </nav>
        </div>

        <button
          onClick={handleLogout}
          className="flex items-center gap-2.5 px-2.5 py-2 rounded-md text-sm font-label text-ink-soft hover:bg-line-soft hover:text-danger transition-colors"
        >
          <LogOut size={15} strokeWidth={2} />
          log out
        </button>
      </aside>

      <main className="flex-1 min-w-0">
        <div className="max-w-4xl mx-auto px-8 py-10">{children}</div>
      </main>
    </div>
  );
}
