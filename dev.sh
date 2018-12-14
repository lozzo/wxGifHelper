sudo docker run -d --name tg_gif-redis -v $(pwd):/usr/local/etc/redis -p 6379:6379 redis redis-server /usr/local/etc/redis/redis.conf
docker run -it --link tg_gif-redis:redis --rm redis redis-cli -h redis -p 6379
