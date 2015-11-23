package main

import("bytes"
       "encoding/gob"
       "fmt"
       "bufio"
       "os"
       "strings")


type Dictionary struct{
    Words map[int][][]string
}

func read()(Dictionary){
    var files []string = []string{"two-syllables.txt", "three-syllables.txt", "four-syllables.txt", "five-syllables.txt", "six-syllables.txt", "seven-syllables.txt", "eight-syllables.txt"}

    dict := new(Dictionary)
    dict.Words = make(map[int][][]string)
    for k,v:= range files{
        file,_ := os.Open("dictionaries/" + v)
        scanner := bufio.NewScanner(file)
        var i int = 0

        //dict.Words[k+2] = make(map[int][]string)

        for scanner.Scan(){
            dict.Words[k+2] = append(dict.Words[k+2], strings.Split(scanner.Text(), ";"))
            i++
        }
    }
    fmt.Println(dict);
    return *dict
}

func write(buffer bytes.Buffer){
    fo,_ := os.Create("out.gob")
    w := bufio.NewWriter(fo)

    w.Write(buffer.Bytes())
    w.Flush()
}


func main(){
    var dataBuffer bytes.Buffer

    enc := gob.NewEncoder(&dataBuffer)

    dict := read()

    err := enc.Encode(dict)

    if(err!=nil){
        fmt.Println("fail")
    }

    write(dataBuffer)
}
