# SUDS
SUDS is a **S**imple **U**DP **D**ata **S**tore.  Basically, you send it properly crafted (it's not hard I promise) JSON over UDP packets, it stores that data in a SQLite DB, and will give it back to you via HTTP.  The concept came about for two reasons:

1. I needed to learn Go, so what better way than writing a simple program.
2. I needed a simple "buffer" area for an IoT project I am working on.  Instead of using bulky client TCP code, I figured the code to barf out UDP packets will be significantly smaller and easier to write.

SUDS is currently in 0.1 version.  It's meant to be minimal, not complex, and intentionally not secure.  (So please don't tell me to put all kinds of security stuff in.)  

## Requirements
SUDS has the following requirements:

- https://github.com/mattn/go-sqlite3

## Using SUDS
This is pretty darn simple

    go get github.com/uberlinuxguy/suds
    go build github.com/uberlinuxguy/suds
    
This will build suds into your $GOPATH/bin directory.  You can run it from there, but where you run it, SUDS will create a `suds.db` file.  You should always make sure you run SUDS from the same directory to maintain state. Perhaps some day I will build in a preferences to set the db path, name, and some other stuff, but hey, this is 0.1 man!
