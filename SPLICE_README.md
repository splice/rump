In order to work with the the older version of Redis we had to make some modifications.

## Build and upload

1. Build the binary: `make build`
2. Build the container: `make build_container`
3. Upload the container: `make upload`

If you need to do a clean build of the container, use `make clean_build_container`

If making changes, make sure to bump the version in the `VERSION` file

## Pull down and run the container

1. Pull container: `git pull quay.io/splice/rump:0.0.2`
2. Spinup the container: `docker run -it quay.io/splice/rump:0.0.2`
3. Example rump run: `time rump -from redis://staging-redis.jopjwe.ng.0001.usw1.cache.amazonaws.com:6379/0 -to redis://staging-redis-2.jopjwe.ng.0001.usw1.cache.amazonaws.com:6379/0`

## Migrations

Here's a link to more complete documentation for running a migration using this tool.

- [Elasticache Redis Migration](https://www.notion.so/splice/Elasticache-Redis-Migration-dbb7d37a75804d0da663023424ac6dbb)
