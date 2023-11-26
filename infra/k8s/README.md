# k8s
## Summary
サンプルとしてminikubeで建てたk8s Cluster

## How to Launch
### 1. launch minikube 
minikube start --driver=docker

### 2. local の docker client を minikube の docker daemon に接続
```
eval $(minikube docker-env)
```

もどしたいとき
```
eval $(minikube docker-env -u)
```

### 3. build
```
make build
```

### 4. apply
```
kubectl apply -f manifest.yamol
```

## Others
### API Service の Port を Local に公開
```
minikube service mangahash-api-service --url
```

### DB Port-Forward
```
kubectl port-forward svc/mangahash-db-service 5432:5432
```

### DB 接続
```
PGOPTIONS="--search_path=app" psql -U mangahash-db-user -d mangahash-db -h localhost
```

### Migration from Local PC
```
goose -dir backend/migrations postgres "user=mangahash-db-user dbname=mangahash-db host=localhost password=tmppass sslmode=disable" up
```

### Restart All Pods
```
kubectl rollout restart deployments/<deployment name>
```

## 参考文献

https://www.tohoho-web.com/ex/kubernetes.html#about
https://redj.hatenablog.com/entry/2022/04/06/120010
