package main

/*
#include <libnetfilter_conntrack/libnetfilter_conntrack.h>
#include <stdio.h>

int event_cb_cgo(int type, struct nf_conntrack *ct, void *data) {
    return event_cb(type, ct);
}
*/
import "C"
