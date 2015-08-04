package main 

import "os"
import "fmt"
import "log"
import "bufio"
import "path/filepath"
import "container/list"
import "strings"
import "os/exec"
import "io/ioutil"
import "time"
import "encoding/json"


type fileinf struct{
	Name string
	Path string
	Packages []string
	ignored bool
	ModTime time.Time
	TemplatedPath string
}
type cachedData struct {
	lastRun time.Time
	files []fileinf
	companyid string 
	productid string
}
var UnownedFiles = list.New() //This is a list of files not owned by anyone
var IgnoredFiles = list.New() //This is a list of files which shall be ignored
var cache cachedData;
// This function reads the configuration files
func readConfiguration() {
	inputFile, err := os.Open("/etc/freezer/freezer.conf")
	if err != nil {
		log.Fatal("Error opening configuration file", err)
	}
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)
 
	// scanner.Scan() advances to the next token returning false if an or was encountered
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text()," ")
		read := true 
		mode := 0
		for i := 0; i < len(tokens) ; i++ {
			
			if strings.HasPrefix(tokens[i],"#") {
				//Ignore this is a comment: mark rest of line as don't read
				
				read =false
			} else {
				if read {
					if i == 0 {
						//This is a command
						switch {
							case tokens[i] == "ignore":
								mode = 1
							default:
								mode = 0
						}
					} else {
						if mode==1 {
							IgnoredFiles.PushFront(tokens[i]);
						}
					}
				}
			}
		}
	}
 
	// When finished scanning if any error other than io.EOF occured
	// it will be returned by scanner.Err().
	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to read configuration file at /etc/freezer/freezer.conf")
		log.Fatal(scanner.Err())
	}
}

func readPrevRun() {
	if _, err := os.Stat("/var/freezer/com.freezer.cache"); os.IsNotExist(err) {
		if _, err := os.Stat("/etc/freezer/freezerfile"); os.IsNotExist(err) {
			fmt.Println("No freezerfile found... starting manual configuration")
			setup()
		}
	}
}
func shouldAnalyse(p string) (b bool) {
	for e := IgnoredFiles.Front(); e != nil; e = e.Next() {
		pattern, ok := e.Value.(string)
		if strings.HasSuffix(pattern,"/")&&(strings.Index(pattern,"*")<0){
			if(strings.HasPrefix(p,pattern)){
				return false
			}
		}else{
			if !ok {
				return false
			}
	 		matched, error := filepath.Match(pattern,p)
	 		//fmt.Println(pattern+": "+p)
	 		//fmt.Println(matched)
			if error != nil {
				return false
			}
			if matched {
				return false
			}
			_,file :=filepath.Split(p)
			matched, error = filepath.Match(file,p)
			if matched {
				return false
			}
		}
	}
	return true
}
func matchesExp(exp,path string){
	/*array= new string[]
	for i :=0; i<len(exp); i++{
		if(exp){

		}
	}*/
}
func debianPkgExists(path string, stats os.FileInfo, err error) error /*(e error,b bool)*/ {
	
		if err!= nil {
			fmt.Println(err);
			return err
		}else{
			if shouldAnalyse(path){
				//fmt.Println(path)
				ownership := exec.Command("dpkg","-S",path)
				_,err := ownership.Output()
				if err !=nil{
				//	fmt.Println(err)
					UnownedFiles.PushFront(path)
				}
				//fmt.Println(string(out))
				
			}
	}
	return nil
}
func modifiedRecently(lastRun time.Time,stats os.FileInfo) {

}
func IsAllowed(name string) bool{
	if(name == "dev"|| name == "proc" || name == "tmp" || name == "var" || name == "mnt" || name == "media" || name == "home" || name == "sys" || name == "root"){
		return false
	}else {
		return true
	}
}
func analyse() {
	
	files, _ := ioutil.ReadDir("/")
    for _, f := range files {
        //fmt.Println(f.Name())
        if f.IsDir() && IsAllowed(f.Name()) {
        	fmt.Println("Analysing /"+f.Name())
        	filepath.Walk("/"+f.Name(),debianPkgExists)
    	}
    //filepath.Walk("/vagrant",debianPkgExists)
    }
    fmt.Println("The following files were found in the root directory but are not managed by dpkg:")
	for e := UnownedFiles.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
	fmt.Println("Would you like to package them as .deb? (y/n)")
	fmt.Println("Checking for config changes...")
	fmt.Println("Checking for randomly misplaced executables")
}
func setup(){
	fmt.Println("Welcome to freezer... this seems to be the first time you are running freezer from this machine");
	
	fmt.Println("We would like to ask you a few questions in order to set up freezer")
	var productid string;
	var companyid string;
	//var url string;
	fmt.Println("Enter your freezer username")
	fmt.Scanln(&companyid)
	fmt.Println("Enter your product name")
	fmt.Scanln(&productid)
	fmt.Println("Using a custom freezer server? (Y/n)")
	fmt.Println("Setting up your freezer distribution")

}
func main() {

	fmt.Println("Freezer starting analysis")
	readConfiguration()
	analyse()

	//TODO: count the number of arguments
	/*
	argsWithoutProg := os.Args[1:]
	switch {
	case argsWithoutProg[0] == "app-store":
		//Sync with app store purchases
	case argsWithoutProg[0] == "package":
		//Initialize			
	case argsWithoutProg[0] == "discover-deps":
		//Initialize		
	case argsWithoutProg[0] == "init":
		//Initialize	
	case argsWithoutProg[0] == "freeze":
		//Analyze filesystem changes detect non-package files and files which belong to packages but have changed
	case argsWithoutProg[0] == "restore":
		//Restore a snapshot of the previous system
	case argsWithoutProg[0] == "ignore":
		//Ignore a file
	case argsWithoutProg[0] == "list":
		//list snapshots
	case argsWithoutProg[0] == "add":
		//Add a path which is not watched
	case argsWithoutProg[0] == "generateKeypair":
		//Generate  a keypair to be used on 
	case argsWithoutProg[0] == "push":	
		//Attempt to find git to push to a server	
	case argsWithoutProg[0] == "mumify":
		//mumifies a package i.e ensures that the system continues to use a specific release of a package even if newer packages are released
	case argsWithoutProg[0] == "remote-exec":
		//remotely executes commands				
	default: 
		//show help message

	}*/
}