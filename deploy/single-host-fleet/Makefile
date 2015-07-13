st:
	fleetctl --endpoint FLEET_HOST list-units

run_redis:
	fleetctl --endpoint FLEET_HOST start redis-discovery.service
	fleetctl --endpoint FLEET_HOST start redis.service
destroy_redis:
	fleetctl --endpoint FLEET_HOST destroy redis-discovery.service
	fleetctl --endpoint FLEET_HOST destroy redis.service

run_rethinkdb:
	fleetctl --endpoint FLEET_HOST start rethinkdb-discovery.service
	fleetctl --endpoint FLEET_HOST start rethinkdb.service

destroy_rethinkdb:
	fleetctl --endpoint FLEET_HOST destroy rethinkdb-discovery.service
	fleetctl --endpoint FLEET_HOST destroy rethinkdb.service

run_api:
	fleetctl --endpoint FLEET_HOST start api.service

destroy_api:
	fleetctl --endpoint FLEET_HOST destroy api.service

run_worker:
	fleetctl --endpoint FLEET_HOST start worker.service

destroy_worker:
	fleetctl --endpoint FLEET_HOST destroy worker.service
