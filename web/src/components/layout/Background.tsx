export function Background({ children }: { children: React.ReactNode }) {
  return (
    <div className="relative min-h-screen bg-[radial-gradient(ellipse_at_top,#0a0a0f_0%,#050506_50%,#020203_100%)] noise-overlay grid-overlay">
      {/* Animated gradient blobs */}
      <div className="fixed inset-0 z-0 overflow-hidden pointer-events-none" aria-hidden="true">
        {/* Primary blob — top center */}
        <div
          className="absolute top-[-200px] left-1/2 -translate-x-1/2 w-[900px] h-[1400px] rounded-full opacity-25 blur-[150px]"
          style={{
            background: 'radial-gradient(circle, var(--color-accent) 0%, transparent 70%)',
            animation: 'float 10s ease-in-out infinite',
          }}
        />
        {/* Secondary blob — left */}
        <div
          className="absolute top-[20%] left-[-100px] w-[600px] h-[800px] rounded-full opacity-15 blur-[120px]"
          style={{
            background: 'radial-gradient(circle, #7c3aed 0%, #ec4899 50%, transparent 70%)',
            animation: 'float-reverse 8s ease-in-out infinite',
          }}
        />
        {/* Tertiary blob — right */}
        <div
          className="absolute top-[40%] right-[-100px] w-[500px] h-[700px] rounded-full opacity-12 blur-[100px]"
          style={{
            background: 'radial-gradient(circle, #4f46e5 0%, #3b82f6 50%, transparent 70%)',
            animation: 'float 9s ease-in-out infinite 2s',
          }}
        />
        {/* Bottom pulse */}
        <div
          className="absolute bottom-[-200px] left-1/2 -translate-x-1/2 w-[800px] h-[600px] rounded-full blur-[120px]"
          style={{
            background: 'radial-gradient(circle, var(--color-accent) 0%, transparent 70%)',
            animation: 'pulse-glow 6s ease-in-out infinite',
          }}
        />
      </div>

      {/* Content layer */}
      <div className="relative z-10">
        {children}
      </div>
    </div>
  );
}
