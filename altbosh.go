// Altoros BOSH (altbosh)
// ALTOROS ARGENTINA S.A.
// Developer: Andres Lucas Garcia fiorini
// Date: 09/12/2016

package main

import (

    "crypto/aes"
    "crypto/cipher"
    "crypto/md5"
    "encoding/hex"
    "bytes"
    "bufio"
    "io"
    "io/ioutil"
    "os"
    "fmt"
    "os/exec"
    "log"
    "strings"
    "time"
//    "strconv"
)

//EncryptPass:  encrypt a password. https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/09.5.html
func EncryptPass(cleanpass string) []byte {
//import "crypto/md5"

    // Need to encrypt a string
    plaintext := []byte( cleanpass )
    
    // aes encryption string
    keytext := "ast.x,e(27)8a#ljzmknm.ahkjkljl;k"

    // fmt.Println(len(keytext))

    // Create the aes encryption algorithm 
    c, err := aes.NewCipher([]byte(keytext))
    if err != nil {
        fmt.Printf("Error: NewCipher(%d bytes) = %s", len(keytext), err)
        os.Exit(-1)
    }

    // Encrypted string
    cfb := cipher.NewCFBEncrypter(c, commonIV)
    ciphertext := make([]byte, len(plaintext))
    cfb.XORKeyStream(ciphertext, plaintext)
    // fmt.Printf("%s=>%x\n", plaintext, ciphertext)

return
}

// DecryptPass: decrypts a password returns a []byte 
func DecryptPass (  encryptedpass []byte ) []byte { 
    // Need to encrypt a string
    ciphertext := []byte( encryptedpass )
    
    // aes encryption string
    keytext := "ast.x,e(27)8a#ljzmknm.ahkjkljl;k"

    // fmt.Println(len(keytext))

    // Create the aes encryption algorithm 
    c, err := aes.NewCipher([]byte(keytext))
    if err != nil {
        fmt.Printf("Error: NewCipher(%d bytes) = %s", len(keytext), err)
        os.Exit(-1)
    }

    // Decrypt strings
    cfbdec := cipher.NewCFBDecrypter(c, commonIV)
    plaintextCopy := make([]byte, len(encryptedpass))
    cfbdec.XORKeyStream(plaintextCopy, ciphertext)
    // fmt.Printf("%x=>%s\n", ciphertext, plaintextCopy)

return
}

// CopyFile Copy file contents ## Src: http://www.devdungeon.com/content/working-files-go#copy
func CopyFile (src string, dst string) {
   
    originalFile, err := os.Open( src )
    if err != nil {
        log.Fatal(err)
    }
    defer originalFile.Close()

    
    newFile, err := os.Create( dst )
    if err != nil {
        log.Fatal(err)
    }
    defer newFile.Close()

    // Copy the bytes to destination from source
    bytesWritten, err := io.Copy(newFile, originalFile)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Copied %d bytes.", bytesWritten)
    
    err = newFile.Sync()
    if err != nil {
        log.Fatal(err)
    }
    os.Chmod ( dst, 0600 )
    if err != nil {
        log.Fatal(err)
    }
}
// SetCf: setup a config file 
func SetCf (name string, home string) {

    timestamp := time.Now().Format(time.RFC850)
//    timestampstr := strconv.Itoa(timestamp)

    oldcfstr := []string{home, ".bosh_config" }
    newcfstr := []string{home, ".altbosh", "target", name}
    bkpcfstr := []string{home, ".altbosh", "bkp", timestamp} //fix timestamp, Doesn't work as expected
    oldcf :=  strings.Join(oldcfstr, "/") 
    bkpcf :=  strings.Join(bkpcfstr, "/") 
    newcf :=  strings.Join(newcfstr, "/") 
    fmt.Printf("old: %s\nnew: %s\nbkp: %s\n", oldcf, newcf, bkpcf)
    os.Rename( oldcf, bkpcf )
    CopyFile( newcf, oldcf )

}

