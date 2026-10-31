# Essay menilai salah jawaban yang tidak menyertakan catatan dalam tanda kurung

- **Tanggal**: 2026-10-10
- **Komponen**: Frontend — Tebak Kanji & Tebak Kotoba (mode essay)

## Ringkasan

Pada mode essay, jawaban pengguna ditandai **salah** padahal seharusnya **benar**
ketika arti menyertakan catatan di dalam tanda kurung.

Contoh: kanji `田中` (knj-018) memiliki arti `"Tanaka (nama orang)"`. Ketika
pengguna menjawab `tanaka`, jawaban dianggap salah.

## Penyebab

Perbandingan jawaban memakai `normalizeAnswer()` yang sudah menangani perbedaan
huruf besar/kecil (`toLowerCase`), jadi bukan masalah kapitalisasi.

Akar masalah sebenarnya: arti pada data menyertakan catatan di dalam tanda
kurung, mis. `"Tanaka (nama orang)"`. Pembanding membandingkan seluruh string,
sehingga `tanaka` tidak sama dengan `tanaka (nama orang)`.

## Solusi

`frontend/src/lib/quiz.js` — `normalizeAnswer()`: buang lebih dulu bagian di
dalam tanda kurung, baik `()` maupun full-width `（）`, sebelum proses
normalisasi lain (trim, lowercase, hapus spasi & tanda baca).

```js
export function normalizeAnswer(value) {
  return String(value ?? "")
    .replace(/[（(][^）)]*[）)]/g, "")
    .trim()
    .toLowerCase()
    .replace(/\s+/g, "")
    .replace(/[.。,，、!！?？'"`’“”\-ー…]/g, "")
}
```

Berlaku untuk Quiz Kanji (arti & cara baca, mis. `"だれ (どなた)"` → `"だれ"`)
dan Quiz Kotoba (arti terjemahan).

## Verifikasi

- `normalizeAnswer("Tanaka (nama orang)") === normalizeAnswer("tanaka")` → true
- `normalizeAnswer("AM (pagi / sebelum tengah hari)") === normalizeAnswer("am")` → true
- `normalizeAnswer("だれ (どなた)") === normalizeAnswer("だれ")` → true
- `bun run build` berhasil.

## Berkas

- `frontend/src/lib/quiz.js`
