# FLibGoLite
[ *русский вариант здесь* ](README_RU.md)

### BETA RELEASE v0.1.x * 
_*This software release has not been tested thoroughly yet but based on __[flibgo](https://github.com/vinser/flibgo.git)__ it does the job_

---

__FLibGoLite__ is easy to use home library OPDS server 

>The Open Publication Distribution System (OPDS) catalog format is a syndication format for electronic publications based on Atom and HTTP. OPDS catalogs enable the aggregation, distribution, discovery, and acquisition of electronic publications. [(Wikipedia)](https://en.wikipedia.org/wiki/Open_Publication_Distribution_System)

__FLibGoLite__ is multiplatform lightweight OPDS server with SQLite database book search index
This __FLibGoLite__ release only supports FB2 publications, both individual files and zip archives.

OPDS-catalog is checked and works with mobile readers FBReader and PocketBook Reader


##  Install and run program
---


   **FLibGoLite** is written in GO as a single executable and doesn't require any prerequsites  
   All you have to do is to download OS and hardware specific build and run it.

|OS        |Hardware              |Program executable          |Tested  |  
|----------|----------------------|----------------------------|:------:|  
|Windows   | Intel, AMD 32-bit    | flibgolite-linux-386.exe   |Yes     |  
|Windows   | Intel, AMD 64-bit    | flibgolite-linux-amd64.exe |Yes     |  
|OS X (MAC)| 64-bit               | flibgolite-darwin-64       |No      |  
|Linux     | Intel, AMD 32-bit    | flibgolite-linux-386       |No      |  
|Linux     | Intel, AMD 64-bit    | flibgolite-linux-amd64     |Yes     |  
|Linux     | ARM 32-bit (armhf)   | flibgolite-linux-arm-6     |Yes     |  
|Linux     | ARM 64-bit (armv8)   | flibgolite-linux-arm64     |Yes     |  
  
For convenience you may rename downloaded program executable to __flibgolite__  

__FLibGoLite__ program may run from command line or may install itself as a service running in background

Usage:

	./flibgolite [OPTION]

	With no OPTION program will run in console mode (Ctrl+C to exit)  
	Caution: Only one OPTION can be used at a time

	OPTION should be one of:

	-service [action]     control FLibGoLite service
		where action is one of: install, start, stop, restart, uninstall, status  
	-reindex              empty book stock index and then scan book stock directory to add books to index
	-config               create default config file in ./config folder for customization and exit
	-help                 display this help and exit
	-version              output version information and exit

Service control option (install, start ...) requires administrator rights.  
On windows you should accept running as administrator, on linux - use `sudo`

Examples:

	./flibgolite                      		Run FLibGoLite in console mode
	sudo ./flibgolite -service install     	Install FLibGoLite as a system service
	sudo ./flibgolite -service start	


At the first run program will create the set of subfolders in current folder

 	flibgolite
	├─── books  
	|    ├─── new   - new book files and/or zip archives with book files should be put here for scan
	|    ├─── stock - library book files and archives are stored here
	|    └─── trash - files that have been processing bugs will come here
	├─── config - contains main configiration file config.yml and genre tree file
	|    └─── locales - localization files
	├─── dbdata - database with book index resides there
	└─── logs - scan and opds rotating logs are there

 ## Advanced usage

   For advanced sutup see config/config.yml selfexplanatory file.
```yml
library:
  # Selfexplained folders
  BOOK_STOCK: "books/stock"
  NEW_ACQUISITIONS: "books/new"
  TRASH: "books/trash"

language:
  # Locales folder. You can add your own locale file there like en.yml or ru.yml
  LOCALES: "config/locales"
  # Default english locale can be changed to "ru" for Russian opds feeds (bookreaders opds menu tree)
  DEFAULT: "en"  

genres:
  TREE_FILE: "config/genres.xml"
  # Alternative genres tree can be used (Russian only, sorry) 
  # TREE_FILE: "config/alt_genres.xml"
  
database:
  DSN: "dbdata/flibgolite.db"
  # Delay before start each new acquisitions processing
  POLL_DELAY: 30 
  # Maximum simultaneous new aquisitios processing threads
  MAX_SCAN_THREADS: 3
  # Accept only these languages puplications. Add others if needed please.
  ACCEPTED_LANGS: "en,ru,ua"

logs:
  # Logs are here
  OPDS: "logs/opds.log"
  SCAN: "logs/scan.log"
  DEBUG: false

opds:
  # OPDS-server port so opds can be found at http://<server name or IP-address or localhost>:8085/opds
  PORT: 8085
  # OPDS feeds entries page size
  PAGE_SIZE: 30

```
---

*Any comments and suggestions are welcome*
   

