VERSION := 1.0

help:
	go run app/services/fairsplit-api/main.go --help

all: fairsplit

fairsplit:
	docker build \
		-f zarf/docker/dockerfile.fairsplit-api \
		-t fairsplit-api:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within k8s/kind

KIND_CLUSTER := fairsplit-starter-cluster

dev-up:
	kind create cluster \
		--image kindest/node:v1.25.3@sha256:f52781bc0d7a19fb6c405c2af83abfeb311f130707a0e219175677e366cc45d1 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml
	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-load:
	kind load docker-image fairsplit-api:$(VERSION) --name $(KIND_CLUSTER)

dev-apply:
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

dev-describe:
	kubectl describe nodes
	kubectl describe svc

dev-describe-deployment:
	kubectl describe deployment --namespace=fairsplit-system fairsplit

dev-describe-fairsplit:
	kubectl describe pod --namespace=fairsplit-system -l app=fairsplit

dev-update: all dev-load dev-restart

dev-update-apply: all dev-load dev-apply
