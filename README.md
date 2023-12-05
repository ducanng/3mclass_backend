# no-name
# Install: Follow make file
```
➜ make all
```
# Dev:
- Load postgres docker container: `docker pull postgres:15` or `docker pull --platform linux/x86_64 postgres:15` for mac m1
- Load redis docker container: `docker pull redis:6.2`
- Set up h local instance:
```
➜  docker run --name posgres-instance \
   -e POSTGRES_PASSWORD=password \
   -e POSTGRES_DB=threemanclass \
   --restart=unless-stopped \
   -p 5432:5432 \
   -d postgres:15
    