// ReadCf: read the content of a cf file
func ReadFl (name string) {

    file, err := os.Open( name )
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
      if( strings.Contains(scanner.Text(), "target: ") )  {
        fmt.Println(scanner.Text())
      }  else  if( strings.Contains(scanner.Text(), "username: ") )  {
        fmt.Println(scanner.Text())
      } /* else  if( strings.Contains(scanner.Text(), "password: ") )  {
        // fmt.Println(scanner.Text())
        passarr := []string(strings.Split(scanner.Text(), " "))
         fmt.Printf("   password encrypted: %s\n", EncryptPass( passarr[5] ) )
        // fmt.Printf("   password decrypted: %s\n", DecryptPass( EncryptPass( passarr[5] ) ) )
      }*/
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

// AddCf: add the content of .bosh_config to the catalog 
func AddCf (home string) {

    now := time.Now()
    ts := now.Unix()
    timestamp := fmt.Sprint(ts)

    curcfstr := []string{home, ".boshconfig" }
    newcfstr := []string{home, ".altbosh", "target", timestamp} //fix timestamp, Doesn't work as expected
    curcf :=  strings.Join(curcfstr, "/") 
    newcf :=  strings.Join(newcfstr, "/") 
    fmt.Printf("new: %s\ncur: %s\n", newcf, curcf)
    CopyFile( curcf, newcf )

}


// ReadCfDir: readdir the content of the catalog 
func ReadCfDir () {

    dir := os.Getenv("HOME")
    foundactive := 0

    fmt.Printf("\x1b[33;1m %s/.bosh_config\t \x1b[0m ", dir )
    boshconfig := dir
    boshconfig += "/.bosh_config"
    hash, err := HashFilemd5( boshconfig )
    currenthash := hash
    if err == nil {
//      fmt.Printf("\x1b[32;1m MD5HASH: %s\t \x1b[0m\n", hash) // This is not elegant, but works
    }
    ReadFl( boshconfig )

    dir += "/.altbosh/target"
    files, err := ioutil.ReadDir(dir)
    if err != nil {
      log.Fatal(err)
    }

    for _, file := range files {
      fmt.Printf("\x1b[33;1m %s\t \x1b[0m ", file.Name())
      filename := []string{dir, file.Name()};
      hash, err := HashFilemd5( strings.Join(filename, "/") )
      if err == nil {
        fmt.Printf("\x1b[32;1m MD5HASH: %s\t", hash)
        if hash == currenthash { 
          foundactive = 1
          fmt.Printf("\x1b[32;1m <<<< Active config\x1b[0m\n", hash) 
        } else {
          fmt.Printf("\x1b[0m\n", hash) 
        }
      }
      ReadFl( strings.Join(filename, "/") )
    }
    if foundactive ==  0 {
        fmt.Printf("No active config file found, you can use the add command.")
    }

}

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
    if e != nil {
        panic(e)
    }
}
// RunCommand: run a command
func RunCommand (  passedtext string ) { 
  cmd := exec.Command(passedtext, "")
  cmd.Stdin = strings.NewReader("some input")
  var out bytes.Buffer
  cmd.Stdout = &out
  err := cmd.Run()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("Command was %s\n", passedtext)
}
func bosh ( passedtext string ) {
   // bosh comamnds    
    
}
// Cli2: eval the command
func Cli2 () {
    scanner := bufio.NewScanner(os.Stdin)
    var text string
    home := os.Getenv("HOME")
    for text != "q" {  // break the loop if text == "q"
        fmt.Print("i% ")
        scanner.Scan()
        text = scanner.Text()
        if text != "q" {
            //fmt.Println("Your text was: ", text)
          if text == "hi" {
            fmt.Println("hi!")
          } else if text == "help" {
            fmt.Print("Interactive Mode Help\n")
            fmt.Print("Available commands are:\n")
            fmt.Print("=======================\n")
            fmt.Print("\thelp:    This help\n")
            fmt.Print("\thi:      say hi\n")
            fmt.Print("\tlist/l:  list bosh targets\n")
            fmt.Print("\tset:     set bosh target\n")
            fmt.Print("\tbosh:    execute a bosh command \n\n")
            fmt.Print("\tdeploy:    Deploy bosh or cf\n\n")
            fmt.Print("\tquit/q   quit\n")
          } else if text == "bosh"  {
              bosh ( text ) 
          } else if text == "add"  {
              AddCf( home )
          } else if text == "clear\n"  {
              RunCommand ( text ) 
          } else if text == "l"  {
            ReadCfDir()
          } else if text == "list"  {
            ReadCfDir()
          } else if strings.HasPrefix(text, "set ")  {
            cmdarr := []string(strings.Split(text, " "))
            fmt.Printf("set bosh target |%s|\n", cmdarr[1] )
            SetCf( cmdarr[1], home )
          } else if text == "q"  {
               os.Exit(0)
          }
        }
    }
}
// Funcion para genera hash md5 de un archivo filePath
// src:  http://www.mrwaggel.be/post/generate-md5-hash-of-a-file/
func HashFilemd5(filePath string) (string, error) {

  var returnMD5String string

  file, err := os.Open(filePath)
  if err != nil {
    return returnMD5String, err
  }
  defer file.Close()

  hash := md5.New()

  if _, err := io.Copy(hash, file); err != nil {
    return returnMD5String, err
  }
  hashInBytes := hash.Sum(nil)[:16]
  returnMD5String = hex.EncodeToString(hashInBytes)

  return returnMD5String, nil

}

// Variable para encriptar
var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

func main() {

   home := os.Getenv("HOME")
   if len(os.Args) > 1 {
    for _,arg := range os.Args {
      if arg == "--help" {
        fmt.Printf ("I need somebody\n")
      } else if arg == "list" {
            ReadCfDir()
            os.Exit(0)
      } else if arg == "l" {
            ReadCfDir()
            os.Exit(0)
      } else if arg == "set" {
          if len(os.Args) > 2 {
            SetCf( os.Args[2], home )
          } else {
            fmt.Printf("You must use a file name as an argument.\n")
          }
            os.Exit(0)
      } else if arg == "-h" {
        fmt.Printf ("I need somebody\n")
      } else if arg == "--version" {
        fmt.Printf ("Version 0.1\n")
      } else if arg == "-v" {
        fmt.Printf ("Debug active\n")
      } else if os.Args[1] == "bosh" {
        fmt.Printf ("bosh command\n")
          if len(os.Args) > 2 {
            SetCf( os.Args[2], home )
          } else {
            fmt.Printf("You provide arguments to the bosh command.\n")
          }
          os.Exit(0)
      } else if arg == "-h" {
        os.Exit(0)
      } else {
        // fmt.Printf("arg %d: %s\n", i, os.Args[i])
      }
    }
    }
    home += "/.altboshrc"
     dat, err := ioutil.ReadFile( home )
    check(err)
    fmt.Print(string(dat))

    // Run interactive commands
    Cli2 ()

}
