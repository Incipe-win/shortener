import { useParams, Link } from '@tanstack/react-router';
import { motion } from 'framer-motion';
import { ArrowLeft, Link2Off } from 'lucide-react';
import { Container } from '@/components/layout/Container';
import { Card } from '@/components/ui/Card';
import { GradientText } from '@/components/ui/GradientText';
import { Button } from '@/components/ui/Button';

export function NotFoundPage() {
  const { surl } = useParams({ from: '/404/$surl' });

  return (
    <div className="pt-24 pb-16">
      <Container className="max-w-2xl">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, ease: [0.16, 1, 0.3, 1] as const }}
        >
          <Card className="p-8 text-center">
            <div className="w-16 h-16 rounded-2xl bg-[var(--color-fg-muted)]/10 border border-[var(--color-border)] flex items-center justify-center mx-auto mb-6">
              <Link2Off className="w-7 h-7 text-[var(--color-fg-muted)]" />
            </div>
            <p className="text-5xl font-bold text-[var(--color-fg-muted)] mb-2">404</p>
            <h1 className="text-2xl font-semibold tracking-tight mb-3">
              <GradientText>短链不存在</GradientText>
            </h1>
            <p className="text-sm text-[var(--color-fg-muted)] mb-2">
              该链接可能已被删除或从未创建。
            </p>
            {surl && (
              <p className="text-xs font-mono text-[var(--color-accent)] mb-6">/{surl}</p>
            )}
            <div className="mt-6">
              <Link to="/">
                <Button variant="secondary">
                  <ArrowLeft className="w-4 h-4" />
                  返回首页
                </Button>
              </Link>
            </div>
          </Card>
        </motion.div>
      </Container>
    </div>
  );
}
