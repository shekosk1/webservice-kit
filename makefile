# Hashicorp Vault
# 	READ THIS: https://developer.hashicorp.com/vault/docs/concepts/tokens
# 	$ export VAULT_TOKEN=mytoken
# 	$ export VAULT_ADDR='http://vault-service.fairsplit-system.svc.cluster.local:8200'
# 	$ vault secrets list
# 	$ vault kv get secret/fairsplit
# 	$ vault kv put secret/fairsplit key="some data"
# 	$ kubectl logs --namespace=fairsplit-system -l app=fairsplit -c init-vault-server
# 	$ curl -H "X-Vault-Token: mytoken" -X GET http://vault-service.fairsplit-system.svc.cluster.local:8200/v1/secret/data/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1
# 	$ curl -H "X-Vault-Token: mytoken" -H "Content-Type: application/json" -X POST -d '{"data":{"pk":"PEM"}}' http://vault-service.fairsplit-system.svc.cluster.local:8200/v1/secret/data/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1
#

VERSION 	 := 1.0
# ==============================================================================
# Build containers

help:
	go run app/services/fairsplit-api/main.go --help  | go run app/tooling/logfmt/main.go

all: fairsplit

fairsplit:
	docker build \
		-f zarf/docker/dockerfile.fairsplit-api \
		-t fairsplit-api:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Run from within k8s/kind
# http://fairsplit-service.fairsplit-system.svc.cluster.local:4000
# http://fairsplit-service.fairsplit-system.svc.cluster.local:3000

KIND_CLUSTER := fairsplit-starter-cluster

GOLANG       := golang:1.19
ALPINE       := alpine:3.17
KIND         := kindest/node:v1.25.3
POSTGRES     := postgres:15-alpine
VAULT        := hashicorp/vault:1.13
ZIPKIN       := openzipkin/zipkin:2.23
TELEPRESENCE := docker.io/datawire/tel2:2.10.4

dev-up:
	kind create cluster \
		--image kindest/node:v1.25.3@sha256:f52781bc0d7a19fb6c405c2af83abfeb311f130707a0e219175677e366cc45d1 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml
	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner
	
	kind load docker-image $(TELEPRESENCE) --name $(KIND_CLUSTER)
	kind load docker-image $(VAULT) --name $(KIND_CLUSTER)

	
	telepresence --context=kind-$(KIND_CLUSTER) helm install
#   telepresence --context=kind-$(KIND_CLUSTER) connect
	sudo telepresence --context=kind-$(KIND_CLUSTER) connect 

dev-down:
	telepresence quit -s
	kind delete cluster --name $(KIND_CLUSTER)

dev-load:
	cd zarf/k8s/dev/fairsplit; kustomize edit set image fairsplit-api-image=fairsplit-api:$(VERSION)
	kind load docker-image fairsplit-api:$(VERSION) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/vault | kubectl apply -f -

	kustomize build zarf/k8s/dev/fairsplit | kubectl apply -f -
	kubectl wait --timeout=120s --namespace=fairsplit-system --for=condition=Available deployment/fairsplit

dev-restart:
	kubectl rollout restart deployment fairsplit --namespace=fairsplit-system

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-logs:
	kubectl logs --namespace=fairsplit-system -l app=fairsplit --all-containers=true -f --tail=100 --max-log-requests=6

dev-logs-fmt:
	kubectl logs --namespace=fairsplit-system -l app=fairsplit --all-containers=true -f --tail=100 --max-log-requests=6 | go run app/tooling/logfmt/main.go -service=FAIRSPLIT-API

dev-describe:
	kubectl describe nodes
	kubectl describe svc

dev-describe-deployment:
	kubectl describe deployment --namespace=fairsplit-system fairsplit

dev-describe-fairsplit:
	kubectl describe pod --namespace=fairsplit-system -l app=fairsplit

dev-describe-tel:
	kubectl describe pod --namespace=ambassador -l app=traffic-manager

dev-update: all dev-load dev-restart

dev-update-apply: all dev-load dev-apply

test-load:
	hey -m GET -c 100 -n 10000 http://fairsplit-service.fairsplit-system.svc.cluster.local:3000/status

test-load-e:
	hey -m GET -c 100 -n 10000 http://fairsplit-service.fairsplit-system.svc.cluster.local:3000/empty

test-load-p:
	hey -m GET -c 100 -n 10000 http://fairsplit-service.fairsplit-system.svc.cluster.local:3000/testpanic