import { motion } from 'framer-motion';
import { Container } from '@/components/layout/Container';
import { GradientText } from '@/components/ui/GradientText';
import { ConvertForm } from '@/components/features/ConvertForm';
import { Link2, Shield, Zap } from 'lucide-react';
import { Card } from '@/components/ui/Card';

const features = [
  {
    icon: Link2,
    title: '智能转链',
    desc: '输入长链接，秒级生成短链。内置 AI 语义分析与安全检测。',
  },
  {
    icon: Shield,
    title: '安全巡检',
    desc: '多层安全机制：域名黑名单、LLM 风险评估、Kafka 实时告警。',
  },
  {
    icon: Zap,
    title: '全链可观测',
    desc: 'OpenTelemetry + Jaeger 端到端追踪，Prometheus 全指标监控。',
  },
];

const fadeUp = {
  hidden: { opacity: 0, y: 24 },
  visible: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: { delay: i * 0.08, duration: 0.6, ease: [0.16, 1, 0.3, 1] as const },
  }),
};

export function HomePage() {
  return (
    <div className="pt-24">
      {/* Hero */}
      <section className="py-16 sm:py-24 lg:py-32">
        <Container className="text-center">
          <motion.div
            initial="hidden"
            animate="visible"
            variants={fadeUp}
            custom={0}
          >
            <p className="text-xs font-mono tracking-widest text-[var(--color-accent)] mb-6 uppercase">
              Smart-Shortener Gateway
            </p>
          </motion.div>

          <motion.h1
            initial="hidden"
            animate="visible"
            variants={fadeUp}
            custom={1}
            className="text-4xl sm:text-5xl lg:text-7xl font-semibold tracking-[-0.03em] leading-tight mb-6"
          >
            <GradientText>让每一个链接</GradientText>
            <br />
            <GradientText accent>更智能、更安全</GradientText>
          </motion.h1>

          <motion.p
            initial="hidden"
            animate="visible"
            variants={fadeUp}
            custom={2}
            className="text-base sm:text-lg lg:text-xl text-[var(--color-fg-muted)] max-w-2xl mx-auto mb-12 leading-relaxed"
          >
            基于微服务与大模型的智能短链安全网关。
            AI 内容分析、零信任安全巡检、Kafka 事件流、全链路可观测。
          </motion.p>

          <motion.div
            initial="hidden"
            animate="visible"
            variants={fadeUp}
            custom={3}
          >
            <ConvertForm />
          </motion.div>
        </Container>
      </section>

      {/* Features */}
      <section className="py-16 sm:py-24 border-t border-[var(--color-border)]">
        <Container>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {features.map((f, i) => (
              <motion.div
                key={f.title}
                initial="hidden"
                whileInView="visible"
                viewport={{ once: true, amount: 0.2 }}
                variants={fadeUp}
                custom={i}
              >
                <Card className="p-6 h-full hover:-translate-y-1 transition-transform duration-300 ease-[var(--ease-expo-out)]">
                  <div className="w-10 h-10 rounded-xl bg-[var(--color-accent)]/10 border border-[var(--color-border-hover)] flex items-center justify-center mb-4">
                    <f.icon className="w-5 h-5 text-[var(--color-accent)]" />
                  </div>
                  <h3 className="text-lg font-semibold text-[var(--color-fg)] mb-2">{f.title}</h3>
                  <p className="text-sm text-[var(--color-fg-muted)] leading-relaxed">{f.desc}</p>
                </Card>
              </motion.div>
            ))}
          </div>
        </Container>
      </section>
    </div>
  );
}
