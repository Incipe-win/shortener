import { useState } from 'react';
import { Link, useNavigate } from '@tanstack/react-router';
import { Menu, X, Link2, LogOut } from 'lucide-react';
import { Container } from './Container';
import { Button } from '@/components/ui/Button';
import { useAuth } from '@/stores/auth';

const publicLinks = [
  { to: '/' as const, label: '首页' },
] as const;

const authLinks = [
  { to: '/dashboard' as const, label: '仪表盘' },
  { to: '/monitor' as const, label: '监控' },
] as const;

export function Navbar() {
  const [open, setOpen] = useState(false);
  const { isAuthenticated, user, logout } = useAuth();
  const navigate = useNavigate();

  const allLinks = isAuthenticated ? [...publicLinks, ...authLinks] : publicLinks;

  const handleLogout = async () => {
    await logout();
    navigate({ to: '/' });
  };

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
          {allLinks.map(({ to, label }) => (
            <Link
              key={to}
              to={to}
              className="px-4 py-2 rounded-lg text-sm text-[var(--color-fg-muted)] hover:text-[var(--color-fg)] hover:bg-[var(--color-surface-hover)] transition-all duration-200"
              activeProps={{ className: 'text-[var(--color-fg)] bg-[var(--color-surface)]' }}
            >
              {label}
            </Link>
          ))}

          {/* Auth area */}
          <div className="ml-3 pl-3 border-l border-[var(--color-border)] flex items-center gap-2">
            {isAuthenticated ? (
              <>
                <span className="text-xs text-[var(--color-fg-muted)]">{user?.username}</span>
                <button
                  onClick={handleLogout}
                  className="p-2 rounded-lg text-[var(--color-fg-muted)] hover:text-[var(--color-fg)] hover:bg-[var(--color-surface-hover)] transition-colors"
                  title="登出"
                >
                  <LogOut className="w-4 h-4" />
                </button>
              </>
            ) : (
              <Link to="/login">
                <Button variant="secondary" className="text-xs px-3 py-1.5">登录</Button>
              </Link>
            )}
          </div>
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
            {allLinks.map(({ to, label }) => (
              <Link
                key={to}
                to={to}
                onClick={() => setOpen(false)}
                className="px-4 py-3 rounded-lg text-sm text-[var(--color-fg-muted)] hover:text-[var(--color-fg)] hover:bg-[var(--color-surface-hover)] transition-colors"
              >
                {label}
              </Link>
            ))}
            <div className="mt-2 pt-2 border-t border-[var(--color-border)]">
              {isAuthenticated ? (
                <button
                  onClick={() => { handleLogout(); setOpen(false); }}
                  className="w-full px-4 py-3 rounded-lg text-sm text-left text-[var(--color-fg-muted)] hover:text-[var(--color-fg)] hover:bg-[var(--color-surface-hover)] transition-colors"
                >
                  登出 ({user?.username})
                </button>
              ) : (
                <Link to="/login" onClick={() => setOpen(false)}>
                  <Button className="w-full">登录</Button>
                </Link>
              )}
            </div>
          </Container>
        </div>
      )}
    </nav>
  );
}
