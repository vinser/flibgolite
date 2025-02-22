FLibGoLite - Just enough for free OPDS 
===

__FLibGoLite__ is an easy-to-install, fast and resource-friendly OPDS service.  

__Detailed multilingual guides are available [here](https://vinser.github.io/flibgolite-docs/)__
### CURRENT STABLE RELEASE v2.2.4

__FLibGoLite__ main features:
- Works with books in FB2 (separate files and zip archives) and EPUB formats
- Ability to convert FB2 format to EPUB when loading a book into the reader
- Multiplatform: Linux, Windows, MacOS, FreeBSD
- Self-sufficiency - does not require installation of additional libraries or applications
- Can be launched as a system service or in a Docker container, as well as from the command line
- Fast inexing and keep persistent data in SQLite database
- High speed of processing new arrivals and saving the catalog in the SQLite database
- Fast and responsive OPDS service with a simple localization
- Well documented

#### Briefly how to use FLibGoLite.

You need:

1. PC, NAS or server with Windows, MacOS or Linux operating system.
2. A reader (device or application) that can work with OPDS catalogs and supports FB2 or EPUB book formats.
FLibGoLite has been tested and works with mobile applications for reading books `PocketBook Reader`, `FBReader`, `Librera Reader`, `Cool Reader`, as well as desktop applications `Foliate` and `Thorium Reader`. You can use any other applications or devices that can read the listed book formats and work with OPDS catalogs.

Follow these [guide](https://vinser.github.io/flibgolite-docs/en/docs/user-guide/) to install FLibGoLite on your PC.

Put your books in FB2 format (zip-archives or separate files) or EPUB in the `books/stock` folder. The service processes them and enters the books details into the catalog.

Next, configure the reader(s) to work with the OPDS directory `http://server:8085/opds`,  
where `server` is the name of your PC or the IP address of the PC type `192.168.0.10`  
After that, you can select and download any of the books stored on the PC in the reader.  
Books can be selected/searched by author and genre, as well as contextual search by author and/or book title.  
For book readers that do not support the FB2 format, books can be converted to EPUB format when loaded.

Thus, you will create a library that will be used by your loved ones with smartphones, reader devices or PCs.

Good luck!

---
___*Suggestions, bug reports and comments are welcome [here](https://github.com/vinser/flibgolite/issues)*___

   

