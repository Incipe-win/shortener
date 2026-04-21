import { useState, useEffect, useCallback } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { motion, AnimatePresence } from 'framer-motion';
import { ArrowRight, Copy, Check, Sparkles, AlertCircle } from 'lucide-react';
import { Link } from '@tanstack/react-router';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card } from '@/components/ui/Card';
import { convertUrl } from '@/lib/api';
import { useAuth } from '@/stores/auth';

const schema = z.object({
  long_url: z.string().url('请输入有效的 URL'),
});

type FormData = z.infer<typeof schema>;

export function ConvertForm() {
  const [result, setResult] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [remaining, setRemaining] = useState<number | null>(null);
  const { isAuthenticated } = useAuth();

  const fetchRemaining = useCallback(async () => {
    if (isAuthenticated) {
      setRemaining(-1); // unlimited
      return;
    }
    try {
      const res = await fetch('/api/convert/remaining', { credentials: 'include' });
      if (res.ok) {
        const data = await res.json();
        setRemaining(data.remaining);
      }
    } catch {
      // ignore
    }
  }, [isAuthenticated]);

  useEffect(() => {
    fetchRemaining();
  }, [fetchRemaining]);

  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
  });

  const onSubmit = async (data: FormData) => {
    setLoading(true);
    setResult(null);
    setError(null);
    try {
      const res = await convertUrl(data.long_url);
      setResult(res.short_url);
      fetchRemaining(); // refresh count
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      }
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = async () => {
    if (!result) return;
    await navigator.clipboard.writeText(result);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <Card className="p-6 sm:p-8 max-w-2xl mx-auto">
      {/* Limit Warning */}
      {!isAuthenticated && remaining !== null && remaining <= 3 && (
        <div className={`mb-4 p-3 rounded-lg flex items-center gap-2 text-sm ${remaining === 0 ? 'bg-red-500/10 text-red-400 border border-red-500/20' : 'bg-yellow-500/10 text-yellow-400 border border-yellow-500/20'}`}>
          <AlertCircle className="w-4 h-4 shrink-0" />
          <span>
            {remaining === 0
              ? '已达创建上限，请'
              : `未注册用户最多创建 3 个短链接（已用 ${3 - remaining}/3），`}
            <Link to="/login" className="underline hover:no-underline">登录/注册</Link>
            {remaining === 0 ? '后继续使用' : '后可无限制创建'}
          </span>
        </div>
      )}

      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
        <div className="flex flex-col sm:flex-row gap-3">
          <Input
            {...register('long_url')}
            placeholder="粘贴你的长链接..."
            error={errors.long_url?.message}
            className="flex-1"
          />
          <Button type="submit" loading={loading} className="shrink-0">
            <Sparkles className="w-4 h-4" />
            生成短链
            <ArrowRight className="w-4 h-4" />
          </Button>
        </div>
      </form>

      {error && (
        <div className="mt-4 p-3 rounded-lg bg-red-500/10 border border-red-500/20 text-sm text-red-400">
          {error}
        </div>
      )}

      <AnimatePresence>
        {result && (
          <motion.div
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -12 }}
            transition={{ duration: 0.3, ease: [0.16, 1, 0.3, 1] }}
            className="mt-6 p-4 rounded-xl bg-[var(--color-bg-elevated)] border border-[var(--color-border-accent)]"
          >
            <div className="flex items-center justify-between gap-3">
              <div className="min-w-0">
                <p className="text-xs font-mono tracking-widest text-[var(--color-fg-muted)] mb-1">短链接</p>
                <p className="text-lg font-semibold text-[var(--color-accent-bright)] truncate">{result}</p>
              </div>
              <Button variant="secondary" onClick={copyToClipboard} className="shrink-0">
                {copied ? <Check className="w-4 h-4 text-[var(--color-safe)]" /> : <Copy className="w-4 h-4" />}
                {copied ? '已复制' : '复制'}
              </Button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </Card>
  );
}
