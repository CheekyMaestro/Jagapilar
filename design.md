---
version: "alpha"
name: JAGAPILAR
description: Platform digital berbasis psikometri untuk deteksi dini "brainroot" (penurunan fokus & adiksi konten digital berkualitas rendah) pada anak. Ditujukan untuk orang tua, guru, dan siswa.
colors:
  primary: "#2FA7AC"
  secondary: "#1B4E73"
  tertiary: "#E84AB1"
  neutral: "#FFFFFF"
  neutral-dark: "#090A4C"
  illustration-teal: "#00CEC1"
  illustration-blue: "#00AAE7"
  illustration-pink: "#FF6AC5"
typography:
  h1:
    fontFamily: Poppins
    fontSize: 3rem
    fontWeight: 800
    lineHeight: 1.05
    letterSpacing: "-0.01em"
  h2:
    fontFamily: Poppins
    fontSize: 1.75rem
    fontWeight: 700
    lineHeight: 1.2
  body-md:
    fontFamily: Inter
    fontSize: 1rem
    fontWeight: 400
    lineHeight: 1.6
  label:
    fontFamily: Inter
    fontSize: 0.95rem
    fontWeight: 600
  button:
    fontFamily: Poppins
    fontSize: 0.95rem
    fontWeight: 700
    letterSpacing: "0.02em"
rounded:
  sm: 8px
  md: 16px
  lg: 28px
  full: 999px
spacing:
  xs: 8px
  sm: 16px
  md: 24px
  lg: 48px
  xl: 80px
components:
  nav-link:
    textColor: "{colors.neutral-dark}"
    typography: "{typography.label}"
  nav-link-active:
    textColor: "{colors.primary}"
  search-pill:
    backgroundColor: "{colors.secondary}"
    textColor: "{colors.neutral}"
    rounded: "{rounded.full}"
    padding: 12px
  button-primary:
    backgroundColor: "{colors.secondary}"
    textColor: "{colors.neutral}"
    typography: "{typography.button}"
    rounded: "{rounded.sm}"
    padding: 16px
  button-primary-hover:
    backgroundColor: "{colors.primary}"
  icon-badge:
    backgroundColor: "{colors.secondary}"
    textColor: "{colors.neutral}"
    rounded: "{rounded.full}"
  illustration-badge:
    backgroundColor: "{colors.tertiary}"
    rounded: "{rounded.full}"
  illustration-figure-teal:
    backgroundColor: "{colors.illustration-teal}"
    rounded: "{rounded.full}"
  illustration-figure-blue:
    backgroundColor: "{colors.illustration-blue}"
    rounded: "{rounded.full}"
  illustration-figure-pink:
    backgroundColor: "{colors.illustration-pink}"
    rounded: "{rounded.full}"
  bottom-bar:
    backgroundColor: "{colors.neutral}"
    rounded: "{rounded.lg}"
---

## Overview

JAGAPILAR memakai bahasa visual "klinis yang hangat" — kesan sebuah platform kesehatan/psikometri yang bisa dipercaya, tapi tetap ramah untuk orang tua dan anak. Gradasi teal-ke-navy yang dalam memberi kesan medis-profesional, sementara aksen pink dan bentuk-bentuk membulat (pill, lingkaran, siluet perisai) menjaga nuansa tetap hangat dan tidak terlalu kaku/korporat. Logo dan ilustrasi hero (tiga figur manusia abstrak membentuk perisai pelindung di sekeliling otak) adalah elemen tanda tangan brand ini: melambangkan perlindungan kolektif orang tua, guru, dan tenaga profesional terhadap tumbuh kembang kognitif anak.

## Colors

