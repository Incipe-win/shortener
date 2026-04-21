import { useState } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { motion } from 'framer-motion';
import { LogIn, Lock, User, UserPlus } from 'lucide-react';
import { Container } from '@/components/layout/Container';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { GradientText } from '@/components/ui/GradientText';
import { useAuth } from '@/stores/auth';

const loginSchema = z.object({
  username: z.string().min(1, '请输入用户名'),
  password: z.string().min(1, '请输入密码'),
});

const registerSchema = z.object({
  username: z.string().min(3, '用户名至少 3 位').max(32).regex(/^[a-zA-Z0-9]+$/, '仅支持字母数字'),
  password: z.string().min(6, '密码至少 6 位'),
  confirmPassword: z.string().min(1, '请确认密码'),
}).refine((data) => data.password === data.confirmPassword, {
  message: '两次密码不一致',
  path: ['confirmPassword'],
});

type LoginFormData = z.infer<typeof loginSchema>;
type RegisterFormData = z.infer<typeof registerSchema>;

export function LoginPage() {
  const navigate = useNavigate();
  const { login, register: registerFn } = useAuth();
  const [mode, setMode] = useState<'login' | 'register'>('login');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const loginForm = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });

  const registerForm = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  });

  const onLogin = async (data: LoginFormData) => {
    setLoading(true);
    setError('');
    try {
      await login(data.username, data.password);
      navigate({ to: '/dashboard' });
    } catch {
      setError('用户名或密码错误');
    } finally {
      setLoading(false);
    }
  };

  const onRegister = async (data: RegisterFormData) => {
    setLoading(true);
    setError('');
    try {
      await registerFn(data.username, data.password);
      navigate({ to: '/dashboard' });
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : '注册失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="pt-24 pb-16 min-h-screen flex items-center">
      <Container className="max-w-md">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, ease: [0.16, 1, 0.3, 1] as const }}
        >
          <div className="text-center mb-8">
            <div className="w-14 h-14 rounded-2xl bg-[var(--color-accent)]/10 border border-[var(--color-border-hover)] flex items-center justify-center mx-auto mb-6">
              {mode === 'login' ? <Lock className="w-6 h-6 text-[var(--color-accent)]" /> : <UserPlus className="w-6 h-6 text-[var(--color-accent)]" />}
            </div>
            <h1 className="text-2xl font-semibold tracking-tight mb-2">
              <GradientText>{mode === 'login' ? '登录管理后台' : '注册账号'}</GradientText>
            </h1>
            <p className="text-sm text-[var(--color-fg-muted)]">
              {mode === 'login' ? '请输入凭证以访问仪表盘和监控面板' : '注册后可无限制创建短链接'}
            </p>
          </div>

          {/* Mode Toggle */}
          <div className="flex rounded-xl overflow-hidden border border-[var(--color-border)] mb-6">
            <button
              className={`flex-1 py-2.5 text-sm font-medium transition-colors ${mode === 'login' ? 'bg-[var(--color-accent)]/10 text-[var(--color-accent)]' : 'text-[var(--color-fg-muted)] hover:bg-[var(--color-surface)]'}`}
              onClick={() => { setMode('login'); setError(''); }}
            >
              <LogIn className="w-4 h-4 inline mr-1" />
              登录
            </button>
            <button
              className={`flex-1 py-2.5 text-sm font-medium transition-colors ${mode === 'register' ? 'bg-[var(--color-accent)]/10 text-[var(--color-accent)]' : 'text-[var(--color-fg-muted)] hover:bg-[var(--color-surface)]'}`}
              onClick={() => { setMode('register'); setError(''); }}
            >
              <UserPlus className="w-4 h-4 inline mr-1" />
              注册
            </button>
          </div>

          <Card className="p-6 sm:p-8">
            {mode === 'login' ? (
              <form onSubmit={loginForm.handleSubmit(onLogin)} className="space-y-5">
                <div>
                  <label className="block text-xs font-medium text-[var(--color-fg-muted)] mb-2">用户名</label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-fg-subtle)]" />
                    <Input
                      {...loginForm.register('username')}
                      placeholder="admin"
                      error={loginForm.formState.errors.username?.message}
                      className="pl-10"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-xs font-medium text-[var(--color-fg-muted)] mb-2">密码</label>
                  <div className="relative">
                    <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-fg-subtle)]" />
                    <Input
                      {...loginForm.register('password')}
                      type="password"
                      placeholder="••••••••"
                      error={loginForm.formState.errors.password?.message}
                      className="pl-10"
                    />
                  </div>
                </div>

                {error && (
                  <p className="text-sm text-[var(--color-danger)] text-center">{error}</p>
                )}

                <Button type="submit" loading={loading} className="w-full">
                  <LogIn className="w-4 h-4" />
                  登录
                </Button>
              </form>
            ) : (
              <form onSubmit={registerForm.handleSubmit(onRegister)} className="space-y-5">
                <div>
                  <label className="block text-xs font-medium text-[var(--color-fg-muted)] mb-2">用户名</label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-fg-subtle)]" />
                    <Input
                      {...registerForm.register('username')}
                      placeholder="3-32 位字母数字"
                      error={registerForm.formState.errors.username?.message}
                      className="pl-10"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-xs font-medium text-[var(--color-fg-muted)] mb-2">密码</label>
                  <div className="relative">
                    <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-fg-subtle)]" />
                    <Input
                      {...registerForm.register('password')}
                      type="password"
                      placeholder="至少 6 位"
                      error={registerForm.formState.errors.password?.message}
                      className="pl-10"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-xs font-medium text-[var(--color-fg-muted)] mb-2">确认密码</label>
                  <div className="relative">
                    <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-fg-subtle)]" />
                    <Input
                      {...registerForm.register('confirmPassword')}
                      type="password"
                      placeholder="再次输入密码"
                      error={registerForm.formState.errors.confirmPassword?.message}
                      className="pl-10"
                    />
                  </div>
                </div>

                {error && (
                  <p className="text-sm text-[var(--color-danger)] text-center">{error}</p>
                )}

                <Button type="submit" loading={loading} className="w-full">
                  <UserPlus className="w-4 h-4" />
                  注册
                </Button>
              </form>
            )}
          </Card>
        </motion.div>
      </Container>
    </div>
  );
}
