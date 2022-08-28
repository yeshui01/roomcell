协议接口说明
websock的消息包格式:压缩前的原始json_str长度(int类型)+压缩内容(编码格式:msg_id(int),char(1),seq(int),json_str);

消息头[消息class(uint16类型,2个字节),消息type(uint16类型,2个字节)]+客户端用请求id(int32类型,4个字节)+Code(uint16,2个字节)+消息体(proto的二进制数据)

type ClientMsg
{
    MsgClass uint16 // 消息大类别
    MsgType uint16  // 消息大类别下的小类型
    SeqID  int32    // 客户单用请求序列id(用作客户端回调)
    Code uint16     // 错误码
    MsgData []byte  // proto的序列化数据
}