- **Primary — Teal (#2FA7AC):** Warna gradasi hero, aksen brand utama, garis bawah nav aktif.
- **Secondary — Navy Dalam (#1B4E73):** Warna "penggerak interaksi" — dipakai di tombol CTA, search pill, dan lingkaran ikon; juga titik akhir gradasi hero (bawah).
- **Tertiary — Pink (#E84AB1):** Aksen kecil saja (wordmark "PILAR", garis dekoratif tipis di bawah bottom bar). Jangan dipakai sebagai warna latar besar.
- **Neutral (#FFFFFF):** Latar header, teks di atas warna gelap, bottom bar.
- **Neutral Dark (#090A4C):** Teks navigasi & body di atas latar putih.
- **Illustration set (#00CEC1 teal terang, #00AAE7 biru langit, #FF6AC5 pink terang):** Khusus untuk tiga figur pada ilustrasi perisai hero — jangan dipakai di komponen UI lain agar ilustrasi tetap terasa sebagai elemen bertutur, bukan pola warna berulang.

Gradasi hero berjalan diagonal dari `{colors.primary}` (kiri-atas) menuju `{colors.secondary}` (kanan-bawah), dengan foto klinis low-opacity (~15–20%) di-blend di baliknya sebagai tekstur, bukan elemen utama.

## Typography

- **H1** — Poppins ExtraBold, huruf besar semua (uppercase), 3 baris pendek, line-height rapat (1.05) untuk kesan headline yang kuat dan mantap.
- **H2 / label section** — Poppins Bold, dipakai jarang, untuk judul section berikutnya (mis. "Services", "Our Doctor").
- **Body** — Inter Regular, line-height lega (1.6) agar paragraf penjelasan psikometri tetap mudah dibaca.
- **Nav label** — Inter SemiBold, ukuran sedang; item aktif ("Home") memakai warna `{colors.primary}` + garis bawah tipis.
- **Button** — Poppins Bold, uppercase, letter-spacing sedikit lebar, memberi kesan tegas & actionable.

## Layout

Struktur halaman: header sticky (tinggi ±90px) → hero dua kolom → bottom bar melayang yang tumpang-tindih (overlap) dengan tepi bawah hero.

```
┌─────────────────────────────────────────────┐
│ [Logo]   Home Services About Doctor  [Search]│  ← header, bg putih
├─────────────────────────────────────────────┤
│                                    · · · · · │
│  HEADLINE 3 BARIS         [ilustrasi         │  ← hero, gradasi
│  paragraf deskripsi        perisai + 3       │     teal→navy
│  [MULAI/BUAT AKUN →]       figur + otak]     │
│  ┌───────────────────────────────────────┐   │
│  │ 👪 Parent   📖 Teacher   🎓 Student    │   │  ← bottom bar putih,
│  └───────────────────────────────────────┘   │     rounded, overlap
└─────────────────────────────────────────────┘
```

Hero terbagi ~60/40: kolom kiri (teks+CTA) rata kiri dengan padding besar (`{spacing.xl}`), kolom kanan berisi pola titik dekoratif di pojok kanan-atas dan ilustrasi perisai di tengah. Bottom bar full-width, dibagi 3 kolom sama rata, dengan padding vertikal (`{spacing.md}`) dan sebuah garis aksen tipis warna `{colors.tertiary}` di tepi paling bawah.

## Elevation & Depth

Bottom bar "melayang" di atas hero dengan shadow lembut (soft, low-opacity, blur besar) — kesan kartu mengambang, bukan menempel rata. Foto latar di dalam hero diberi lapisan gradasi warna di atasnya (overlay) supaya teks putih tetap kontras tinggi; tidak ada elemen lain yang memakai shadow tebal — kedalaman visual sengaja dijaga minim di luar bottom bar ini.

## Shapes

- **Pill / full-round** (`{rounded.full}`): search bar, tombol ikon, lingkaran badge Parent/Teacher/Student.
- **Rounded rectangle sedang** (`{rounded.sm}`–`{rounded.md}`): tombol CTA utama, kartu-kartu section berikutnya.
- **Rounded besar** (`{rounded.lg}`): bottom bar mengambang.
- **Siluet perisai (signature shape):** bentuk perisai/crest yang dibentuk oleh tiga figur manusia — motif ini adalah elemen visual paling khas dari brand dan sebaiknya dipertahankan sebagai ilustrasi hero di setiap halaman utama, bukan diulang sebagai pola dekoratif biasa.

## Components

- **nav-link / nav-link-active:** Label navigasi Inter SemiBold; state aktif berwarna `{colors.primary}` dengan garis bawah.
- **search-pill:** Kapsul pencarian navy penuh dengan ikon kaca pembesar putih dan placeholder putih transparan.
- **button-primary / button-primary-hover:** Tombol CTA utama, navy dengan teks putih bold uppercase dan ikon panah dalam kotak putih terpisah di ujung kanan; hover berubah ke `{colors.primary}`.
- **icon-badge:** Lingkaran navy solid berisi ikon putih (dipakai di search bar dan tiga badge Parent/Teacher/Student).
- **illustration-badge:** Elemen aksen bulat kecil berwarna `{colors.tertiary}`, dipakai sangat terbatas (mis. titik logo).
- **bottom-bar:** Kartu putih rounded besar berisi 3 kolom icon-badge + label, mengambang di atas hero.

## Do's and Don'ts

- **Do** pertahankan gradasi teal→navy sebagai identitas hero di semua halaman (Services, About Us, Our Doctor, dst.) agar brand konsisten.
- **Do** jaga rounding besar di semua elemen interaktif (pill, lingkaran, kartu) — ini yang membuat brand kesehatan-anak terasa ramah, bukan klinis-dingin.
- **Do** batasi pink hanya sebagai aksen kecil, bukan warna dominan — supaya kesan kredibel/klinis tetap terjaga.
- **Don't** memakai sudut tajam (sharp corner) pada tombol atau kartu — bertentangan dengan tone ramah brand ini.
- **Don't** menaruh teks langsung di atas foto latar tanpa lapisan gradasi — kontras akan turun drastis.
- **Don't** menambah warna aksen baru di luar palet ini; disiplin warna adalah bagian dari kesan "dapat dipercaya" untuk produk kesehatan anak.
