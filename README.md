<p align="center"><img src="https://image.ibb.co/buoQzz/Logo_T1.jpg" alt="Project Logo" width="300"></p>

[![Go Report Card](https://goreportcard.com/badge/github.com/mucanyu/eksisozluk-go)](https://goreportcard.com/report/github.com/mucanyu/eksisozluk-go)
[![GoDoc](https://godoc.org/github.com/mucanyu/eksisozluk-go?status.png)](https://godoc.org/github.com/mucanyu/eksisozluk-go)


# eksisozluk-go 🌢
Komut satırından `ek$isözlük` gündemini takip etmenizi ve başlık içerisindeki entryleri okumanızı sağlayan uygulama.

**Go** ile geliştirildi.

### Nasıl?

#### -> Kullanırım
- İşletim sisteminize göre çalıştırılabilir dosyayı [indirin](https://github.com/mucanyu/eksisozluk-go/releases/tag/v1.0.0).
- İsim değişikliği yaptıktan sonra çalıştırma yetkisi verin.

    ```
    $ chmod +x eksisozluk-go
    ```

#### -> Derlerim

Gereklilikler:
`Go 1.7 ve üzeri`
```
$ go get github.com/mucanyu/eksisozluk-go
$ cd $GOPATH/src/github.com/mucanyu/eksisozluk-go
$ go install
```

### Komutlar
```
eksisozluk-go gundem [-kategori=#] [-limit=#] [-sayfa=#]
eksisozluk-go baslik <BAŞLIK> [-sukela] [-limit=#] [-sayfa=#]
```

### Katkıda Bulunun
Bilimum pull request kabul edilir.
