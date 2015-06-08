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
#include <arpa/inet.h>
#include <netinet/in.h>

typedef int (*cb)(enum nf_conntrack_msg_type type,
                                            struct nf_conntrack *ct,
                                            void *data);
int event_cb_cgo(int type, struct nf_conntrack *ct, void *data);
*/
import "C"

import(
    "time"
    "net"
    "errors"
    "encoding/binary"
    "bytes"
)

var flow_messages = make(chan FlowRecord, 128)

func nfct_check_mark(ct *C.struct_nf_conntrack, mark uint32) (bool) {
 if C.nfct_attr_is_set(ct, C.ATTR_MARK)>0 {
    var markbuf *uint32
    markbuf =  (*uint32)(C.nfct_get_attr(ct, C.ATTR_MARK))
    return mark == *markbuf
 }
 return false
}

func nfct_get_ip(ct *C.struct_nf_conntrack, attr4 uint32, attr6 uint32) (string, error) {
 if C.nfct_attr_is_set(ct, attr4)>0 {
    var ipbuf *uint32
    ipbuf =  (*uint32)(C.nfct_get_attr(ct, attr4))
    return net.IPv4(byte(*ipbuf), byte(*ipbuf>>8), byte(*ipbuf>>16), byte(*ipbuf>>24)).String(), nil
 }
 if C.nfct_attr_is_set(ct, attr6)>0 {
     return "", errors.New("Not implemented")
 }
 return "", errors.New("No such attribute")
}

func nfct_get_port(ct *C.struct_nf_conntrack, attr uint32) (int, error) {
 if C.nfct_attr_is_set(ct, attr)>0 {
    var portbuf *uint16
    portbuf = (*uint16)(C.nfct_get_attr(ct, attr))

    buf := new(bytes.Buffer)
    binary.Write(buf, binary.BigEndian, *portbuf)
    binary.Read(buf, binary.LittleEndian, portbuf)
    return (int)(*portbuf), nil
 }
 return 0, errors.New("No such attribute")
}

func nfct_get_proto(ct *C.struct_nf_conntrack) (string, error) {
 if C.nfct_attr_is_set(ct, C.ATTR_ORIG_L4PROTO)>0 {
    var protobuf *uint8
    protobuf =  (*uint8)(C.nfct_get_attr(ct, C.ATTR_ORIG_L4PROTO))
    return IANAProtocols[*protobuf], nil
 }
 return "", errors.New("No such attribute")
}

//export event_cb
func event_cb(t C.int, ct *C.struct_nf_conntrack) C.int {
  if configuration.Mark != 0 {
      if !nfct_check_mark(ct, configuration.Mark) {
	    return C.NFCT_CB_CONTINUE
      }
  }

  result := FlowRecord { 
    TS : time.Now(),
  }
  var err error

  result.Src, err = nfct_get_ip(ct, C.ATTR_IPV4_SRC, C.ATTR_IPV6_SRC)
  if (err != nil) {
    //Don't know how to handle connections without address
    return C.NFCT_CB_CONTINUE
  }
  result.Dst, err = nfct_get_ip(ct, C.ATTR_IPV4_DST, C.ATTR_IPV6_DST)
  if (err != nil) {
    //Don't know how to handle connections without address
    return C.NFCT_CB_CONTINUE
  }  

  result.Sport,err = nfct_get_port(ct, C.ATTR_PORT_SRC); 
  if err!= nil {
	//Don't know how to handle connections without port
        return C.NFCT_CB_CONTINUE
  }
  result.Dport,err = nfct_get_port(ct, C.ATTR_PORT_DST); 
  if err!= nil {
	//Don't know how to handle connections without port
        return C.NFCT_CB_CONTINUE
  }
  result.Proto, err = nfct_get_proto(ct)
  if err!= nil {
	result.Proto = "err" //At least we will know, that it happened
  }

  flow_messages <- result    

  return C.NFCT_CB_CONTINUE
}
