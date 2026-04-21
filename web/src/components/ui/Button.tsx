import { type ButtonHTMLAttributes } from 'react';

type Variant = 'primary' | 'secondary' | 'ghost';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  loading?: boolean;
}

const styles: Record<Variant, string> = {
  primary: [
    'bg-[var(--color-accent)] text-white',
    'shadow-[var(--shadow-accent-glow)]',
    'hover:bg-[var(--color-accent-bright)] hover:shadow-[0_0_0_1px_rgba(94,106,210,0.6),0_6px_20px_rgba(94,106,210,0.4),inset_0_1px_0_0_rgba(255,255,255,0.25)]',
    'active:scale-[0.98] active:shadow-[0_0_0_1px_rgba(94,106,210,0.4),0_2px_8px_rgba(94,106,210,0.2)]',
  ].join(' '),
  secondary: [
    'bg-[var(--color-surface)] text-[var(--color-fg)]',
    'shadow-[var(--shadow-inner-highlight)]',
    'hover:bg-[var(--color-surface-hover)] hover:shadow-[inset_0_1px_0_0_rgba(255,255,255,0.12),0_0_20px_rgba(94,106,210,0.05)]',
    'active:scale-[0.98]',
  ].join(' '),
  ghost: [
    'bg-transparent text-[var(--color-fg-muted)]',
    'hover:bg-[var(--color-surface)] hover:text-[var(--color-fg)]',
    'active:scale-[0.98]',
  ].join(' '),
};

export function Button({ variant = 'primary', loading, children, className = '', disabled, ...rest }: ButtonProps) {
  return (
    <button
      className={[
        'relative inline-flex items-center justify-center gap-2',
        'px-5 py-2.5 rounded-lg text-sm font-medium',
        'transition-all duration-200 ease-[var(--ease-expo-out)]',
        'focus:outline-none focus:ring-2 focus:ring-[var(--color-accent)]/50 focus:ring-offset-2 focus:ring-offset-[var(--color-bg-base)]',
        'disabled:opacity-50 disabled:pointer-events-none',
        styles[variant],
        className,
      ].join(' ')}
      disabled={disabled || loading}
      {...rest}
    >
      {loading && (
        <svg className="animate-spin -ml-1 w-4 h-4" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="3" className="opacity-25" />
          <path d="M4 12a8 8 0 018-8" stroke="currentColor" strokeWidth="3" strokeLinecap="round" className="opacity-75" />
        </svg>
      )}
      {children}
    </button>
  );
}
