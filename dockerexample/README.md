# Deploy golang menggunakan docker

## 1. Build image
```
docker build -t go-docker-example .
```

Command di atas akan melakukan proses build Image pada file yang ada di dalam . yang merupakan isi folder project. Project akan di-build ke sebuah Image dengan nama adalah `go-docker-example`. Flag `-t` digunakan untuk menentukan nama Image.

Gunakan command docker images untuk menampilkan list semua image yang ada di lokal.
```
$ docker images
REPOSITORY             TAG       IMAGE ID       CREATED             SIZE
go-docker-example   latest    5c9bff44d224   About an hour ago   316MB
```

## 2. Access Local Database

Untuk host database, jika ingin menggunakan host yang ada di local pastikan config mysql supoort alamat ip selain localhost. Anda bisa mengedit file `*.cnf`, yang di ubuntu 23.10 ada di `/etc/mysql/mysql.conf.d/mysqld.cnf`. Lalu ubah `bind-address` menjadi
```
bind-address = 0.0.0.0
```

Lalu anda bisa gunakan host `host.docker.internal` atau bisa menggunakan ip dengan command 
```
ip addr show docker0
```
yang hasilnya bisa seperti:
```
3: docker0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default 
    link/ether 02:42:6a:70:8c:47 brd ff:ff:ff:ff:ff:ff
    inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0
       valid_lft forever preferred_lft forever
    inet6 fe80::42:6aff:fe70:8c47/64 scope link 
       valid_lft forever preferred_lft forever
```
maka yang perlu anda ambil sebagai host iyalah di baris `inet`, yaitu ip `172.17.0.1`. 

Catatan: dari kasus saya, host yang diset adalah `172.17.0.1`, tapi yang dipanggil dari golang `172.17.0.2`. Jadi pastikan user mysql bisa akses dari ip tersebut.


## 3. Create Container

Image yang sudah siap, bisa dibuat `container` nya menggunakan basis image `go-docker-example`.

```
docker container create --name go-container-example -e PORT=8080 -e MYSQL_HOST="172.17.0.1" -e MYSQL_USERNAME="firman" -e MYSQL_PASSWORD="mautauaja" -p 8080:8080 go-docker-example
```

Command di atas akan menjalankan sebuah proses yang isinya kurang lebih berikut:

1. Buat container baru dengan nama `go-container-example`.
2. Flag `--name` digunakan untuk menentukan nama container.
3. Sewaktu pembuatan container, env var PORT di-set dengan nilai adalah `8080`.
4. env var INSTANCE_ID juga di set di-set, nilai adalah teks `my first instance`.
5. Flag `-e` digunakan untuk menge-set env var. Flag ini bisa dituliskan banyak kali sesuai kebutuhan.
6. Kemudian port `8080` yang ada di luar network docker (yaitu di host/laptop/komputer kita) di map ke port 8080 yang ada di dalam container.
7. Flag `-p` digunakan untuk mapping port antara host dan container. Bagian ini biasa disebut dengan `expose port`.
8. Proses pembuatan container dilakukan dengan Image `go-docker-example` digunakan sebagai basis image.

## 3. Start Container

Sekarang container sudah dibuat, lanjut untuk start container tersebut, caranya menggunakan command 
```
docker container start go-container-example
```

Jika sudah, coba cek di browser aplikasi web hello world, harusnya sudah bisa diakses.

Gunakan command ini untuk memunculkan list container yang sedand running atau aktif. Untuk menampilkan semua container (aktif maupun non-aktif), cukup dengan menambahkan flag -a atau --all.
```
docker container ls
```

## 4. Stop Container
Untuk stop container bisa dengan command:
```
docker container stop go-container-example
```

## 5. Hapus Container
Untuk hapus container bisa dengan command:
```
docker container rm go-container-example
```

## 6. Hapus Image
Untuk hapus image bisa dengan command:
```
docker image rm go-docker-example
```