# 🕵️ ThorScraper

Tor ağı üzerinden güvenli web scraping ve screenshot aracı. Dark web sitelerini tarayıp analiz eder.

## ✨ Özellikler

- 🔒 **Tor Entegrasyonu**: Otomatik Tor bağlantısı (Port 9050/9150)
- 🌐 **Toplu Tarama**: Yüzlerce siteyi tek seferde tara
- 📸 **Screenshot**: Her sitenin otomatik ekran görüntüsü
- 📝 **HTML Kaydetme**: Sayfa içeriğini yerel kayıt
- 📊 **Detaylı Loglama**: Tüm işlemlerin kayıtları
- 🎯 **Seçici Tarama**: İstediğin siteleri seç ve tara
- ⚡ **Hızlı & Güvenli**: Proxy üzerinden anonim tarama

## 📋 Gereksinimler

### 1. Tor Servisi (Zorunlu)

ThorScraper'ın çalışması için aktif bir Tor bağlantısı gereklidir. İki seçenek:

#### Seçenek A: Tor Browser (Kolay - Önerilen)

1. [Tor Browser'ı indir](https://www.torproject.org/download/)
2. Tor Browser'ı çalıştır ve bağlantıyı bekle
3. **Tor Browser açık kaldığı sürece** ThorScraper çalışır (Port: 9150)

#### Seçenek B: Tor Expert Bundle (Gelişmiş)

**Windows:**
1. [Tor Expert Bundle](https://www.torproject.org/download/tor/) indir
2. `tor.exe` dosyasını çalıştır
3. Komut satırında çalışır durumda kal (Port: 9050)

**Linux:**
```bash
# Tor servisini kur
sudo apt install tor

# Servisi başlat
sudo systemctl start tor
sudo systemctl enable tor

# Durumu kontrol et
sudo systemctl status tor
```

**macOS:**
```bash
# Homebrew ile kur
brew install tor

# Çalıştır
tor
```

## 🚀 Kurulum

### Yöntem 1: Hazır Binary (Windows - ÖNERİLEN)

1. Repository'yi klonla veya ZIP olarak indir
2. `ThorScraper.exe` dosyası hazır!
3. Tor Browser'ı çalıştır
4. `ThorScraper.exe` dosyasını çalıştır

```powershell
git clone https://github.com/KULLANICI_ADIN/ThorScraper.git
cd ThorScraper
.\ThorScraper.exe
```

### Yöntem 2: Kaynak Koddan Derle (İsteğe Bağlı)

**Gereksinimler:**
- Go 1.21 veya üzeri

```bash
go mod download
go build -o ThorScraper.exe
```

**Linux için:**
```bash
go build -o ThorScraper
chmod +x ThorScraper
./ThorScraper
```

**macOS için:**
```bash
GOOS=darwin GOARCH=amd64 go build -o ThorScraper
./ThorScraper
```

## 📖 Kullanım

### Adım 1: Tor Servisini Başlat

**Tor Browser ile (Kolay):**
1. Tor Browser'ı aç
2. "Connect" butonuna tıkla
3. Bağlantı kurulana kadar bekle (yeşil soğan simgesi)
4. Tor Browser'ı **kapatma** - arka planda çalışsın

**Tor Servisi ile:**
```bash
# Windows
tor.exe

# Linux/Mac
tor
```

### Adım 2: Hedef URL'leri Ekle

`targets.yaml` dosyasını düzenle:

```yaml
# İsim | URL formatı
DarkNetArmy | http://darknet77vonbqeatf...onion/
GhostHub | http://aniozgjggq2pzxzn...onion/

# Sadece URL formatı (otomatik isim)
http://example123456789.onion/
```