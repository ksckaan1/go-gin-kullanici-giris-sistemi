package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	anahtar = "kullanıcı"
)

func main() {
	r := kurulum()
	//kurulum fonksiyonunda adres yönlendirmelerimizi yapacağız

	r.Use(gin.Logger())
	//terminalden logları görmek için logger kullandık.

	if err := kurulum().Run(":8080"); err != nil {
		log.Fatal("Sunucu başlatma hatası:", err)
	}
	//klasik go tarzı hata yakalama yaptık.
}

func kurulum() *gin.Engine {
	r := gin.New()
	//Yeni bir router oluşturduk.

	r.Use(sessions.Sessions("oturum", sessions.NewCookieStore([]byte("bilgi"))))
	//Tarayıcımızda barındıracağımız çerezi oluşturduk.

	r.LoadHTMLGlob("sablonlar/*")
	//şablon klasörümüzü ekledik.

	r.Static("/static", "./static")
	//statik dosyalarımız için sunucuda dizin oluşturduk.
	//bu durumda style.css için işe yarayacak

	r.GET("/giris", giris)
	//giriş sayfamıza yönlendirme yaptık.

	r.POST("/giris-kontrol", girisKontrol)
	//giriş bilgilerini kontrol edeceğimiz yönlendirmemiz.
	//giriş bilgileri formdan geleceği için post metodu ile oluşturduk.

	r.GET("/cikis", cikis)
	//çıkış işlemi için yönlendirmemiz

	hesap := r.Group("/hesap")
	//oturum açmış kullanıcılara özel olan kapsayıcı adresimiz

	hesap.Use(oturumKontrol)
	//hesap dizini altındaki tüm adreslere oturum kontrolü yapalım.
	{
		hesap.GET("/ben", ben)
		//kullanıcı adımızı gösterecek yönlendirme

		hesap.GET("/durum", durum)
		//oturum açık mı değil mi diye gösteren yönlendirme
	}
	return r
}

//giriş formumuzun bulunduğu sayfa
func giris(c *gin.Context) {
	c.HTML(200, "giris.html", gin.H{})
	//giris.html dosyasını şablon olarak ekledik
}

//formdan gelen kullanıcı bilgilerini kontrol eden yönlendirmemiz
func girisKontrol(c *gin.Context) {
	oturum := sessions.Default(c)
	//oturum bilgilerimizi kaydedeceğimiz çerez deposunu aldık.

	kullaniciAdi := c.PostForm("kullaniciadi")
	//formdan gelen kullaniciadi bilgisi

	sifre := c.PostForm("sifre")
	//formdan gelen sifre bilgisi

	if strings.Trim(kullaniciAdi, " ") == "" || strings.Trim(sifre, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Hata": "Kullanıcı adı ve şifre boş olamaz"})
		return //hata varsa fonksiyon sonlanacak
	}
	//kullaniciadi ve sifre bilgilerinin uygunluğunu kontrol ettik.

	if kullaniciAdi != "kaan" || sifre != "kaan123" {
		c.JSON(http.StatusUnauthorized, gin.H{"hata": "Yanlış kullanıcı adı veya şifre"})
		return //hata varsa fonksiyon sonlanacak
	}
	//kullanıcı adı veya şifre yanlışsa hata vermesini istedik.
	//(Opsiyonel) Bu kısımda veri tabanı bağlantısı ile kontrol edilebilir.

	oturum.Set(anahtar, kullaniciAdi)
	// çerezimizin anahtar bölümüne kullaniciAdi atanacak
	if err := oturum.Save(); err != nil {
		//çerezi kaydediyoruz.
		c.JSON(http.StatusInternalServerError, gin.H{"Hata": "Oturumu kaydederken sunucu hatası oluştu"})
		return //hata varsa fonksiyon sonlanacak
	}
	c.JSON(http.StatusOK, gin.H{"mesaj": "Başarıyla giriş yapıldı"})
	//Bir sıkıntı ile karşılaşılmadığında giriş yapmış olacağız.
}

//kullanıcının çıkış yapacağı yönlendirme
func cikis(c *gin.Context) {
	oturum := sessions.Default(c)
	//çerez oturumumuzu aldık

	kullanici := oturum.Get(anahtar)
	//kullanıcı adımızın kayıtlı olduğu anahtarı aldık.

	if kullanici == nil {
		//anahtar içerisinde bir kayıt yoksa oturum zaten açılmamız demektir
		c.JSON(http.StatusBadRequest, gin.H{"Hata": "Şuanda oturum açılmamış"})
		return //hata varsa fonksiyon sonlanacak
	}

	oturum.Delete(anahtar)
	//anahtarımız içerisindeki kullanıcı adını kaldırdık.

	if err := oturum.Save(); err != nil {
		//ve çerez oturumunu kaydettik.
		c.JSON(http.StatusInternalServerError, gin.H{"Hata": "Çıkış yaparken sunucu hatası oluştu"})
		return //hata varsa fonksiyon sonlanacak
	}
	c.JSON(http.StatusOK, gin.H{"Mesaj": "Oturumdan başarıyla çıkıldı"})
	//bir sıkıntı çıkmadığında çıkış yağtığımıza dair mesaj görüntüledik.
}

//oturumun açık olup olmadığını kontrol eden fonksiyonumuz
func oturumKontrol(c *gin.Context) {
	oturum := sessions.Default(c)
	//oturum bilgilerimizi aldık

	kullanici := oturum.Get(anahtar)
	//çerezdeki kullanıcı adımızı aldık.

	if kullanici == nil {
		//çerezde kullanıcı adı bulunmuyorsa oturum açılmamış demektir
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Hata": "Erişim yok"})
		return //hata varsa fonksiyon sonlanacak
	}
	c.Next()
	//Eğer oturum açılmış ise Next fonksiyonu ile alt adreslere yönlendirme mümkün oluyor
}

//kullanıcı adımızın gösterileceği yönlendirme
func ben(c *gin.Context) {
	oturum := sessions.Default(c)
	//oturum bilgilerimizi aldık

	kullanici := oturum.Get(anahtar)
	//anahtar çerezinden kullanıcı adımızı aldık.

	c.JSON(http.StatusOK, gin.H{"Kullanıcı": kullanici})
	//kullanıcı adımızı gösterdik
}

//oturumumuzun durumunu gösteren yönlendirme
func durum(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Durum": "Şuanda oturum açılmış durumda"})
	//oturum açık ise mesaj gösterdik.
}
