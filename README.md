## connect to server:
1. pertama jika ingin tersambung ke server :
   ```sh
   ssh -p 60606 jose@croot.ypbpi.or.id
   pw : joseganteng123
   ```
   ```sh
   jika lokal maka dapat melalui ip:
   10.14.200.210
   ```
note.
jangan sampe salah lebih dari 3 kali karena akun akan di block jika slaah lebih dari 3 kali

2. kemudian jika ingin cek debv:
   ```sh
	pertama:
	cd /etc/nginx/sites-enabled
	kedua:
	ls
	Ketiga:
	cek system yang ingin dilakukan pengecekan
   ```
3. jika ingin masuk ke nano system:
   ```sh
	sudo nano (nama server)
   ```
4. jika sudah melakukan perubahan pada server:
   ```sh
	sudo nginx -t
   ```
5. kemudian restart system server:
   ```sh
	sudo system ctl restart nginx
   ```
6. jika ingin cek logs dari server:
   ```sh
	masuk docker dahulu:
	ssh -p 60606 docker@10.14.200.20
	password:
	1
	kemudian code cek logs:
	docker logs -f (nama server)
   ```