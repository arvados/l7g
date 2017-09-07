package main

import "fmt"
import "os"
import "os/exec"
//import "sys"

import "syscall"
import "bytes"

import "strings"
import "strconv"

func help_exit() {
  fmt.Fprintf(os.Stderr, "l7g\n")
  fmt.Fprintf(os.Stderr, "  assembly [fn] path\n")
  fmt.Fprintf(os.Stderr, "  assembly-range [fn] path\n")
  fmt.Fprintf(os.Stderr, "  tagset [fn] path\n")
  os.Exit(0)
}

var REF_NAME string

func init() {
  REF_NAME = "hg19"
}

var BGZIP string = "/usr/bin/bgzip"

func process_assembly(ifn, path string) error {

  idx_ifn := ifn + ".fwi"
  if strings.HasSuffix(ifn, ".gz") {
    li := strings.LastIndex(ifn, ".gz")
    idx_ifn = ifn[:li] + ".fwi"
  }

  cmd := exec.Command("grep", "-P", "^" + "[^:]*:[^:]*:" + path + "\t", idx_ifn)

  var out bytes.Buffer
  cmd.Stdout = &out
  err := cmd.Run()
  if err!=nil { return err }

  v := strings.Split( strings.Trim( out.String(), "\t\n " ), "\t" )
  sz_str := v[1]
  beg_str := v[2]

  args := []string{"-c", "-b", beg_str, "-s", sz_str, ifn}
  env := os.Environ()

  return syscall.Exec(BGZIP, args, env)
}

func assembly_end(ifn, path string) (int, int, string, string, error) {

  idx_ifn := ifn + ".fwi"
  if strings.HasSuffix(ifn, ".gz") {
    li := strings.LastIndex(ifn, ".gz")
    idx_ifn = ifn[:li] + ".fwi"
  }

  cmd := exec.Command("grep", "-P", "^" + "[^:]*:[^:]*:" + path + "\t", idx_ifn)

  var out bytes.Buffer
  cmd.Stdout = &out
  err := cmd.Run()
  if err!=nil { return 0,0,"","",err }

  v := strings.Split( strings.Trim( out.String(), "\t\n " ), "\t" )
  sz_str := v[1]
  beg_str := v[2]

  info_fields := strings.Split( v[0], ":" )
  ref_name := info_fields[0]
  chrom_name := info_fields[1]


  cmd = exec.Command(BGZIP, "-c", "-b", beg_str, "-s", sz_str, ifn)

  cmd_out,e := cmd.Output()
  if e!=nil { return 0,0,"","",fmt.Errorf(fmt.Sprintf("assembly_end bgzip error: %v", e)) }
  a := strings.Split(strings.Trim(string(cmd_out), "\n"), "\n")
  if len(a)==0 { return 0, 0, "","", fmt.Errorf("no output") }
  n := len(a)

  step_pos := strings.Split( strings.Replace(a[n-1], " ", "", -1), "\t" )
  step,e := strconv.ParseInt(step_pos[0], 16, 64)
  if e!=nil { return 0,0,"","",e}
  pos,e := strconv.ParseInt(step_pos[1], 10, 64)
  if e!=nil { return 0,0,"","",e}
  return int(step),int(pos),ref_name,chrom_name,nil
}

func assembly_range(ifn, path string) error {
  env := os.Environ() ; _ = env

  i_path,e := strconv.ParseInt(path, 16, 64)
  if e!=nil { return fmt.Errorf(fmt.Sprintf("invalid path: %s. %v", path, e)) }

  bfr_path := fmt.Sprintf("%04x", i_path-1)

  _,pos_bfr,_,_,e := assembly_end(ifn, bfr_path)
  if e!=nil { pos_bfr = 0 }

  step_aft,pos_aft,ref_name,chrom_name,e := assembly_end(ifn, path)
  if e!=nil { return fmt.Errorf(fmt.Sprintf("assembly_end error: %v", e))  }

  if pos_bfr > pos_aft { pos_bfr = 0 }

  fmt.Printf("nstep\tbeg\tend\tchrom_name\tref_name\n")
  fmt.Printf("%d\t%d\t%d\t%s\t%s\n", step_aft, pos_bfr, pos_aft, chrom_name, ref_name)
  return nil
}

/*
func process_tagset_fa(ifn, path string) error {
  idx_ifn := ifn + "."

  cmd := exec.Command("grep", "-P", "^" + REF_NAME + ":.*:" + path + "\t", idx_ifn)
  var out bytes.Buffer
  cmd.Stdout = &out
  err := cmd.Run()
  if err!=nil { return err }

  v := strings.Split( strings.Trim( out.String(), "\t\n " ), "\t" )
  sz_str := v[1]
  beg_str := v[2]

  args := []string{"-c", "-b", beg_str, "-s", sz_str, ifn}
  env := os.Environ()


  args := []string{"faidx", "", beg_str, "-s", sz_str, ifn}
  env := os.Environ()
  return syscall.Exec("/usr/local/bin/samtools", args, env)
}
*/

func main() {
  if len(os.Args)<3 {
    help_exit()
  }

  assembly_fn := os.Getenv("L7G_ASSEMBLY")
  tagset_fn := os.Getenv("L7G_TAGSET") ; _ = tagset_fn

  if len(os.Getenv("L7G_REF"))!=0 {
    REF_NAME = os.Getenv("L7G_REF")
  }

  if os.Args[1] == "assembly" {
    path := os.Args[2]
    if len(os.Args)==4 {
      assembly_fn = os.Args[2]
      path = os.Args[3]
    }

    if len(assembly_fn)==0 { help_exit() }

    e := process_assembly(assembly_fn, path)
    if e!=nil {
      fmt.Fprintf(os.Stderr, "process_assembly: %v", e)
      os.Exit(1)
    }
  } else if os.Args[1] == "assembly-range" {
    path := os.Args[2]
    if len(os.Args)==4 {
      assembly_fn = os.Args[2]
      path = os.Args[3]
    }

    if len(assembly_fn)==0 { help_exit() }

    e := assembly_range(assembly_fn, path)
    if e!=nil {
      fmt.Fprintf(os.Stderr, "%v\n", e)
      os.Exit(1)
    }


    /*
  } else if os.Args[1] == "tagset" {
    path := os.Args[2]
    if len(os.Args)==4 {
      tagset_fn = os.Args[2]
      path = os.Args[3]
    }

    if len(tagset_fn)==0 { help_exit() }

    e:=process_tagset_fa(tagset_fn, path)
    if e!=nil {
      fmt.Fprintf(os.Stderr, "%v", e)
      os.Exit(1)
    }
    */

  }

}
