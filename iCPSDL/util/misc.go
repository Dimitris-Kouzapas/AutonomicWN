package util

import "strings"
import "errors"

import "io/ioutil"
import "os/exec"
import "runtime"

import (
    "io"
    "os"
)

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/binary"
)


func CapitaliseFirst(s string) string {
    if len(s) == 0 {
        return ""
    }
    // Assuming the first character is ASCII and thus one byte.
    firstChar := strings.ToUpper(s[:1])
    return firstChar + s[1:]
}

func SequenceHash(strings ...string) string {
    hasher := sha256.New() // Create a new SHA256 hasher

    for _, s := range strings {
        hasher.Write([]byte(s)) // Update the hash with each string
    }

    // Calculate the hash of the concatenated bytes and convert it to a hex string
    return hex.EncodeToString(hasher.Sum(nil))
}

func CommutativeHash(strings []string) string {
    var hashResult uint64
    for _, s := range strings {
        hasher := sha256.New()
        hasher.Write([]byte(s))
        hashBytes := hasher.Sum(nil)
        // Take the first 8 bytes to convert to uint64
        thisHash := binary.BigEndian.Uint64(hashBytes[:8])
        hashResult ^= thisHash // XOR to combine hash values
    }

    b := make([]byte, 8)
    // Put the uint64 into the byte slice using BigEndian
    binary.BigEndian.PutUint64(b, hashResult)

    // Convert the final hash result to a hexadecimal string
    return hex.EncodeToString(b) //binary.BigEndian.AppendUint64(nil, hashResult))
}


/*****************************************************************************
 * set data structure
 *****************************************************************************/

type Set struct {
    elements map[interface{}]bool
}

func NewSet() (this *Set) {
    this = new(Set)
    this.elements = make(map[interface{}]bool)
    return
}

func (this *Set) Add(elems ...interface{}) (*Set) {
    for _, elem := range elems {
        this.elements[elem] = true
    }
    return this
}

func (this *Set) Remove(elem interface{}) (*Set) {
    delete(this.elements, elem)
    return this
}

func (this *Set) Contains(elem interface{}) bool {
    _, ok := this.elements[elem]
    return ok
}

func (this *Set) Size() int {
    return len(this.elements)
}

func (this *Set) Elements() []interface{} {
    elements := make([]interface{}, 0)
    for _, elem := range this.elements {
        elements = append(elements, elem)
    }
    return elements
}
/*func (this *Set) Iterator() func() (interface{}, bool) {
    keys := make([]interface{}, 0)
    for key := range this.elements {
        keys = append(keys, key)
    }
    size := len(keys)
    index := 0
    return  func() (interface{}, bool) {
                if index < size {
                    elem := keys[index]
                    index++
                    return elem, true
                }
                return nil, false
            }
}*/

/*func (this *set) Iterate(f func(interface{})) {
    for element := range s.elements {
        f(element)
    }
}*/


/*****************************************************************************
 * file manipulation
 *****************************************************************************/

func ReadFile(filename string) (string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return "", err
    }
    defer file.Close() // Ensure file is closed after function returns

    // Read file contents into a byte slice
    data, err := io.ReadAll(file)
    if err != nil {
        return "", err // Handle error reading the file
    }

    // Convert byte slice to string
    return string(data), nil
}

func CreateFile(filename string, data string) error {
    file, err := os.Create(filename)
    if err != nil { return err }
    _, err = file.Write([]byte(data))
    if err != nil { return err }
    err = file.Close()
    if err != nil { return err }
    return nil
}

func CreateTmpFile(data string, extension string) (string, error) {
    tmpFile, err := ioutil.TempFile("", "output-*." + extension)
    if err != nil { return "", err }

    _, err = tmpFile.Write([]byte(data))
    if err != nil { return "", err }
    fileName := tmpFile.Name()
    err = tmpFile.Close()
    if err != nil { return "", err }
    return fileName, nil
}


func OpenBrowser(filename string) error {
    var cmd *exec.Cmd
    switch runtime.GOOS {
        case "windows":
            cmd = exec.Command("cmd", "/c", "start", filename)
        case "darwin":
            cmd = exec.Command("open", filename)
        case "linux":
            cmd = exec.Command("xdg-open", filename)
        default:
            return errors.New("Unsupported platform.")
    }
    //err = cmd.Run()
    //if err != nil {
        //return err
    //}
    return cmd.Run()
}
