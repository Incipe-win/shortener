import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { motion, AnimatePresence } from 'framer-motion';
import { ArrowRight, Copy, Check, Sparkles } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card } from '@/components/ui/Card';
import { convertUrl } from '@/lib/api';

const schema = z.object({
  long_url: z.string().url('请输入有效的 URL'),
});

type FormData = z.infer<typeof schema>;

export function ConvertForm() {
  const [result, setResult] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [loading, setLoading] = useState(false);

  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
  });

  const onSubmit = async (data: FormData) => {
    setLoading(true);
    setResult(null);
    try {
      const res = await convertUrl(data.long_url);
      setResult(res.short_url);
    } catch (err) {
      console.error(err);
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
