import { useState, useEffect, useCallback } from 'react';
import { motion } from 'framer-motion';
import { Activity, Cpu, MessageSquare, ShieldAlert, RefreshCw } from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { Container } from '@/components/layout/Container';
import { Card } from '@/components/ui/Card';
import { GradientText } from '@/components/ui/GradientText';
import { Button } from '@/components/ui/Button';
import { fetchMetrics } from '@/lib/api';

interface MetricPoint {
  time: string;
  requests: number;
  kafkaProduce: number;
  kafkaConsume: number;
}

function parsePrometheusValue(text: string, metricName: string): number {
  const regex = new RegExp(`^${metricName}(?:\\{[^}]*\\})?\\s+(\\d+\\.?\\d*)`, 'm');
  const match = text.match(regex);
  return match ? parseFloat(match[1]) : 0;
}

function MetricCard({ icon: Icon, label, value, subtitle }: {
  icon: typeof Activity; label: string; value: string; subtitle?: string;
}) {
  return (
    <Card className="p-5">
      <div className="flex items-center gap-3 mb-2">
        <div className="w-8 h-8 rounded-lg bg-[var(--color-accent)]/10 border border-[var(--color-border-hover)] flex items-center justify-center">
          <Icon className="w-4 h-4 text-[var(--color-accent)]" />
        </div>
        <span className="text-xs text-[var(--color-fg-muted)] uppercase tracking-wider">{label}</span>
      </div>
      <p className="text-xl font-semibold font-mono text-[var(--color-fg)]">{value}</p>
      {subtitle && <p className="text-xs text-[var(--color-fg-muted)] mt-1">{subtitle}</p>}
    </Card>
  );
}

const chartStyle = {
  grid: { stroke: 'rgba(255,255,255,0.04)' },
  axis: { stroke: 'rgba(255,255,255,0.1)', fontSize: 11, fill: '#8A8F98' },
  tooltip: { backgroundColor: '#0a0a0c', border: '1px solid rgba(255,255,255,0.06)', borderRadius: '8px' },
};

export function MonitorPage() {
  const [history, setHistory] = useState<MetricPoint[]>([]);
  const [currentMetrics, setCurrentMetrics] = useState({ requests: 0, kafkaProduce: 0, kafkaConsume: 0, blocked: 0 });
  const [refreshing, setRefreshing] = useState(false);

  const pollMetrics = useCallback(async () => {
    try {
      setRefreshing(true);
      const text = await fetchMetrics();
      const requests = parsePrometheusValue(text, 'shortener_convert_total');
      const kafkaProduce = parsePrometheusValue(text, 'shortener_kafka_produce_total');
      const kafkaConsume = parsePrometheusValue(text, 'shortener_kafka_consume_total');
      const blocked = parsePrometheusValue(text, 'shortener_safety_blocked_total');

      setCurrentMetrics({ requests, kafkaProduce, kafkaConsume, blocked });
      setHistory((prev) => {
        const next = [...prev, { time: new Date().toLocaleTimeString('zh-CN', { hour12: false }), requests, kafkaProduce, kafkaConsume }];
        return next.slice(-30); // Keep last 30 data points
      });
    } catch {
      // Metrics endpoint may not be available
    } finally {
      setRefreshing(false);
    }
  }, []);

  useEffect(() => {
    pollMetrics();
    const interval = setInterval(pollMetrics, 10_000);
    return () => clearInterval(interval);
  }, [pollMetrics]);

  return (
    <div className="pt-24 pb-16">
      <Container>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, ease: [0.16, 1, 0.3, 1] as const }}
        >
          <div className="flex items-center justify-between mb-8">
            <div>
              <p className="text-xs font-mono tracking-widest text-[var(--color-accent)] mb-2 uppercase">Monitor</p>
              <h1 className="text-3xl font-semibold tracking-tight">
                <GradientText>实时监控</GradientText>
              </h1>
            </div>
            <Button variant="secondary" onClick={pollMetrics} loading={refreshing}>
              <RefreshCw className="w-4 h-4" />
              刷新
            </Button>
          </div>

          {/* Metrics Cards */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
            <MetricCard icon={Activity} label="转链请求" value={String(currentMetrics.requests)} subtitle="总计" />
            <MetricCard icon={MessageSquare} label="Kafka 生产" value={String(currentMetrics.kafkaProduce)} subtitle="消息数" />
            <MetricCard icon={Cpu} label="Kafka 消费" value={String(currentMetrics.kafkaConsume)} subtitle="已处理" />
            <MetricCard icon={ShieldAlert} label="安全拦截" value={String(currentMetrics.blocked)} subtitle="总计" />
          </div>

          {/* Chart */}
          <Card className="p-6">
            <h2 className="text-lg font-semibold text-[var(--color-fg)] mb-6">趋势（每10秒刷新）</h2>
            <div className="h-[300px]">
              {history.length > 1 ? (
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={history}>
                    <CartesianGrid strokeDasharray="3 3" {...chartStyle.grid} />
                    <XAxis dataKey="time" {...chartStyle.axis} />
                    <YAxis {...chartStyle.axis} />
                    <Tooltip
                      contentStyle={chartStyle.tooltip}
                      labelStyle={{ color: '#EDEDEF' }}
                      itemStyle={{ color: '#8A8F98' }}
                    />
                    <Line type="monotone" dataKey="requests" name="转链请求" stroke="var(--color-accent)" strokeWidth={2} dot={false} />
                    <Line type="monotone" dataKey="kafkaProduce" name="Kafka 生产" stroke="#7c3aed" strokeWidth={2} dot={false} />
                    <Line type="monotone" dataKey="kafkaConsume" name="Kafka 消费" stroke="var(--color-safe)" strokeWidth={2} dot={false} />
                  </LineChart>
                </ResponsiveContainer>
              ) : (
                <div className="h-full flex items-center justify-center text-[var(--color-fg-muted)]">
                  数据采集中，请稍候...
                </div>
              )}
            </div>
          </Card>
        </motion.div>
      </Container>
    </div>
  );
}
