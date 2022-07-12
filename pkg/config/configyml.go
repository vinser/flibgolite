package config

const CONFIG_YML = `
library:
  # Selfexplained folders
  STOCK: "books/stock"
  NEW: "books/new"
  TRASH: "books/trash"

genres:
  TREE_FILE: "config/genres.xml"
  # Alternative genres tree can be used (Russian only, sorry) 
  # TREE_FILE: "config/alt_genres.xml"
  
database:
  DSN: "dbdata/books.db"
  # Delay before start each new acquisitions processing
  POLL_DELAY: 30 
  # Maximum simultaneous new aquisitios processing threads
  MAX_SCAN_THREADS: 3

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

locales:
  # Locales folder. You can add your own locale file there like en.yml, ru.yml, uk.yml
  DIR: "config/locales"
  # Default english locale for opds feeds (bookreaders opds menu tree) can be changed to:
  # "uk" for Ukrainian, 
  # "ru" for Russian 
  DEFAULT: "en"
  # Accept only these languages publications. Add others if needed please.
  ACCEPTED: "en, ru, uk"
`
