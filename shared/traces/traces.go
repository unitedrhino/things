package traces

import (
	"context"
	ztrace "github.com/zeromicro/go-zero/core/trace"
	otletrace "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

/*链路追踪接口，满足opentracing规范，通过ctx 将链路tid和spanid传递*/
/*使用方法：
  将上游ctx作为参数传入，tp为yaml中配置的jaeger的信息的全局对象，
  instrumentationName和spanName可自由命名，会在jaeger可视化中显示出来，后续定制命名规范。
  返回本节点的ctx和span信息，以备下游继续传递使用

示例：
		1、
			dd.yaml中配置：
			Telemetry:
			  Name: dd.rpc                                --service name
			  Endpoint: http://127.0.0.1:14268/api/traces --推送jaeger的已远程url地址
			  Sampler: 1.0                                --默认1.0
			  Batcher: jaeger                             --可选择 jaeger或者zipkin
		2、
            ctx_sub0 := context.Background()
			ctx_sub1, span1 := tr.StartSpan(ctx_sub0, "fun1", "target1")
			defer span1.End() -- 接口结束前必须调用span.End() 否则无法将本节点链路信息推送jaeger
			logx.Infof("traceid:%s, spanid:%s", span1.SpanContext().TraceID(), span1.SpanContext().SpanID())

----------- 再将上游的ctx_sub1 传递到下游节点形成父子关系：

            ctx_sub2, span2 := tr.StartSpan(ctx_sub1, "fun2", "target2")
			defer span2.End()
			logx.Infof("traceid:%s, spanid:%s", span2.SpanContext().TraceID(), span2.SpanContext().SpanID())
*/
func StartSpan(ctx context.Context, method, target string) (context.Context, trace.Span) {
	name, attr := ztrace.SpanInfo(method, target)
	tr := otletrace.GetTracerProvider().Tracer(ztrace.TraceName)
	return tr.Start(ctx, name, trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attr...))
}
