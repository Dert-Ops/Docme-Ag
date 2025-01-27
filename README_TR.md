# Docme-Ag Yapay Zeka Dökümantasyon ve Versiyonlama Asistanı

## Genel Bakış
Bu proje, yazılım geliştirme süreçlerini takip eden ve yöneten bir yapay zeka dökümantasyon ve versiyonlama asistanı geliştirmeyi amaçlamaktadır:
- Kod geliştirme sürecini izler.
- Uygun aşamalarda otomatik olarak commit işlemi yapar.
- Gerekli yerlerde versiyonlama gerçekleştirir.
- Projenin dökümantasyonunu tutar.

## Özellikler
- **Otomatik Git Commit ve Versiyonlama**: Yapılan kod değişikliklerini takip eder ve gerektiğinde commit işlemi gerçekleştirir.
- **Akıllı Dökümantasyon Yönetimi**: Proje belgelerini otomatik olarak günceller ve düzenler.
- **Kolay Entegrasyon**: GitHub ile sorunsuz çalışır.
- **Özelleştirilebilir İş Akışı**: Kullanıcıların commit politikalarını tanımlamasına olanak tanır.
- **Çoklu Dil Desteği**: Farklı programlama dillerinde geliştirme süreçlerini takip eder.

## Commit Mesaj Standardı
Bu proje, **Conventional Commits** standardını takip etmektedir. AI asistanı otomatik olarak commit mesajlarını bu formata uygun şekilde oluşturmaktadır:

```
<type>(<scope>): <description>
```

- `type`: Değişiklik türü (örn. `feat`, `fix`, `chore`, `docs`, `style`, `refactor`, `test`).
- `scope`: Etkilenen modül veya bileşen.
- `description`: Yapılan değişiklik hakkında kısa bir açıklama.

**Örnekler:**
```
feat(api): yeni kimlik doğrulama sistemi eklendi
fix(ui): navbar'daki düzen sorunu giderildi
docs(readme): kurulum talimatları güncellendi
```

## Versiyonlama Standardı
Bu proje, **Semantic Versioning (SemVer)** modelini kullanmaktadır ve AI asistanı yapılan değişikliklere bağlı olarak versiyon güncellemelerini otomatik olarak belirler:

```
MAJOR.MINOR.PATCH
```

- **MAJOR**: Geriye dönük uyumluluğu bozan değişiklikler.
- **MINOR**: Yeni özellikler (geri uyumlu).
- **PATCH**: Hata düzeltmeleri ve küçük iyileştirmeler.

**Örnekler:**
```
1.0.0  # İlk stabil sürüm
1.1.0  # Yeni özellik eklendi
1.1.1  # Küçük bir hata düzeltildi
2.0.0  # Büyük değişiklik yapıldı, geri uyumluluk bozuldu
```

## Kurulum
1. Depoyu klonlayın:
   ```sh
   git clone https://github.com/yourusername/your-repo.git
   cd your-repo
   ```
2. Bağımlılıkları yükleyin:
   ```sh
   pip install -r requirements.txt
   ```
3. Asistanı çalıştırın:
   ```sh
   python main.py
   ```

## Kullanım
- Asistan, kod değişikliklerini sürekli izler.
- Tanımlanan commit politikasına göre otomatik commit yapar ve Conventional Commits formatını kullanır.
- Proje dökümantasyonu ve versiyon kayıtlarını tutar.
- Semantic Versioning modeline göre otomatik versiyon güncellemesi yapar.

## Katkıda Bulunma
Katkılar memnuniyetle karşılanır! Lütfen değişiklikleri tartışmak için bir issue açın veya pull request gönderin.

## Lisans
Bu proje MIT Lisansı ile lisanslanmıştır. Ayrıntılar için [LICENSE](LICENSE) dosyasına bakabilirsiniz.
