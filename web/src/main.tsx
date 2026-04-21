import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App'
import { initTracing } from './lib/tracing'
import { useAuth } from './stores/auth'

// Initialize OpenTelemetry tracing
initTracing()

// Check auth status on app load
useAuth.getState().checkAuth()

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
