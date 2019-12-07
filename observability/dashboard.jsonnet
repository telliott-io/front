local grafana = import 'grafonnet/grafana.libsonnet';

grafana.dashboard.new(
  'Front',
  schemaVersion=16,
)
.addPanel(
    grafana.graphPanel.new(
        'Index Requests',
        format='short',
        datasource='Prometheus',
        span=2,
    )
    .addTarget(
        grafana.prometheus.target('sum by (job)(rate(front_serve_index_count[1m]))')
    ), gridPos={
    x: 0,
    y: 0,
    w: 24,
    h: 10,
  }
)
.addPanel(
    grafana.graphPanel.new(
        'Index Latency',
        format='ms',
        datasource='Prometheus',
        span=2,
    )
    .addTarget(
        grafana.prometheus.target('histogram_quantile(0.95, sum(rate(get_projects_duration_ms_bucket[5m])) by (le))')
    ), gridPos={
    x: 0,
    y: 0,
    w: 24,
    h: 10,
  }
)