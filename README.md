# no-name
# Install: Follow make file
```
➜ make all
```
# Dev:
- Load mysql docker container: `docker pull mysql` or `docker pull --platform linux/x86_64 mysql` for mac m1
- Load redis docker container: `docker pull redis:6.2`
- Set up a local instance:
```
➜  docker run --name mysql-instance \
   -e MYSQL_ROOT_PASSWORD=password \
   -e MYSQL_DATABASE=noname \
   -e MYSQL_USER=aemi \
   --restart=unless-stopped \
   -p 3336:3306 \
   -d mysql:latest
