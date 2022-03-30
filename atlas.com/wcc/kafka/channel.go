package kafka

type ChannelHandlerFunc[E any] func(byte, byte) HandlerFunc[E]
