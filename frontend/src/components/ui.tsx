export function PageHeader({
  eyebrow,
  title,
  description,
}: {
  eyebrow: string;
  title: string;
  description?: string;
}) {
  return (
    <div className="mb-8">
      <div className="font-label text-xs text-accent mb-2">&gt; {eyebrow}</div>
      <h1 className="text-2xl font-semibold tracking-tight text-ink">
        {title}
      </h1>
      {description && (
        <p className="text-sm text-ink-soft mt-1.5 max-w-lg">{description}</p>
      )}
    </div>
  );
}

export function Card({
  children,
  className = "",
}: {
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <div
      className={`border border-line rounded-lg bg-paper p-5 ${className}`}
    >
      {children}
    </div>
  );
}

export function Button({
  children,
  variant = "primary",
  className = "",
  ...props
}: React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "ghost";
}) {
  const base =
    "inline-flex items-center justify-center gap-2 text-sm font-medium rounded-md px-3.5 py-2 transition-colors disabled:opacity-40 disabled:cursor-not-allowed";

  const variants: Record<string, string> = {
    primary: "bg-ink text-paper hover:bg-accent",
    secondary: "border border-line text-ink hover:border-ink",
    ghost: "text-ink-soft hover:text-ink",
  };

  return (
    <button className={`${base} ${variants[variant]} ${className}`} {...props}>
      {children}
    </button>
  );
}

export function Input(props: React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      {...props}
      className={`border border-line rounded-md px-3 py-2 text-sm bg-paper text-ink placeholder:text-ink-soft/60 focus:border-accent outline-none transition-colors ${
        props.className ?? ""
      }`}
    />
  );
}

export function EmptyState({
  title,
  description,
  action,
}: {
  title: string;
  description: string;
  action?: React.ReactNode;
}) {
  return (
    <div className="border border-dashed border-line rounded-lg py-14 px-6 text-center">
      <p className="text-sm font-medium text-ink">{title}</p>
      <p className="text-sm text-ink-soft mt-1 max-w-sm mx-auto">
        {description}
      </p>
      {action && <div className="mt-4 flex justify-center">{action}</div>}
    </div>
  );
}

export function ProgressBar({ value }: { value: number }) {
  return (
    <div className="h-1.5 w-full rounded-full bg-line-soft overflow-hidden">
      <div
        className="h-full bg-accent rounded-full transition-all"
        style={{ width: `${Math.min(100, Math.max(0, value))}%` }}
      />
    </div>
  );
}
