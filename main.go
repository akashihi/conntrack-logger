/*
    conntrack-logger
    Copyright (C) 2015 Denis V Chapligin <akashihi@gmail.com>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

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

  //Start puring
  go purgeDB(configuration)

  //Connect to Netlink
  ct_handle, err := C.nfct_open(NFNL_SUBSYS_CTNETLINK, NF_NETLINK_CONNTRACK_NEW)
  if ct_handle == nil {
    panic(err)
  }
  defer C.nfct_close(ct_handle)

  //Stop netlink on signal
  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
  go func() {
        <-sigs
        log.Print("Exiting...")
	C.nfct_close(ct_handle)
	os.Exit(0)
  }()

  //Increase bufffer
  bsize := C.nfnl_rcvbufsiz(C.nfct_nfnlh(ct_handle), CT_BUFF_SIZE);
  log.Print("Netlink buffer set to: ", bsize)

  //Link netlink and processing function
  C.nfct_callback_register(ct_handle, NFCT_T_NEW, (C.cb)(unsafe.Pointer(C.event_cb_cgo)), nil);

  //Start even processing!
  C.nfct_catch(ct_handle)
}
