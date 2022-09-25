package entiy

import (
	"encoding/json"
	"fmt"
	"time"
)

type ProcessIO struct {
	Rchar               int `json:"rchar"`                 // 读出的总字节数，read或者pread()中的长度参数总和（pagecache中统计而来，不代表实际磁盘的读入）
	Wchar               int `json:"wchar"`                 // 写入的总字节数，write或者pwrite中的长度参数总和
	Syscr               int `json:"syscr"`                 // read()或者pread()总的调用次数
	Syscw               int `json:"syscw"`                 // write()或者pwrite()总的调用次数
	ReadBytes           int `json:"read_bytes"`            // 实际从磁盘中读取的字节总数
	WriteBytes          int `json:"write_bytes"`           // 实际写入到磁盘中的字节总数
	CancelledWriteBytes int `json:"cancelled_write_bytes"` // 由于截断pagecache导致应该发生而没有发生的写入字节数（可能为负数）
}

type ProcessStat struct {
	Pid         string //    pid		： 进程ID.
	Comm        string //    comm	: task_struct结构体的进程名
	State       string //    state	: 进程状态, 此处为S
	Ppid        string //    ppid	: 父进程ID （父进程是指通过fork方式，通过clone并非父进程）
	Pgrp        string //    pgrp	：进程组ID
	Session     string //    session	：进程会话组ID
	Tty_nr      string //    tty_nr	：当前进程的tty终点设备号
	Tpgid       string //    tpgid	：控制进程终端的前台进程号
	Flags       int    //    flags	：进程标识位，定义在include/linux/sched.h中的PF_*
	Minflt      int    //    minflt	： 次要缺页中断的次数，即无需从磁盘加载内存页. 比如COW和匿名页
	Cminflt     int    //    cminflt	：当前进程等待子进程的minflt
	Majflt      int    //    majflt	：主要缺页中断的次数，需要从磁盘加载内存页. 比如map文件
	Cmajflt     int    //    cmajflt	：当前进程等待子进程的majflt
	Utime       int    //    utime	: 该进程处于用户态的时间，单位jiffies
	Stime       int    //    stime	: 该进程处于内核态的时间，单位jiffies
	Cutime      int    //    cutime	：当前进程等待子进程的utime
	Cstime      int    //    cstime	: 当前进程等待子进程的utime
	Priority    int    //    priority: 进程优先级, 此次等于10.
	Nice        int    //    nice	: nice值，取值范围[19, -20]
	Num_threads int    //    num_threads: 线程个数, 此处等于221
	Itrealvalue int    //    itrealvalue: 该字段已废弃，恒等于0
	Starttime   int    //    starttime	：自系统启动后的进程创建时间，单位jiffies
	Vsize       int    //    vsize		：进程的虚拟内存大小，单位为bytes
	Rss         int    //    rss		: 进程独占内存+共享库，单位pages
	Rsslim      int    //    rsslim		: rss大小上限
}

