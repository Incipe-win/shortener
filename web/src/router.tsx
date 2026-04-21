import { createRouter, createRoute, createRootRoute, Outlet, redirect } from '@tanstack/react-router';
import { Background } from '@/components/layout/Background';
import { Navbar } from '@/components/layout/Navbar';
import { HomePage } from '@/pages/Home';
import { PreviewPage } from '@/pages/Preview';
import { LoginPage } from '@/pages/Login';
import { DashboardPage } from '@/pages/Dashboard';
import { MonitorPage } from '@/pages/Monitor';
import { useAuth } from '@/stores/auth';

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

// ── Public Routes ────────────────────────

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: HomePage,
});

const previewRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/preview/$surl',
  component: PreviewPage,
});

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  component: LoginPage,
});

// ── Auth-Protected Routes ────────────────

const dashboardRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/dashboard',
  beforeLoad: () => {
    const { isAuthenticated } = useAuth.getState();
    if (!isAuthenticated) {
      throw redirect({ to: '/login' });
    }
  },
  component: DashboardPage,
});

const monitorRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/monitor',
  beforeLoad: () => {
    const { isAuthenticated } = useAuth.getState();
    if (!isAuthenticated) {
      throw redirect({ to: '/login' });
    }
  },
  component: MonitorPage,
});

// ── Router ───────────────────────────────

export const router = createRouter({
  routeTree: rootRoute.addChildren([
    indexRoute,
    previewRoute,
    loginRoute,
    dashboardRoute,
    monitorRoute,
  ]),
  defaultPreload: 'intent',
});

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
