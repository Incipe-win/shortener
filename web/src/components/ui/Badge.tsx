import type { ReactNode } from 'react';

type Level = 'safe' | 'warning' | 'danger' | 'pending';

const levelStyles: Record<Level, string> = {
  safe: 'bg-[var(--color-safe)]/15 text-[var(--color-safe)] border-[var(--color-safe)]/30',
  warning: 'bg-[var(--color-warning)]/15 text-[var(--color-warning)] border-[var(--color-warning)]/30',
  danger: 'bg-[var(--color-danger)]/15 text-[var(--color-danger)] border-[var(--color-danger)]/30',
  pending: 'bg-white/5 text-[var(--color-fg-muted)] border-white/10',
};

const levelLabels: Record<Level, string> = {
  safe: '安全',
  warning: '警告',
  danger: '危险',
  pending: '待检',
};

interface BadgeProps {
  level: Level;
  children?: ReactNode;
  className?: string;
}

export function Badge({ level, children, className = '' }: BadgeProps) {
  return (
    <span
      className={[
        'inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full',
        'text-xs font-medium border',
        'transition-colors duration-200',
        levelStyles[level],
        className,
      ].join(' ')}
    >
      <span className={`w-1.5 h-1.5 rounded-full ${level === 'safe' ? 'bg-[var(--color-safe)]' : level === 'warning' ? 'bg-[var(--color-warning)]' : level === 'danger' ? 'bg-[var(--color-danger)]' : 'bg-[var(--color-fg-muted)]'}`} />
      {children ?? levelLabels[level]}
    </span>
  );
}
