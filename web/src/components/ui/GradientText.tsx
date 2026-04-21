import type { ReactNode } from 'react';

interface GradientTextProps {
  children: ReactNode;
  className?: string;
  accent?: boolean;
}

export function GradientText({ children, className = '', accent = false }: GradientTextProps) {
  if (accent) {
    return (
      <span
        className={`bg-gradient-to-r from-[var(--color-accent)] via-indigo-400 to-[var(--color-accent)] bg-clip-text text-transparent bg-[length:200%_auto] ${className}`}
        style={{ animation: 'shimmer 3s linear infinite' }}
      >
        {children}
      </span>
    );
  }

  return (
    <span className={`bg-gradient-to-b from-white via-white/95 to-white/70 bg-clip-text text-transparent ${className}`}>
      {children}
    </span>
  );
}
