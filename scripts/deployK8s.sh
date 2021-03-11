gcloud container clusters get-credentials icsharing --region asia-east1-b --project red-atlas-303101
kubectl apply -f kubernetes/env-configmap.yaml
kubectl apply -f kubernetes/env-secrets.yaml
kubectl apply -f kubernetes/mongo.yaml
kubectl apply -f kubernetes/redis.yaml
kubectl apply -f kubernetes/icsharing.yaml