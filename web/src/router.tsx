import { createRouter, createRoute, createRootRoute, Outlet } from '@tanstack/react-router';
import { Background } from '@/components/layout/Background';
import { Navbar } from '@/components/layout/Navbar';
import { HomePage } from '@/pages/Home';

// Root layout
const rootRoute = createRootRoute({
  component: () => (
    <Background>
      <Navbar />
      <main>
        <Outlet />
      </main>
    </Background>
  ),
});

// Public routes
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: HomePage,
});

const dashboardRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/dashboard',
  component: () => <div className="pt-24 text-center text-[var(--color-fg-muted)]">Dashboard — Phase 4</div>,
});

const monitorRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/monitor',
  component: () => <div className="pt-24 text-center text-[var(--color-fg-muted)]">Monitor — Phase 4</div>,
});

// Router
export const router = createRouter({
  routeTree: rootRoute.addChildren([indexRoute, dashboardRoute, monitorRoute]),
  defaultPreload: 'intent',
});

// Type registration
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
