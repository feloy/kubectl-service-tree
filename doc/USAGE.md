
## Usage

The following assumes you have the plugin installed via

```shell
kubectl krew install service-tree
```

### Show ingresses, services and backends in current namespace

```shell
kubectl service-tree
```

### Show ingresses, services and backends in `demo` namespace

```shell
kubectl service-tree -n demo
```

## Notes

- Only services of type `NodePort` and `LoadBalancer` are displayed.

- When a Pod contains several containers and the containers declare their ports (with `ports.containerPort`), the name of the concerned container is displayed.
