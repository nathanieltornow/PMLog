#pragma once

#include <string>

using color_t = int8_t;
using seqnum_t = int64_t;

// implements an RPC-client to the order layer
class OrderLayerProxy {
   public:
    OrderLayerProxy(std::string order_ip);

    seqnum_t get_order(color_t color, seqnum_t lsn);
};