### REDIS相关
sudo docker run -d --name tg_gif-redis -v $(pwd):/usr/local/etc/redis -p 6379:6379 redis redis-server /usr/local/etc/redis/redis.conf
sudo docker run -it --link tg_gif-redis:redis --rm redis redis-cli -h redis -p 6379
### MYSQL相关
sudo docker run -d \
	-p 3306:3306 \
	--name gif-sql \
	--restart always \
	-v $(pwd)/mysql.d:/etc/mysql/conf.d \
	-e MYSQL_ROOT_PASSWORD=AwMTEyMCAzNjMzIGYyMTAgYmI \
	-e MYSQL_DATABASE=tg_gif \
	-e MYSQL_USER=tg_gif \
	-e MYSQL_PASSWORD=lYmYwIGUzOGEgOTBkZCBjNGRlIDNk \
	mysql:5.7 \
	--character-set-server=utf8mb4 \
	--collation-server=utf8mb4_unicode_ci
sudo docker run -it --rm --link gif-sql:mysql -v $(pwd)/mysql.d:/etc/mysql/conf.d --rm mysql:5.7 sh -c 'exec mysql -h 172.17.0.1 -P 3306 -u root -p tg_gif'
