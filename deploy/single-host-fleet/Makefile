st:
	fleetctl list-units

run_redis:
	fleetctl start redis-discovery.service
	fleetctl start redis.service
destroy_redis:
	fleetctl destroy redis-discovery.service
	fleetctl destroy redis.service

run_rethinkdb:
	fleetctl start rethinkdb-discovery.service
	fleetctl start rethinkdb.service

destroy_rethinkdb:
	fleetctl destroy rethinkdb-discovery.service
	fleetctl destroy rethinkdb.service

run_api:
	fleetctl start api.service

destroy_api:
	fleetctl destroy api.service

run_worker:
	fleetctl start worker.service

destroy_worker:
	fleetctl destroy worker.service
