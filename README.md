# Request examples
## Upload
```bash
curl --location --request POST 'http://127.0.0.1:8080/' \
   --header 'Content-Type: multipart/form-data' \
   --form 'file=@/home/forty/Downloads/audio.wav'
```

## Download
```bash
curl --location --request GET 'http://127.0.0.1:8080/?hash=1257cfb4e5ac35d5c32f5103691321001775609b'
```

## Delete
```bash
curl --location --request DELETE 'http://127.0.0.1:8080/?hash=1257cfb4e5ac35d5c32f5103691321001775609b'
```