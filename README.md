# Kubernetes Tutorial

This tutorial will demonostrate how to run a simple Go server that reads and writes
to a Postgres database. The application will run on the [MVP Studio](https://github.com/MVPStudio) 
Kubernetes cluster and take advantage of the Kubernetes secrets service.

Prerequisites
- Docker
- Kubernetes cluster on MVP Studio. [Instructions here](https://github.com/MVPStudio/k8)
- Go (for local development). [Installation instructions here](https://golang.org/doc/install)

---

## Kubectl commands
- create secret with
  ```
  kubectl create secret generic pgsecrets --from-literal=username=postgres --from-literal=password=password -n hello-world
  ```
- more info on secrets is in the [Kubernetes docs](https://kubernetes.io/docs/concepts/configuration/secret/)
- deploy with
  ```
  kubectl apply -f kube/
  ```
- shut down app with
  ```
  kubectl delete -f kube/
  ```
  