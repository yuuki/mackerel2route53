.PHONY: setup
setup:
	(cd stacks/route53; make create)
	(cd stacks/mackerel-webhook-gateway; make deploy)