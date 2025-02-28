For more details, please goto https://github.com/helm/charts/tree/master/stable/grafana

# Grafana Helm Chart

* Installs the web dashboarding system [Grafana](http://grafana.org/)

## TL;DR;

```console
$ helm install stable/grafana --name grafana --namespace monitoring -f values.yaml
```

Compared to the upstream `helm/charts/stable/grafana/values.yaml`, the supplied `values.yaml` in this folder adds two datasource configurations:  
```
datasources:
  datasources.yaml:
    apiVersion: 1
    datasources:
    - name: Prometheus
      type: prometheus
      url: http://prometheus-server
      access: proxy
      isDefault: true
    - name: InfluxDB
      type: influxdb
      database: _internal
      url: http://influxdb-influxdb:8086
      access: proxy
      isDefault: false
```
