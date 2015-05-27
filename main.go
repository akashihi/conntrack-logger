package main

/*
#include <libnetfilter_conntrack/libnetfilter_conntrack.h>
#include <errno.h>

#cgo LDFLAGS: -lnetfilter_conntrack -lnfnetlink

typedef int (*cb)(enum nf_conntrack_msg_type type,
                                            struct nf_conntrack *ct,
                                            void *data);
int event_cb_cgo(int type, struct nf_conntrack *ct, void *data);
*/
import "C"

import(
    "log"
    "log/syslog"
    "os"
    "os/signal"
    "syscall"
    "unsafe"
)

const (
    NF_NETLINK_CONNTRACK_NEW = 0x00000001
    NFNL_SUBSYS_CTNETLINK = 1
    CT_BUFF_SIZE = 8388608
    NFCT_CB_CONTINUE = 1
    NFCT_T_NEW = 1
)

var xml_messages  = make(chan []byte, 128)
var flow_messages = make(chan FlowRecord, 128)

//export event_cb
func event_cb(t C.int, ct *C.struct_nf_conntrack) C.int {
  b := make([]byte, 1024)

  C.nfct_snprintf((*C.char)(unsafe.Pointer(&b[0])), 1024, ct, NFCT_T_NEW, 1, 0)

  //At this stage we have XML structure in the b
  xml_messages <- b

  return NFCT_CB_CONTINUE
}

func main() {
  //Prepare logging
  logwriter, e := syslog.New(syslog.LOG_NOTICE, "conntrack-logger")
  if e == nil {
    log.SetOutput(logwriter)
  }

  //Prepare configuration
  configuration := config()

  //Start parsing and database writing
  go parseNF(xml_messages, flow_messages)
  go writeDB(configuration, flow_messages)

  ct_handle, err := C.nfct_open(NFNL_SUBSYS_CTNETLINK, NF_NETLINK_CONNTRACK_NEW)
  if ct_handle == nil {
    panic(err)
  }
  defer C.nfct_close(ct_handle)

  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
  go func() {
        <-sigs
        log.Print("Exiting...")
	C.nfct_close(ct_handle)
	os.Exit(0)
  }()

  bsize := C.nfnl_rcvbufsiz(C.nfct_nfnlh(ct_handle), CT_BUFF_SIZE);
  log.Print("Netlink buffer set to: ", bsize)

  C.nfct_callback_register(ct_handle, NFCT_T_NEW, (C.cb)(unsafe.Pointer(C.event_cb_cgo)), nil);

  C.nfct_catch(ct_handle)
}
