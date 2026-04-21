import { useParams } from '@tanstack/react-router';
import { useQuery } from '@tanstack/react-query';
import { motion } from 'framer-motion';
import { ExternalLink, Shield, Brain, Tag, Loader2 } from 'lucide-react';
import { Container } from '@/components/layout/Container';
import { Card } from '@/components/ui/Card';
import { Badge } from '@/components/ui/Badge';
import { GradientText } from '@/components/ui/GradientText';
import { Button } from '@/components/ui/Button';
import { previewLink } from '@/lib/api';

export function PreviewPage() {
  const { surl } = useParams({ from: '/preview/$surl' });

  const { data, isLoading, error } = useQuery({
    queryKey: ['preview', surl],
    queryFn: () => previewLink(surl),
  });

  return (
    <div className="pt-24 pb-16">
      <Container className="max-w-3xl">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, ease: [0.16, 1, 0.3, 1] as const }}
        >
          <p className="text-xs font-mono tracking-widest text-[var(--color-accent)] mb-4 uppercase">
            链接预览
          </p>
          <h1 className="text-3xl sm:text-4xl font-semibold tracking-tight mb-8">
            <GradientText>/{surl}</GradientText>
          </h1>

          {isLoading && (
            <div className="flex items-center justify-center py-20 text-[var(--color-fg-muted)]">
              <Loader2 className="w-6 h-6 animate-spin mr-3" />
              加载中...
            </div>
          )}

          {error && (
            <Card className="p-8 text-center">
              <p className="text-[var(--color-danger)]">加载失败：{(error as Error).message}</p>
            </Card>
          )}

          {data && (
            <div className="space-y-6">
              {/* Risk Level */}
              <Card className="p-6">
                <div className="flex items-center gap-3 mb-4">
                  <div className="w-9 h-9 rounded-xl bg-[var(--color-accent)]/10 border border-[var(--color-border-hover)] flex items-center justify-center">
                    <Shield className="w-4 h-4 text-[var(--color-accent)]" />
                  </div>
                  <h2 className="text-lg font-semibold text-[var(--color-fg)]">安全评估</h2>
                </div>
                <Badge level={data.risk_level as 'safe' | 'warning' | 'danger' | 'pending'} />
              </Card>

              {/* AI Summary */}
              <Card className="p-6">
                <div className="flex items-center gap-3 mb-4">
                  <div className="w-9 h-9 rounded-xl bg-[var(--color-accent)]/10 border border-[var(--color-border-hover)] flex items-center justify-center">
                    <Brain className="w-4 h-4 text-[var(--color-accent)]" />
                  </div>
                  <h2 className="text-lg font-semibold text-[var(--color-fg)]">AI 摘要</h2>
                </div>
                <p className="text-[var(--color-fg-muted)] leading-relaxed">
                  {data.summary || '暂无摘要'}
                </p>
              </Card>

              {/* Keywords */}
              {data.keywords && data.keywords.length > 0 && (
                <Card className="p-6">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="w-9 h-9 rounded-xl bg-[var(--color-accent)]/10 border border-[var(--color-border-hover)] flex items-center justify-center">
                      <Tag className="w-4 h-4 text-[var(--color-accent)]" />
                    </div>
                    <h2 className="text-lg font-semibold text-[var(--color-fg)]">关键词</h2>
                  </div>
                  <div className="flex flex-wrap gap-2">
                    {data.keywords.map((kw) => (
                      <span
                        key={kw}
                        className="px-3 py-1 rounded-full text-xs bg-[var(--color-surface)] text-[var(--color-fg-muted)] border border-[var(--color-border)]"
                      >
                        {kw}
                      </span>
                    ))}
                  </div>
                </Card>
              )}

              {/* Original URL */}
              <Card className="p-6">
                <p className="text-xs font-mono tracking-widest text-[var(--color-fg-muted)] mb-2 uppercase">原始链接</p>
                <p className="text-sm text-[var(--color-fg)] break-all mb-4">{data.long_url}</p>
                <Button
                  variant="secondary"
                  onClick={() => window.open(data.long_url, '_blank')}
                >
                  <ExternalLink className="w-4 h-4" />
                  访问原始链接
                </Button>
              </Card>
            </div>
          )}
        </motion.div>
      </Container>
    </div>
  );
}
