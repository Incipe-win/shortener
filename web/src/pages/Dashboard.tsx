import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { motion } from 'framer-motion';
import { Link2, MousePointerClick, ShieldAlert, TrendingUp, Search, ExternalLink } from 'lucide-react';
import { Link } from '@tanstack/react-router';
import { Container } from '@/components/layout/Container';
import { Card } from '@/components/ui/Card';
import { Badge } from '@/components/ui/Badge';
import { GradientText } from '@/components/ui/GradientText';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import type { LinksResponse, LinkItem } from '@/types';

// Simulated fetch - will connect to real API once backend is ready
async function fetchLinks(page: number, search: string): Promise<LinksResponse> {
  const params = new URLSearchParams({ page: String(page), page_size: '10' });
  if (search) params.set('search', search);
  const res = await fetch(`/api/links?${params}`, { credentials: 'include' });
  if (!res.ok) throw new Error('Failed to fetch');
  return res.json();
}

function StatsCard({ icon: Icon, label, value, trend }: {
  icon: typeof Link2; label: string; value: string; trend?: string;
}) {
  return (
    <Card className="p-6">
      <div className="flex items-center justify-between mb-3">
        <div className="w-9 h-9 rounded-xl bg-[var(--color-accent)]/10 border border-[var(--color-border-hover)] flex items-center justify-center">
          <Icon className="w-4 h-4 text-[var(--color-accent)]" />
        </div>
        {trend && (
          <span className="text-xs text-[var(--color-safe)] flex items-center gap-1">
            <TrendingUp className="w-3 h-3" /> {trend}
          </span>
        )}
      </div>
      <p className="text-2xl font-semibold text-[var(--color-fg)]">{value}</p>
      <p className="text-xs text-[var(--color-fg-muted)] mt-1">{label}</p>
    </Card>
  );
}

function LinkRow({ link }: { link: LinkItem }) {
  return (
    <motion.tr
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="border-b border-[var(--color-border)] hover:bg-[var(--color-surface-hover)] transition-colors"
    >
      <td className="py-3 px-4">
        <Link
          to="/preview/$surl"
          params={{ surl: link.surl }}
          className="text-[var(--color-accent)] hover:text-[var(--color-accent-bright)] font-mono text-sm transition-colors"
        >
          /{link.surl}
        </Link>
      </td>
      <td className="py-3 px-4 max-w-xs truncate text-sm text-[var(--color-fg-muted)]">
        {link.lurl}
      </td>
      <td className="py-3 px-4">
        <Badge level={link.risk_level as 'safe' | 'warning' | 'danger' | 'pending'} />
      </td>
      <td className="py-3 px-4 text-sm text-[var(--color-fg)] font-mono">
        {link.click_count}
      </td>
      <td className="py-3 px-4 text-xs text-[var(--color-fg-muted)]">
        {new Date(link.create_at).toLocaleDateString('zh-CN')}
      </td>
      <td className="py-3 px-4">
        <button
          onClick={() => window.open(link.lurl, '_blank')}
          className="p-1.5 rounded-lg text-[var(--color-fg-muted)] hover:text-[var(--color-fg)] hover:bg-[var(--color-surface)] transition-colors"
        >
          <ExternalLink className="w-4 h-4" />
        </button>
      </td>
    </motion.tr>
  );
}

export function DashboardPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');

  const { data, isLoading } = useQuery({
    queryKey: ['links', page, search],
    queryFn: () => fetchLinks(page, search),
    retry: false,
  });

  return (
    <div className="pt-24 pb-16">
      <Container>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, ease: [0.16, 1, 0.3, 1] as const }}
        >
          <div className="mb-8">
            <p className="text-xs font-mono tracking-widest text-[var(--color-accent)] mb-2 uppercase">Dashboard</p>
            <h1 className="text-3xl font-semibold tracking-tight">
              <GradientText>链接仪表盘</GradientText>
            </h1>
          </div>

          {/* Stats Grid */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
            <StatsCard icon={Link2} label="总链接数" value={data ? String(data.total) : '—'} trend="+12%" />
            <StatsCard icon={MousePointerClick} label="今日点击" value="—" trend="+8%" />
            <StatsCard icon={ShieldAlert} label="安全拦截" value="—" />
            <StatsCard icon={TrendingUp} label="转化率" value="—" />
          </div>

          {/* Links Table */}
          <Card className="overflow-hidden">
            <div className="p-4 border-b border-[var(--color-border)] flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
              <h2 className="text-lg font-semibold text-[var(--color-fg)]">全部链接</h2>
              <div className="relative w-full sm:w-64">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-fg-subtle)]" />
                <Input
                  placeholder="搜索链接..."
                  value={search}
                  onChange={(e) => { setSearch(e.target.value); setPage(1); }}
                  className="pl-10 py-2"
                />
              </div>
            </div>

            <div className="overflow-x-auto">
              <table className="w-full text-left">
                <thead>
                  <tr className="border-b border-[var(--color-border)] text-xs text-[var(--color-fg-muted)] uppercase tracking-wider">
                    <th className="py-3 px-4 font-medium">短链</th>
                    <th className="py-3 px-4 font-medium">原始链接</th>
                    <th className="py-3 px-4 font-medium">风险</th>
                    <th className="py-3 px-4 font-medium">点击</th>
                    <th className="py-3 px-4 font-medium">创建时间</th>
                    <th className="py-3 px-4 font-medium"></th>
                  </tr>
                </thead>
                <tbody>
                  {isLoading ? (
                    <tr><td colSpan={6} className="py-12 text-center text-[var(--color-fg-muted)]">加载中...</td></tr>
                  ) : data?.list?.length ? (
                    data.list.map((link) => <LinkRow key={link.id} link={link} />)
                  ) : (
                    <tr><td colSpan={6} className="py-12 text-center text-[var(--color-fg-muted)]">
                      {search ? '无搜索结果' : '暂无链接数据，请先创建短链'}
                    </td></tr>
                  )}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {data && data.total > 10 && (
              <div className="p-4 border-t border-[var(--color-border)] flex items-center justify-between">
                <p className="text-xs text-[var(--color-fg-muted)]">
                  共 {data.total} 条 · 第 {page} 页
                </p>
                <div className="flex gap-2">
                  <Button variant="ghost" disabled={page <= 1} onClick={() => setPage(page - 1)}>上一页</Button>
                  <Button variant="ghost" disabled={page * 10 >= data.total} onClick={() => setPage(page + 1)}>下一页</Button>
                </div>
              </div>
            )}
          </Card>
        </motion.div>
      </Container>
    </div>
  );
}
