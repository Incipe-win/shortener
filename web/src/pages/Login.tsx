import { useState } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { motion } from 'framer-motion';
import { LogIn, Lock, User } from 'lucide-react';
import { Container } from '@/components/layout/Container';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { GradientText } from '@/components/ui/GradientText';
import { useAuth } from '@/stores/auth';

const schema = z.object({
  username: z.string().min(1, '请输入用户名'),
  password: z.string().min(1, '请输入密码'),
});

type FormData = z.infer<typeof schema>;

export function LoginPage() {
  const navigate = useNavigate();
  const loginFn = useAuth((s) => s.login);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
  });

  const onSubmit = async (data: FormData) => {
    setLoading(true);
    setError('');
    try {
      await loginFn(data.username, data.password);
      navigate({ to: '/dashboard' });
    } catch {
      setError('用户名或密码错误');
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
              <Lock className="w-6 h-6 text-[var(--color-accent)]" />
            </div>
            <h1 className="text-2xl font-semibold tracking-tight mb-2">
              <GradientText>登录管理后台</GradientText>
            </h1>
            <p className="text-sm text-[var(--color-fg-muted)]">请输入凭证以访问仪表盘和监控面板</p>
          </div>

          <Card className="p-6 sm:p-8">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
              <div>
                <label className="block text-xs font-medium text-[var(--color-fg-muted)] mb-2">用户名</label>
                <div className="relative">
                  <User className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-fg-subtle)]" />
                  <Input
                    {...register('username')}
                    placeholder="admin"
                    error={errors.username?.message}
                    className="pl-10"
                  />
                </div>
              </div>

              <div>
                <label className="block text-xs font-medium text-[var(--color-fg-muted)] mb-2">密码</label>
                <div className="relative">
                  <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-fg-subtle)]" />
                  <Input
                    {...register('password')}
                    type="password"
                    placeholder="••••••••"
                    error={errors.password?.message}
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
          </Card>
        </motion.div>
      </Container>
    </div>
  );
}
