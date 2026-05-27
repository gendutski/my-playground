# TICKET WAR SIMULATION

Setelah menonton [Cara Web Flashsale Biar Gak Down](https://www.youtube.com/watch?v=srYIazbZ_-8), yaitu membuat `Slot Pool` dengan redis `INCR`, `DECR`, dan set time span untuk yang berhasil masuk. Ada yang mengganjal dengan cara tersebut

Jika user sudah mendapatkan slot, tapi abandon / close browser, maka command `DECR` tidak akan pernah dijalankan. Akibatnya counter menjadi “bocor”.

Contoh:

```
limit 1000
300 orang abandon
counter tetap terisi 300
```

Padahal slot sebenarnya kosong.

## Solusi

Solusi sederhana untuk masalah itu, yaitu dengan menggunakan atomic lua script yang memproses beberapa command redis sekaligus

## Cara menggunakan

Jalankan server

```
php -S localhost:8000
```

Atau copy folder ini ke wwwroot seperti `/var/www/html `
