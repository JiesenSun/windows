package client;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.annotations.Expose;
import java.nio.ByteBuffer;
import socket.java_socket.BaseData;

/**
 * Created by wuxiangan on 2016/1/7.
 */
public class LoginT extends BaseData {
    @Expose
    public long uid;
    @Expose
    public long sid;
    @Expose
    public String password;

    public void toJson(ByteBuffer byteBuffer) {
        Gson gson = new GsonBuilder().excludeFieldsWithoutExposeAnnotation().create();
        byteBuffer.put(gson.toJson(this).getBytes());
    }

    public BaseData fromJson(ByteBuffer byteBuffer) {
        int dataSize = byteBuffer.limit() - byteBuffer.position();
        byte[] dataBuf = new byte[dataSize];
        byteBuffer.get(dataBuf);
        Gson gson = new GsonBuilder().excludeFieldsWithoutExposeAnnotation().create();
        return  gson.fromJson(new String(dataBuf), this.getClass());
    }
}
