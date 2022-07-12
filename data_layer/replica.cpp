#include "data_layer.h"

namespace data_layer {

Replica::Replica(IStorage storage, OrderLayerProxy order_layer)
    : storage{storage}, order_layer{order_layer} {
    // TODO
}

AppendResponse Replica::append(AppendRequest app_req) {
    // TODO
}

ReadResponse Replica::read(ReadRequest read_req) {
    // TODO
}

}  // namespace data_layer