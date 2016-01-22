package test;

import android.test.InstrumentationTestCase;

import com.google.protobuf.InvalidProtocolBufferException;

import protocol.Protocol;
import socket.Receiver;
import socket.Socket;

/**
 * Created by wuxiangan on 2016/1/21.
 */
public class TestRegister extends InstrumentationTestCase {
    public void testRegister() {
        Socket socket = new Socket();
        if (!socket.connect("192.168.20.25", 8888)) {
            System.out.println("connect server failed!!!");
            return;
        }
        Protocol.register_request.Builder builder = Protocol.register_request.newBuilder();
        builder.setUsername("18702759796");
        builder.setPassword("wuxiangan");
        Protocol.register_request register_request = builder.build();
        socket.send(10001, register_request.toByteArray(), new Receiver() {
            @Override
            public void handle(byte[] data) {
                Protocol.register_response register_response = null;
                try{
                    register_response = Protocol.register_response.parseFrom(data);
                } catch (InvalidProtocolBufferException e){
                    e.printStackTrace();
                    return;
                }
                System.out.println(register_response.toString());
            }
        });
        while(true){}
    }
}
