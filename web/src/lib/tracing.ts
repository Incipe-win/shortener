import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-web';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch';
import { DocumentLoadInstrumentation } from '@opentelemetry/instrumentation-document-load';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { W3CTraceContextPropagator } from '@opentelemetry/core';

export function initTracing() {
  try {
    const exporter = new OTLPTraceExporter({
      url: 'http://localhost:4318/v1/traces',
    });

    const provider = new WebTracerProvider({
      spanProcessors: [new BatchSpanProcessor(exporter)],
    });

    provider.register({
      propagator: new W3CTraceContextPropagator(),
    });

    registerInstrumentations({
      instrumentations: [
        new FetchInstrumentation({
          propagateTraceHeaderCorsUrls: [/localhost:8888/, /\/api\//],
        }),
        new DocumentLoadInstrumentation(),
      ],
    });

    console.log('[OTel] Browser tracing initialized → Jaeger OTLP HTTP :4318');
  } catch (err) {
    console.warn('[OTel] Failed to initialize tracing:', err);
  }
}
