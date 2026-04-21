import { createRouter, createRoute, createRootRoute, Outlet, redirect } from '@tanstack/react-router';
import { Background } from '@/components/layout/Background';
import { Navbar } from '@/components/layout/Navbar';
import { HomePage } from '@/pages/Home';
import { PreviewPage } from '@/pages/Preview';
import { LoginPage } from '@/pages/Login';
import { DashboardPage } from '@/pages/Dashboard';
import { MonitorPage } from '@/pages/Monitor';
import { NotFoundPage } from '@/pages/NotFound';
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
  beforeLoad: async () => {
    await useAuth.getState().checkAuth();
    if (useAuth.getState().isAuthenticated) {
      throw redirect({ to: '/dashboard' });
    }
  },
  component: LoginPage,
});

// ── Auth-Protected Routes ────────────────

const dashboardRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/dashboard',
  beforeLoad: async () => {
    await useAuth.getState().checkAuth();
    if (!useAuth.getState().isAuthenticated) {
      throw redirect({ to: '/login' });
    }
  },
  component: DashboardPage,
});

const monitorRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/monitor',
  beforeLoad: async () => {
    await useAuth.getState().checkAuth();
    if (!useAuth.getState().isAuthenticated) {
      throw redirect({ to: '/login' });
    }
  },
  component: MonitorPage,
});

const notFoundRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/404/$surl',
  component: NotFoundPage,
});

// ── Router ───────────────────────────────

export const router = createRouter({
  routeTree: rootRoute.addChildren([
    indexRoute,
    previewRoute,
    loginRoute,
    dashboardRoute,
    monitorRoute,
    notFoundRoute,
  ]),
  defaultPreload: 'intent',
});

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
