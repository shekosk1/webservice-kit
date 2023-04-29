VERSION := 1.0

all: fairsplit

fairsplit:
	docker build \
		-f zarf/docker/dockerfile.fairsplit-api \
		-t fairsplit-api:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.