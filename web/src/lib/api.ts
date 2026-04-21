import type { ConvertRequest, ConvertResponse, PreviewResponse } from '@/types';

const API_BASE = '/api';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: 'include', // HTTP-only cookie auth
    ...init,
  });
  if (!res.ok) {
    const text = await res.text().catch(() => 'Unknown error');
    throw new Error(`${res.status}: ${text}`);
  }
  return res.json();
}

// ── Public APIs ──────────────────────────────────

export async function convertUrl(longUrl: string): Promise<ConvertResponse> {
  return request<ConvertResponse>('/convert', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ long_url: longUrl } satisfies ConvertRequest),
  });
}

export async function previewLink(surl: string): Promise<PreviewResponse> {
  return request<PreviewResponse>(`/preview/${surl}`);
}

// ── Auth APIs ────────────────────────────────────

export async function login(username: string, password: string) {
  return request<{ message: string; username: string; user_id: number }>('/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
}

export async function register(username: string, password: string) {
  return request<{ message: string; username: string; user_id: number }>('/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
}

export async function logout() {
  return request<{ message: string }>('/auth/logout', { method: 'POST' });
}

export async function getMe() {
  return request<{ username: string; user_id: number }>('/auth/me');
}

export async function getUnregisteredRemaining(): Promise<{ remaining: number }> {
  try {
    const res = await fetch('/api/convert/remaining', { credentials: 'include' });
    if (!res.ok) return { remaining: 0 };
    return res.json();
  } catch {
    return { remaining: 0 };
  }
}

// ── Protected APIs ───────────────────────────────

export async function fetchMetrics(): Promise<string> {
  const res = await fetch(`${API_BASE}/metrics`, { credentials: 'include' });
  return res.text();
}
