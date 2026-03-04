# FLibGoLite — OPDS for Keenetic Routers (aarch64/mips)
===

**FLibGoLite** is an optimized fork of the original project, specifically adapted for stable performance on Keenetic routers and other resource-constrained embedded systems.

Original project and full documentation:
* [vinser/flibgolite](https://github.com/vinser/flibgolite)
* [Official Documentation](https://vinser.github.io/flibgolite-docs/en/docs/user-guide/)

### Key Differences from the Original Version:
* **Enhanced Stability**: Fixed critical errors (panics) that occurred when book files were missing from archives or deleted from the storage.
* **Resource Control**: Implemented a semaphore system to limit CPU and RAM usage during book scanning, cover generation, and file conversion.
* **Scalability**: Optimized SQLite performance and configuration to handle massive libraries (hundreds of thousands of books) on low-power hardware.
* **Language Logic**: Refactored localization handling. Book delivery is now independent of the interface language, ensuring consistent navigation across different OPDS clients.

---

### Installation on Keenetic

A correctly configured **Entware** environment is required on your router.

1. Copy the `flibgolite` binary to the `/opt/bin` directory.
2. Set execution permissions: `chmod +x /opt/bin/flibgolite`.
3. Copy the startup script [S99flibgolite](./S99flibgolite) to the `/opt/etc/init.d` directory.
4. Set execution permissions for the script: `chmod +x /opt/etc/init.d/S99flibgolite`.
5. Start the service to generate the default configuration: `/opt/etc/init.d/S99flibgolite start`.
6. Wait for **60 seconds**.
7. Verify that the process is running (using the `ps` command), then stop it: `/opt/etc/init.d/S99flibgolite stop`.
8. Open the configuration file at `/opt/bin/config/config.yml`, locate the `STOCK:` line, and enter the path to your book directory.
   * *Example:* `STOCK: /tmp/mnt/YOUR_DISK_ID/books/fb2`
9. Restart the service: `/opt/etc/init.d/S99flibgolite start`.

Scanning will begin after approximately 60 seconds. Once indexing is complete, your library is ready for use.

---

### Compatibility and Requirements
* **Tested on**: Keenetic Hopper (KN-1012).
* **Tested Clients**: `AlReaderX`, `FBReader`, `Cool Reader`.
* **Important Note**: A **SWAP partition** on the Entware drive is highly recommended for stable operation when indexing large libraries!

---

**This is an independent fork aimed at maximum performance in resource-limited environments.**

___*Suggestions and bug reports are welcome in the [Issues](https://github.com/alanneverland/flibgolite-keenetic-aarch64/issues) section of this repository.*___