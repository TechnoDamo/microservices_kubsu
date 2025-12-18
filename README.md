# microservices_kubsu
NodePort type
ClusterIP type
ConfigMap
Namespace
labels
Liveness and Readiness Probes
Kubernetes Deployment
CronJobs/schedulers in kubernetess
Helm Chart

microservices:
1. auth-service
    1. POST /v1/register - creating new user
    request:
        login string 
        password string
    response:
        token string
    2. POST /v1/login - loging new user
    request:
        login string
        password string
    response:
        token string
    3. POST /v1/logout - logout user
    request:
        login string
        token string
    response:
        empty
    4. GET /v1/validate/token
    request:
        token string
    response:
        user_id string

2. profile-service
3. notification-service
    POST /v1/send/notification/{user_id}
    request:
        message string
    response: 
        empty
    GET /v1/notifications
    request:
        user_id
        token
    response
        message

4. report-service

helm repo add bitnami https://charts.bitnami.com/bitnami[citation:1]
kubectl create configmap postgres-init-script --from-file=init-db.sql

helm install kubsu-project-db bitnami/postgresql -f postgres-values.yaml

kubectl get pods  # Should show "Running"
kubectl logs observability-db-postgresql-0 | grep "PostgreSQL init process complete"

kubectl delete pod kubsu-project-db
helm uninstall kubsu-project-db
kubectl delete pvc data-kubsu-project-db-postgresql-0
kubsu-project-db-postgresql.default.svc.cluster.local

helm uninstall kubsu-project-db
kubectl delete pvc data-kubsu-project-db-postgresql-0
kubectl delete configmap postgres-init-script
kubectl create configmap postgres-init-script --from-file=init-db.sql
helm install kubsu-project-db bitnami/postgresql -f postgres-values.yaml
minikube service kubsu-project-db-postgresql --url

#deletion
# Uninstall the Helm release (removes StatefulSet, Services, etc.)
helm uninstall kubsu-project-db
# Delete the Persistent Volume Claim (WARNING: This PERMANENTLY deletes all database data)
kubectl delete pvc data-kubsu-project-db-postgresql-0
# Delete your custom initialization ConfigMap
kubectl delete configmap postgres-init-script

#recreatioin
# Create your ConfigMap from the current init-db.sql file
kubectl create configmap postgres-init-script --from-file=init-db.sql
# Install PostgreSQL with your configuration
helm install kubsu-project-db bitnami/postgresql -f postgres-values.yaml