migrate_up:
	docker exec -i $$(docker ps | grep loadbalancer-load_balancer | awk '{{ print $$1 }}') ./migrator
