import { useState } from 'react';
import { Link } from '@tanstack/react-router';
import { Menu, X, Link2 } from 'lucide-react';
import { Container } from './Container';

const navLinks = [
  { to: '/', label: '首页' },
  { to: '/dashboard', label: '仪表盘' },
  { to: '/monitor', label: '监控' },
] as const;

export function Navbar() {
  const [open, setOpen] = useState(false);

  return (
    <nav className="fixed top-0 inset-x-0 z-50 border-b border-[var(--color-border)] bg-[var(--color-bg-base)]/80 backdrop-blur-xl">
      <Container className="flex items-center justify-between h-16">
        {/* Logo */}
        <Link to="/" className="flex items-center gap-2 text-[var(--color-fg)] font-semibold text-lg tracking-tight">
          <div className="w-8 h-8 rounded-lg bg-[var(--color-accent)] flex items-center justify-center">
            <Link2 className="w-4 h-4 text-white" />
          </div>
          Shortener
        </Link>

        {/* Desktop links */}
        <div className="hidden md:flex items-center gap-1">
          {navLinks.map(({ to, label }) => (
            <Link
              key={to}
              to={to}
              className="px-4 py-2 rounded-lg text-sm text-[var(--color-fg-muted)] hover:text-[var(--color-fg)] hover:bg-[var(--color-surface-hover)] transition-all duration-200"
              activeProps={{ className: 'text-[var(--color-fg)] bg-[var(--color-surface)]' }}
            >
              {label}
            </Link>
          ))}
        </div>

        {/* Mobile toggle */}
        <button
          onClick={() => setOpen(!open)}
          className="md:hidden p-2 rounded-lg text-[var(--color-fg-muted)] hover:text-[var(--color-fg)] hover:bg-[var(--color-surface-hover)] transition-colors"
          aria-label="Toggle menu"
        >
          {open ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
        </button>
      </Container>

      {/* Mobile menu */}
      {open && (
        <div className="md:hidden bg-[var(--color-bg-base)]/95 backdrop-blur-xl border-b border-[var(--color-border)]">
          <Container className="py-4 flex flex-col gap-1">
            {navLinks.map(({ to, label }) => (
              <Link
                key={to}
                to={to}
                onClick={() => setOpen(false)}
                className="px-4 py-3 rounded-lg text-sm text-[var(--color-fg-muted)] hover:text-[var(--color-fg)] hover:bg-[var(--color-surface-hover)] transition-colors"
              >
                {label}
              </Link>
            ))}
          </Container>
        </div>
      )}
    </nav>
  );
}