type ProcessStatus struct {
	Name                     string `json:"Name"`                       //: 进程对应的应用名称
	Umask                    string `json:"Umask"`                      //: 0022
	State                    string `json:"State"`                      //: D (disk sleep)-----------------------表示此时线程处于sleeping，并且是uninterruptible状态的wait。
	Tgid                     string `json:"Tgid"`                       //: 157-----------------------------------线程组的主pid为157。
	Ngid                     string `json:"Ngid"`                       //: 0
	Pid                      string `json:"Pid"`                        //: 159------------------------------------线程自身的pid为159。
	PPid                     string `json:"PPid"`                       //	: 1-------------------------------------线程组是由init进程创建的。
	TracerPid                string `json:"TracerPid"`                  //: 0
	Uid                      string `json:"Uid"`                        //: 0 0 0 0
	Gid                      string `json:"Gid"`                        //: 0 0 0 0
	FDSize                   string `json:"FDSize"`                     //: 256---------------------------------表示到目前为止进程使用过的描述符总数。
	Groups                   string `json:"Groups"`                     //: 0 10
	VmPeak                   string `json:"VmPeak"`                     //: 1393220 kB--------------------------虚拟内存峰值大小。
	VmSize                   string `json:"VmSize"`                     //: 1390372 kB--------------------------当前使用中的虚拟内存，小于VmPeak。
	VmLck                    string `json:"VmLck"`                      //: 0 kB
	VmPin                    string `json:"VmPin"`                      //: 0 kB
	VmHWM                    string `json:"VmHWM"`                      //: 47940 kB-----------------------------RSS峰值。
	VmRSS                    string `json:"VmRSS"`                      //: 47940 kB-----------------------------RSS实际使用量=RSSAnon+RssFile+RssShmem。
	RssAnon                  string `json:"RssAnon"`                    //: 38700 kB
	RssFile                  string `json:"RssFile"`                    //: 9240 kB
	RssShmem                 string `json:"RssShmem"`                   //	: 0 kB
	VmData                   string `json:"VmData"`                     //: 366648 kB--------------------------进程数据段共366648KB。
	VmStk                    string `json:"VmStk"`                      //: 132 kB------------------------------进程栈一共132KB。
	VmExe                    string `json:"VmExe"`                      //: 84 kB-------------------------------进程text段大小84KB。
	VmLib                    string `json:"VmLib"`                      //: 11488 kB----------------------------进程lib占用11488KB内存。
	VmPTE                    string `json:"VmPTE"`                      //: 1220 kB
	VmSwap                   string `json:"VmSwap"`                     //: 0 kB
	Threads                  string `json:"Threads"`                    //: 40-------------------------------进程中一个40个线程。
	SigQ                     string `json:"SigQ"`                       //: 0/3142------------------------------进程信号队列最大3142，当前没有pending状态的信号。
	SigPnd                   string `json:"SigPnd"`                     //: 0000000000000000------------------没有进程pending，所以位图为0。
	ShdPnd                   string `json:"ShdPnd"`                     //: 0000000000000000
	SigBlk                   string `json:"SigBlk"`                     //: 0000000000000000
	SigIgn                   string `json:"SigIgn"`                     //: 0000000000000006------------------被忽略的信号，对应信号为SIGINT和SIGQUIT，这两个信号产生也不会进行处理。
	SigCgt                   string `json:"SigCgt"`                     //: 0000000180000800------------------已经产生的信号位图，对应信号为SIGUSR2、以及实时信号32和33。
	CapInh                   string `json:"CapInh"`                     //: 0000000000000000
	CapPrm                   string `json:"CapPrm"`                     //: 0000003fffffffff
	CapEff                   string `json:"CapEff"`                     //: 0000003fffffffff
	CapBnd                   string `json:"CapBnd"`                     //: 0000003fffffffff
	CapAmb                   string `json:"CapAmb"`                     //: 0000000000000000
	CpusAllowed              string `json:"Cpus_allowed"`               //: 1---------------------------仅在第1个cpu上执行。
	CpusAllowedList          string `json:"Cpus_allowed_list"`          //: 0
	VoluntaryCtxtSwitches    string `json:"voluntary_ctxt_switches"`    //: 2377-------------线程主动切换2377次。
	NonVoluntaryCtxtSwitches string `json:"nonvoluntary_ctxt_switches"` //: 5 ---------------线程被动切换5次。
}

type ProcessInfo struct {
	Name           string  `json:"name"`
	Pid            string  `json:"pid"`
	CpuUtilization float64 `json:"cpuUtilization"` // Cpu使用率
	//ReadBytes      int     `json:"readBytes"`    // to ioProcess 实际从磁盘中读取的字节总数(这里if=/dev/zero 所以没有实际的读入字节数)
	//WriteBytes     int     `json:"writeBytes"`   // to ioProcess 实际写入到磁盘中的字节总数
	PhyRSS    int   `json:"phyRSS"`         // to processStatRaw vmRss
	VmSize    int   `json:"vmRSS"`          // to processStatRaw vmSize
	Threads   int   `json:"threadCount"`    // to processStatusRaw threads
	Rchar     int   `json:"readCharCount"`  // 读出的总字节数，read或者pread()中的长度参数总和（pagecache中统计而来，不代表实际磁盘的读入）
	Wchar     int   `json:"writeCharCount"` // 写入的总字节数，write或者pwrite中的长度参数总和
	TimeStamp int64 `json:"timeStamp"`      //
}

func (i *ProcessInfo) ToJson() string {
	str, _ := json.Marshal(i)
	return string(str)
}

func (i *ProcessInfo) ToString() string {
	return fmt.Sprintf("name:%s pid:%s cpuUtilizetion:%f phyRss:%d vmRss:%d threadCount:%d readCharCount:%d writeCharCount:%d timeStamp:%d", i.Name, i.Name, i.CpuUtilization, i.PhyRSS, i.VmSize, i.Threads, i.Rchar, i.Wchar, time.Now().UnixNano())
}
