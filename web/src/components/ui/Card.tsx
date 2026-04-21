import { useRef, useState, type ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  className?: string;
  spotlight?: boolean;
}

export function Card({ children, className = '', spotlight = true }: CardProps) {
  const ref = useRef<HTMLDivElement>(null);
  const [mousePos, setMousePos] = useState({ x: 0, y: 0 });
  const [isHovered, setIsHovered] = useState(false);

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!ref.current || !spotlight) return;
    const rect = ref.current.getBoundingClientRect();
    setMousePos({ x: e.clientX - rect.left, y: e.clientY - rect.top });
  };

  return (
    <div
      ref={ref}
      onMouseMove={handleMouseMove}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      className={[
        'relative overflow-hidden rounded-2xl',
        'bg-gradient-to-b from-white/[0.08] to-white/[0.02]',
        'border border-[var(--color-border)]',
        'shadow-[var(--shadow-card)]',
        'hover:shadow-[var(--shadow-card-hover)]',
        'hover:border-[var(--color-border-hover)]',
        'transition-all duration-300 ease-[var(--ease-expo-out)]',
        className,
      ].join(' ')}
    >
      {/* Top edge glow line */}
      <div className="absolute top-0 inset-x-0 h-px bg-gradient-to-r from-transparent via-white/20 to-transparent" />

      {/* Mouse-tracking spotlight */}
      {spotlight && (
        <div
          className="absolute inset-0 pointer-events-none transition-opacity duration-300"
          style={{
            opacity: isHovered ? 1 : 0,
            background: `radial-gradient(300px circle at ${mousePos.x}px ${mousePos.y}px, var(--color-accent-glow), transparent 70%)`,
          }}
        />
      )}

      {/* Content */}
      <div className="relative z-10">
        {children}
      </div>
    </div>
  );
}
