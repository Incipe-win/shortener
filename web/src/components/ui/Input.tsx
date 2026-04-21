import { forwardRef, type InputHTMLAttributes } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  error?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ className = '', error, ...rest }, ref) => {
    return (
      <div className="w-full">
        <input
          ref={ref}
          className={[
            'w-full px-4 py-3 rounded-lg text-sm',
            'bg-[var(--color-bg-input)] text-gray-100 placeholder:text-gray-500',
            'border border-[var(--color-border-hover)]',
            'focus:outline-none focus:border-[var(--color-accent)] focus:ring-2 focus:ring-[var(--color-accent)]/20',
            'transition-all duration-200 ease-[var(--ease-expo-out)]',
            error ? 'border-[var(--color-danger)]' : '',
            className,
          ].join(' ')}
          {...rest}
        />
        {error && (
          <p className="mt-1.5 text-xs text-[var(--color-danger)]">{error}</p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';
