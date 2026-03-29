# FLibGoLite 📚 — Your personal OPDS library, made simple

FLibGoLite is a lightweight, fast, and easy-to-set-up OPDS service for your home library. No fuss, no bloat — just your books, available anywhere.

Detailed guides (in multiple languages) are available [here](https://vinser.github.io/flibgolite-docs/).

![GitHub Release](https://img.shields.io/github/v/release/vinser/flibgolite?label=release&sort=semver)
![Docker](https://img.shields.io/docker/pulls/vinser/flibgolite?logo=docker)

## Why you'll like it ✨
- **Read what you love:** Supports EPUB and FB2 (files or zip) — no format stress.
- **FB2 → EPUB on the fly:** Got FB2 books but your reader prefers EPUB? I've got you covered — FLibGoLite converts them automatically when you download.
- **Runs anywhere:** Linux, Windows, MacOS, FreeBSD — pick your platform.
- **No dependencies:** Just one self-contained binary. Download and run.
- **Docker-ready:** Prefer containers? There's a pre-built image waiting for you.
- **Your setup, your way:** Run as a system service, in Docker, or just from the terminal.
- **Fast & light:** Quick indexing, SQLite storage, and a snappy OPDS service that doesn't hog resources.
- **Speaks your language:** Built-in localization, easy to switch.
- **Docs that actually help:** Clear, practical guides — no guessing required.

## Get started in minutes 🚀

### What you'll need
- A PC, NAS, or server (Windows, MacOS, or Linux).
- Any OPDS-compatible reader app or device that handles EPUB or FB2.
  - *Tested and happy with:* `PocketBook Reader`, `FBReader`, `Librera Reader`, `Cool Reader` (mobile), `Foliate`, `Thorium Reader` (desktop).
  - *Got something else?* If it speaks OPDS and reads EPUB/FB2, it'll probably work just fine.

### Choose how you want to run it

#### 💻 Option 1: Run the binary directly
Prefer to go old-school? Follow this [guide](https://vinser.github.io/flibgolite-docs/en/docs/user-guide/) for your OS.
1. Drop your books (EPUB, or FB2 as zip/loose files) into `books/stock`.
2. FLibGoLite will index them automatically — no manual cataloging needed.
3. Point your reader to: `http://server:8085/opds`
   - Replace `server` with your PC's hostname or IP (e.g., `192.168.0.10`).
4. Browse, search by author/genre/title, and start reading.
   - Need EPUB? FB2 books convert automatically on download if your reader doesn't support them.

#### 🐳 Option 2: Docker (quick & clean)
I've built a Docker image so you can skip the setup hassle:
```bash
docker run -d -p 8085:8085 -v ./books:/app/books/stock vinser/flibgolite
```
*   `-p 8085:8085` — exposes the service port.
*   `-v ./books:/app/books/stock` — points the container to your books folder.

That's it. Your library is now live at `http://localhost:8085/opds`.

###  And you're done! 🎉
Now your whole household can enjoy your library — from phones, e-readers, or desktops. Share it with family, set it up for friends, or just keep it for yourself.

Happy reading! 📖

---
Found a bug? Have an idea? Let me know [here](https://github.com/vinser/flibgolite/issues).
